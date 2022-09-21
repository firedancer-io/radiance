package main

import (
	"flag"
	"fmt"
	"go.firedancer.io/radiance/pkg/tpu"
	"github.com/gagliardetto/solana-go"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"sort"
)

var (
	flagSigverify = flag.Bool("sigverify", false, "Verify signatures")
)

// pkcon install libpcap-devel

// readPCAP reads a PCAP file and returns a channel of packets.
func readPCAP(file string) chan []byte {
	packets := make(chan []byte)
	go func() {
		defer close(packets)
		handle, err := pcap.OpenOffline(file)
		if err != nil {
			panic(err)
		}
		defer handle.Close()
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			udpLayer := packet.Layer(layers.LayerTypeUDP)
			if udpLayer == nil {
				continue
			}
			packets <- udpLayer.LayerPayload()
		}
	}()
	return packets
}

func main() {
	flag.Parse()
	if flag.NArg() > 1 || flag.NArg() == 0 {
		fmt.Println("Usage: pcap [file]")
		return
	}
	packets := readPCAP(flag.Arg(0))

	signerCount := make(map[solana.PublicKey]uint)

	// replace by hyperloglog or similar structure if memory usage ever becomes an issue
	signatureCount := make(map[solana.Signature]bool)

	n := 0
	invalid := 0

	for p := range packets {
		n++

		// filter impossibly small packets
		if len(p) < 10 {
			invalid++
			continue
		}

		tx, err := tpu.ParseTx(p)
		if err != nil {
			log.Printf("%d: %v %x", n, err, p)
			invalid++
			continue
		}

		if *flagSigverify {
			ok := tpu.VerifyTxSig(tx)
			if !ok {
				fmt.Printf("bad signature on %s\n", tx.Signatures[0])
				continue
			}
		}

		if len(tx.Signatures) > 0 {
			signatureCount[tx.Signatures[0]] = true
		}

		signers := tpu.ExtractSigners(tx)
		for _, signer := range signers {
			signerCount[signer]++
		}
	}

	// sort by count
	var longTail, longTailCnt uint

	var keys []solana.PublicKey
	for k := range signerCount {
		if signerCount[k] < 10 {
			longTail++
			longTailCnt += signerCount[k]
			continue
		}
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return signerCount[keys[i]] > signerCount[keys[j]]
	})
	for _, k := range keys {
		fmt.Printf("%s\t%d\n", k, signerCount[k])
	}

	fmt.Printf("other signers (<10 pkts, %d total)\t%d\n", longTail, longTailCnt)

	log.Printf("%d packets", n)
	log.Printf("%d invalid packets", invalid)
	log.Printf("%d unique signatures", len(signatureCount))
	log.Printf("%d unique signers", len(signerCount))
	log.Printf("packets per signature: %.02f", float64(n)/float64(len(signatureCount)))
}
