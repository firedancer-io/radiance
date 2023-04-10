package dumpblock

import (
	"context"
	"github.com/spf13/cobra"
	"log"
	"os"

	"google.golang.org/protobuf/encoding/prototext"

	"go.firedancer.io/radiance/pkg/ledger_bigtable"
)

var Cmd = cobra.Command{
	Use:   "dumpblock",
	Short: "Dump a block from bigtable in Protobuf text format",
	Run:   run,
}

var (
	flagBlock uint64
	flagDump  bool
)

var flags = Cmd.Flags()

func init() {
	flags.Uint64Var(&flagBlock, "block", 0, "Block number to dump")
	flags.BoolVar(&flagDump, "dump", true, "Dump the block to stdout as prototxt")
}

func run(_ *cobra.Command, _ []string) {
	if flagBlock == 0 {
		log.Fatal("Must specify block number")
	}

	ctx := context.Background()

	btClient, err := ledger_bigtable.MainnetClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create bigtable client: %v", err)
	}

	table := btClient.Open(ledger_bigtable.BlocksTable)
	rowKey := ledger_bigtable.SlotKey(flagBlock)

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

	log.Printf("Fetched block %v with %d txs", flagBlock, len(block.Transactions))

	if flagDump {
		if _, err := os.Stdout.Write(b); err != nil {
			panic(err)
		}
	}
}
