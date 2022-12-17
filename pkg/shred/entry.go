package shred

import (
	"fmt"

	"github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type Entry struct {
	NumHashes uint64
	Hash      solana.Hash
	Txns      []solana.Transaction
}

func (en *Entry) UnmarshalWithDecoder(decoder *bin.Decoder) (err error) {
	// read the number of hashes:
	if en.NumHashes, err = decoder.ReadUint64(bin.LE); err != nil {
		return fmt.Errorf("failed to read number of hashes: %w", err)
	}
	// read the hash:
	_, err = decoder.Read(en.Hash[:])
	if err != nil {
		return fmt.Errorf("failed to read hash: %w", err)
	}
	// read the number of transactions:
	numTxns, err := decoder.ReadUint64(bin.LE)
	if err != nil {
		return fmt.Errorf("failed to read number of transactions: %w", err)
	}
	if numTxns > uint64(decoder.Remaining()) {
		return fmt.Errorf("not enough bytes to read %d transactions", numTxns)
	}
	// read the transactions:
	en.Txns = make([]solana.Transaction, numTxns)
	for i := uint64(0); i < numTxns; i++ {
		if err = en.Txns[i].UnmarshalWithDecoder(decoder); err != nil {
			return fmt.Errorf("failed to read transaction %d: %w", i, err)
		}
	}
	return
}
