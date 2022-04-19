package main

import (
	"context"
	"flag"
	"log"
	"os"

	"cloud.google.com/go/bigtable"
	"google.golang.org/protobuf/encoding/prototext"

	"github.com/certusone/radiance/pkg/ledger_bigtable"
)

var (
	flagBlock = flag.Uint64("block", 0, "Block number to dump")
	flagDump  = flag.Bool("dump", true, "Dump the block to stdout as prototxt")
)

func init() {
	flag.Parse()

	if *flagBlock == 0 {
		log.Fatal("Must specify block number")
	}
}

func main() {
	ctx := context.Background()

	btClient, err := bigtable.NewClient(ctx, "mainnet-beta", "solana-ledger")
	if err != nil {
		log.Fatalf("Could not create client: %v", err)
	}

	table := btClient.Open(ledger_bigtable.BlocksTable)
	rowKey := ledger_bigtable.SlotKey(*flagBlock)

	row, err := table.ReadRow(ctx, rowKey)
	if err != nil {
		log.Fatalf("Could not read row: %v", err)
	}

	block, err := ledger_bigtable.ParseRow(row)
	if err != nil {
		log.Fatalf("Could not parse row: %v", err)
	}

	if block == nil {
		log.Fatalf("Block not found")
	}

	b, err := prototext.MarshalOptions{
		Multiline: true,
		Indent:    "\t",
	}.Marshal(block)
	if err != nil {
		log.Fatalf("Could not marshal block: %v", err)
	}

	log.Printf("Fetched block %v with %d txs", *flagBlock, len(block.Transactions))

	if *flagDump {
		if _, err := os.Stdout.Write(b); err != nil {
			panic(err)
		}
	}
}
