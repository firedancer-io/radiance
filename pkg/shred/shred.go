package shred

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type Shred interface {
	CommonHeader() *CommonHeader
	Data() ([]byte, bool)
	DataComplete() bool
}

const (
	LegacyCodeID = uint8(0b0101_1010)
	LegacyDataID = uint8(0b1010_0101)
	MerkleMask   = uint8(0xF0)
	MerkleCodeID = uint8(0x40)
	MerkleDataID = uint8(0x80)
)

const (
	FlagShredTickReferenceMask = uint8(0b0011_1111)
	FlagDataCompleteShred      = uint8(0b0100_0000)
	FlagLastShredInSlot        = uint8(0b1100_0000)
)

// NewShredFromSerialized creates a shred object from the given buffer.
//
// The original slice may be deallocated after this function returns.
func NewShredFromSerialized(shred []byte, version int) Shred {
	if len(shred) < 65 {
		return nil
	}
	variant := shred[64]
	switch {
	case variant == LegacyCodeID:
		return LegacyCodeFromPayload(shred)
	case variant == LegacyDataID:
		if version <= 1 {
			return LegacyDataFromPayload(shred)
		} else {
			return LegacyDataV2FromPayload(shred)
		}
	case variant&MerkleMask == MerkleCodeID:
		return MerkleCodeFromPayload(shred)
	case variant&MerkleMask == MerkleDataID:
		return MerkleDataFromPayload(shred)
	default:
		return nil
	}
}

type CommonHeader struct {
	Signature   solana.Signature
	Variant     uint8
	Slot        uint64
	Index       uint32
	Version     uint16
	FECSetIndex uint32
}

type DataHeader struct {
	ParentOffset uint16
	Flags        uint8
}

func (d *DataHeader) LastInSlot() bool {
	return d.Flags&FlagLastShredInSlot != 0
}

type DataV2Header struct {
	ParentOffset uint16
	Flags        uint8
	Size         uint16
}

func (d *DataV2Header) LastInSlot() bool {
	return d.Flags&FlagLastShredInSlot != 0
}

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
