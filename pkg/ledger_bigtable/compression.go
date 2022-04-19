package ledger_bigtable

import (
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/klauspost/compress/zstd"
)

type bigtableCompression uint

const (
	bigtableCompressionNone  bigtableCompression = 0
	bigtableCompressionBzip2 bigtableCompression = 1
	bigtableCompressionGzip  bigtableCompression = 2
	bigtableCompressionZstd  bigtableCompression = 3
)

func (c bigtableCompression) String() string {
	switch c {
	case bigtableCompressionBzip2:
		return "bzip2"
	case bigtableCompressionGzip:
		return "gzip"
	case bigtableCompressionZstd:
		return "zstd"
	default:
		return "none"
	}
}

func (c bigtableCompression) Uncompress(b []byte) ([]byte, error) {
	r := bytes.NewReader(b)
	var o io.Reader
	var err error

	switch c {
	case bigtableCompressionBzip2:
		o = bzip2.NewReader(r)
	case bigtableCompressionGzip:
		o, err = gzip.NewReader(r)
	case bigtableCompressionZstd:
		o, err = zstd.NewReader(r)
	case bigtableCompressionNone:
		return b, nil
	default:
		return nil, fmt.Errorf("unknown compression type: %d", c)
	}
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(o)
}
