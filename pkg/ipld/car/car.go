// Package car implements the CARv1 file format.
package car

import (
	"bytes"

	"github.com/ipfs/go-cid"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/datamodel"
	"github.com/multiformats/go-multicodec"
)

// IdentityCID is the "zero-length "identity" multihash with "raw" codec".
//
// This is the best-practices placeholder value to refer to a non-existent or unknown object.
var IdentityCID cid.Cid

func init() {
	id, err := cid.Cast([]byte{0x01, 0x55, 0x00, 0x00})
	if err != nil {
		panic("failed to create zero-length identity multihash with raw codec lmfao")
	}
	IdentityCID = id
}

// Block is a length-cid-data tuple.
// These make up most of CARv1.
//
// See https://ipld.io/specs/transport/car/carv1/#data
type Block struct {
	Length int
	Cid    cid.Cid
	Data   []byte
}

// NewBlockFromRaw creates a new CIDv1 with the given multicodec contentType on the fly.
func NewBlockFromRaw(data []byte, contentType uint64) Block {
	cidBuilder := cid.V1Builder{
		Codec:  contentType,
		MhType: uint64(multicodec.Sha2_256),
	}
	id, err := cidBuilder.Sum(data)
	if err != nil {
		// Something is wrong with go-cid if this fails.
		panic("failed to construct CID: " + err.Error())
	}
	return Block{
		Length: id.ByteLen() + len(data),
		Cid:    id,
		Data:   data,
	}
}

func NewBlockFromCBOR(node datamodel.Node, contentType uint64) (Block, error) {
	// TODO: This could be rewritten as zero-copy
	var buf bytes.Buffer
	if err := dagcbor.Encode(node, &buf); err != nil {
		return Block{}, err
	}
	return NewBlockFromRaw(buf.Bytes(), contentType), nil
}

// TotalLen returns the total length of the block, including the length prefix.
func (b Block) TotalLen() int {
	return leb128Len(uint64(b.Length)) + b.Length
}

// leb128Len is like len(leb128.FromUInt64(x)).
// But without an allocation, therefore should be preferred.
func leb128Len(x uint64) (n int) {
	for {
		x >>= 7
		if x == 0 {
			return
		}
		n++
	}
}
