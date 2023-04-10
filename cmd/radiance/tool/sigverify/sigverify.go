package main

import (
	"github.com/spf13/cobra"
	"io"
	"os"

	"github.com/gagliardetto/solana-go"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "sigverify",
	Short: "Parse txn and verify signature",
	Run:   run,
}

var (
	flagPublicKey string
	flagSignature string
)

var flags = Cmd.Flags()

func init() {
	flags.StringVar(&flagPublicKey, "pk", "", "Base58-encoded public key")
	flags.StringVar(&flagSignature, "sig", "", "Base58-encoded signature to verify")
}

func run(_ *cobra.Command, _ []string) {
	// Read all of stdin
	msg, err := io.ReadAll(os.Stdin)
	if err != nil {
		klog.Exitf("Failed to read stdin: %v", err)
	}

	// Parse the public key
	pk, err := solana.PublicKeyFromBase58(flagPublicKey)
	if err != nil {
		klog.Exitf("Failed to parse public key: %v", err)
	}

	// Parse the signature
	sig, err := solana.SignatureFromBase58(flagSignature)
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
