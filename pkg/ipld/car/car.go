package car

import (
	"bytes"
	"io"

	"github.com/filecoin-project/go-leb128"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car"
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

type OutStream interface {
	WriteBlock(Block) error
}

// Writer produces CARv1 files with size tracking.
//
// The implementation is kinda memory-efficient.
// Needs up to IPLD block size plus peanuts of memory.
//
// # Rationale
//
// github.com/ipld/go-car/v2 is not helpful because it wants to traverse a complete IPLD link system.
// However we create an IPLD link system (Merkle-DAG) on the fly in a single pass as we read the chain.
// CARv1 is simple enough that we can roll a custom block writer, so no big deal.
type Writer struct {
	out countingWriter
}

// NewWriter creates a new CARv1 Writer and writes the header.
func NewWriter(out io.Writer) (*Writer, error) {
	w := &Writer{out: newCountingWriter(out)}

	// Deliberately using the go-car v0 library here.
	// go-car v2 doesn't seem to expose the CARv1 header format.
	hdr := car.CarHeader{
		Roots:   []cid.Cid{IdentityCID}, // placeholder
		Version: 1,
	}
	if err := car.WriteHeader(&hdr, w.out); err != nil {
		return nil, err
	}

	return w, nil
}

// WriteBlock writes out a length-CID-value tuple.
func (w *Writer) WriteBlock(b Block) (err error) {
	if _, err = w.out.Write(leb128.FromUInt64(uint64(b.Length))); err != nil {
		return err
	}
	if _, err = w.out.Write(b.Cid.Bytes()); err != nil {
		return err
	}
	_, err = w.out.Write(b.Data)
	return
}

func (w *Writer) Offset() int64 {
	return w.out.written()
}

// countingWriter wraps io.Writer, but counts number of written bytes.
// Not thread safe.
type countingWriter struct {
	io.Writer
	n *int64
}

func newCountingWriter(w io.Writer) countingWriter {
	return countingWriter{
		Writer: w,
		n:      new(int64),
	}
}

func (c countingWriter) Write(data []byte) (n int, err error) {
	n, err = c.Writer.Write(data)
	*c.n += int64(n)
	return
}

// written returns number of bytes written so far.
func (c countingWriter) written() int64 {
	return *c.n
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
