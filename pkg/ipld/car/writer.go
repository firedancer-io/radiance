package car

import (
	"io"

	"github.com/filecoin-project/go-leb128"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car"
)

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

// Written returns the number of bytes written so far.
func (w *Writer) Written() int64 {
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
