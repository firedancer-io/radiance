package ledger_bigtable

import (
	"context"

	"cloud.google.com/go/bigtable"
)

const (
	// BlocksTable is the canonical name of the
	// ConfirmedBlocks table in Solana's BigTable instance.
	BlocksTable = "blocks"
)

func MainnetClient(ctx context.Context) (*bigtable.Client, error) {
	return bigtable.NewClient(ctx, "mainnet-beta", "solana-ledger")
}
