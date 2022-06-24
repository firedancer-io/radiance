package main

import (
	"flag"
	"io/ioutil"
	"os"

	"github.com/gagliardetto/solana-go"
	"k8s.io/klog/v2"
)

var (
	flagPublicKey = flag.String("pk", "", "Base58-encoded public key")
	flagSignature = flag.String("sig", "", "Base58-encoded signature to verify")
)

func init() {
	flag.Parse()
}

func main() {
	// Read all of stdin
	msg, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		klog.Exitf("Failed to read stdin: %v", err)
	}

	// Parse the public key
	pk, err := solana.PublicKeyFromBase58(*flagPublicKey)
	if err != nil {
		klog.Exitf("Failed to parse public key: %v", err)
	}

	// Parse the signature
	sig, err := solana.SignatureFromBase58(*flagSignature)
	if err != nil {
		klog.Exitf("Failed to parse signature: %v", err)
	}

	// Verify the signature
	if !sig.Verify(pk, msg) {
		klog.Exitf("Signature verification failed")
	} else {
		klog.Info("Signature verification succeeded")
	}
}
