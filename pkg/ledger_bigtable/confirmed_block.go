package ledger_bigtable

import (
	"fmt"

	"cloud.google.com/go/bigtable"
	"google.golang.org/protobuf/proto"

	"go.firedancer.io/radiance/third_party/solana_proto/confirmed_block"
)

const (
	compressionHeaderSize = 4
)

func SlotKey(slot uint64) string {
	return fmt.Sprintf("%016x", slot)
}

func ParseRow(row bigtable.Row) (*confirmed_block.ConfirmedBlock, error) {
	var block confirmed_block.ConfirmedBlock
	x := row["x"]
	for _, item := range x {
		if item.Column == "x:proto" {
			b, err := bigtableCompression(item.Value[0]).Uncompress(item.Value[compressionHeaderSize:])
			if err != nil {
				return nil, fmt.Errorf("failed to uncompress block: %w", err)
			}

			if err := proto.Unmarshal(b, &block); err != nil {
				return nil, fmt.Errorf("failed to unmarshal block: %w", err)
			}
			return &block, nil
		}
	}

	return nil, fmt.Errorf("no proto message in row") // might be bincode?
}
