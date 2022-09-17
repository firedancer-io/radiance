package shred

import "github.com/gagliardetto/solana-go"

type Shred interface {
	CommonHeader() *CommonHeader
	DataHeader() *DataHeader
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

const (
	LegacyDataHeaderLen = 86
)

// NewShredFromSerialized creates a shred object from the given buffer.
//
// The original slice may be deallocated after this function returns.
func NewShredFromSerialized(shred []byte) Shred {
	if len(shred) < 65 {
		return nil
	}
	variant := shred[64]
	switch {
	case variant == LegacyCodeID:
		return LegacyCodeFromPayload(shred)
	case variant == LegacyDataID:
		return LegacyDataFromPayload(shred)
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
	// Size uint16 (?)
}

func (d *DataHeader) LastInSlot() bool {
	return d.Flags&FlagLastShredInSlot != 0
}

type Entry struct {
	NumHashes uint64
	Hash      solana.Hash
	NumTxns   uint64 `bin:"sizeof=Txns"`
	Txns      []solana.Transaction
}
