package main

import (
	"fmt"
	"github.com/certusone/tpuproxy/pkg/tpu"
	"github.com/gagliardetto/solana-go"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"sort"
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
	packets := readPCAP("fixtures/tpu.pcap")

	signerCount := make(map[solana.PublicKey]uint)

	for p := range packets {
		tx, err := tpu.ParseTx(p)
		if err != nil {
			fmt.Println(err)
			continue
		}

		ok := tpu.VerifyTxSig(tx)
		if !ok {
			fmt.Printf("bad signature on %s", tx.Signatures[0])
			continue
		}

		signers := tpu.ExtractSigners(tx)
		for _, signer := range signers {
			signerCount[signer]++
		}
	}

	// sort by count
	var keys []solana.PublicKey
	for k := range signerCount {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool {
		return signerCount[keys[i]] > signerCount[keys[j]]
	})
	for _, k := range keys {
		fmt.Printf("%s\t%d\n", k, signerCount[k])
	}
}
