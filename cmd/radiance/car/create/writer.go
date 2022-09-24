package create

import (
	"io"

	"github.com/filecoin-project/go-leb128"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car"
	"github.com/multiformats/go-multicodec"
)

// Multicodec IDs of Solana-related IPLD blocks.
const (
	// Solana transaction in canonical serialization (ledger/wire-format).
	// Stable and upgradable format, unlikely to change soon.
	SolanaTx = 0xc00001
)

// block is a length-cid-data tuple.
// These make up most of CARv1.
//
// See https://ipld.io/specs/transport/car/carv1/#data
type block struct {
	length int
	cid    cid.Cid
	data   []byte
}

// Creates a new CIDv1 with the given multicodec contentType on the fly.
func newBlockFomRaw(data []byte, contentType uint64) block {
	cidBuilder := cid.V1Builder{
		Codec:  contentType,
		MhType: multicodec.Sha2_256,
	}
	id, err = cidBuilder.Sum(data)
	if err != nil {
		// Something is wrong with go-cid if this fails.
		panic("failed to construct CID: " + err.Error())
	}
	return block{
		length: id.ByteLen() + len(data),
		cid:    id,
		data:   data,
	}
}

// totalLen returns the total length of the block, including the length prefix.
func (b block) totalLen() int {
	return leb128Len(uint64(b.length)) + b.length
}

// writer produces CARv1 files with size tracking.
//
// The implementation is kinda memory-efficient.
// Needs up to IPLD block size plus peanuts of memory.
type writer struct {
	out countingWriter
}

// newWriter creates a new CARv1 writer and writes the header.
func newWriter(out io.Writer) (*writer, error) {
	w := &writer{out: newCountingWriter(out)}

	// We don't know the root CID.
	// As per best practices we write the
	// "zero-length "identity" multihash with "raw" codec"
	root, _ := cid.Cast([]byte{0x01, 0x55, 0x00, 0x00})
	if root == nil {
		panic("failed to create zero-length identity multihash with raw codec lmfao")
	}

	// Deliberately using the go-car v0 library here.
	// go-car v2 doesn't seem to expose the CARv1 header format.
	hdr := car.CarHeader{
		Roots:   []cid.Cid{root},
		Version: 1,
	}
	if err := car.WriteHeader(&hdr, w.out); err != nil {
		return nil, err
	}

	return w, nil
}

// writeBlock writes out a length-CID-value tuple.
func (w *writer) writeBlock(b block) error {
	if _, err := w.out.Write(leb128.FromUInt64(b.length)); err != nil {
		return err
	}
	if _, err := w.out.Write(b.cid.Bytes()); err != nil {
		return err
	}
	return w.out.Write(b.data)
}

// countingWriter wraps writer, but counts number of written bytes.
// Not thread safe.
type countingWriter struct {
	io.Writer
	n *uint64
}

func newCountingWriter(w io.Writer) countingWriter {
	return countingWriter{
		Writer: w,
		n:      new(uint64),
	}
}

func (c countingWriter) Write(data []byte) (n int, err error) {
	n, err = c.Writer.Write(data)
	*c.n += n
	return
}

// written returns number of bytes written so far.
func (c countingWriter) written() int {
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
