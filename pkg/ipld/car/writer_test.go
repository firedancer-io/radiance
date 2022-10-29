package car

import (
	"bytes"
	"context"
	"io"
	"testing"

	"github.com/ipfs/go-cid"
	ds "github.com/ipfs/go-datastore"
	blockstore "github.com/ipfs/go-ipfs-blockstore"
	cbornode "github.com/ipfs/go-ipld-cbor"
	"github.com/ipld/go-car"
	"github.com/multiformats/go-multicodec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewWriter(t *testing.T) {
	var buf bytes.Buffer
	w, err := NewWriter(&buf)
	require.NoError(t, err)
	require.NotNil(t, w)

	// Ensure that CAR header has been written.
	assert.Equal(t, []byte{
		0x19, // length 0x19 follows
		0xa2, // map with two items

		0x65, 0x72, 0x6f, 0x6f, 0x74, 0x73, // map key: "roots"
		0x81,                                           // map value: array of one root
		0xd8, 0x2a, 0x45, 0x00, 0x01, 0x55, 0x00, 0x00, // CID 00 01 55 00 00 (identity)

		0x67, 0x76, 0x65, 0x72, 0x73, 0x69, 0x6f, 0x6e, // map key: "version"
		0x01, // integer 0x01
	}, buf.Bytes())

	// Ensure CBOR can be deserialized.
	var x struct {
		Roots   []cid.Cid `refmt:"roots"`
		Version uint64    `refmt:"version"`
	}
	cbornode.RegisterCborType(x)
	err = cbornode.DecodeInto(buf.Bytes()[1:], &x)
	require.NoError(t, err)
	assert.Equal(t, []cid.Cid{IdentityCID}, x.Roots)
	assert.Equal(t, uint64(1), x.Version)
}

func TestNewWriter_Error(t *testing.T) {
	var mock mockWriter
	mock.err = io.ErrClosedPipe

	w, err := NewWriter(&mock)
	assert.Nil(t, w)
	assert.Same(t, mock.err, err)
}

func TestWriter(t *testing.T) {
	var buf bytes.Buffer
	w, err := NewWriter(&buf)
	require.NoError(t, err)
	require.NotNil(t, w)

	// Write a bunch of data
	require.NoError(t, w.WriteBlock(NewBlockFromRaw([]byte{}, uint64(multicodec.Raw))))
	require.NoError(t, w.WriteBlock(NewBlockFromRaw([]byte("hello world"), uint64(multicodec.Raw))))

	assert.Equal(t, int64(buf.Len()), w.Written())

	// Load using ipld/go-car library
	ctx := context.Background()
	store := blockstore.NewBlockstore(ds.NewMapDatastore())
	ch, err := car.LoadCar(ctx, store, &buf)
	require.NoError(t, err)
	assert.Equal(t, &car.CarHeader{
		Roots:   []cid.Cid{IdentityCID},
		Version: 1,
	}, ch)

	allKeys, err := store.AllKeysChan(ctx)
	require.NoError(t, err)
	var keys []cid.Cid
	for key := range allKeys {
		keys = append(keys, key)
		t.Log(key.String())
	}
	assert.Len(t, keys, 2)

	b, err := store.Get(ctx, cid.MustParse("bafkreihdwdcefgh4dqkjv67uzcmw7ojee6xedzdetojuzjevtenxquvyku"))
	require.NoError(t, err)
	assert.Equal(t, b.RawData(), []byte{})

	b, err = store.Get(ctx, cid.MustParse("bafkreifzjut3te2nhyekklss27nh3k72ysco7y32koao5eei66wof36n5e"))
	require.NoError(t, err)
	assert.Equal(t, b.RawData(), []byte("hello world"))
}

func TestCountingWriter(t *testing.T) {
	var mock mockWriter

	w := newCountingWriter(&mock)
	assert.Equal(t, int64(0), w.written())

	// successful write
	mock.n, mock.err = 5, nil
	n, err := w.Write([]byte("hello"))
	assert.Equal(t, mock.n, n)
	assert.Equal(t, mock.err, err)
	assert.Equal(t, int64(5), w.written())

	// partial write
	mock.n, mock.err = 3, io.ErrClosedPipe
	n, err = w.Write([]byte("hello"))
	assert.Equal(t, mock.n, n)
	assert.Equal(t, mock.err, err)
	assert.Equal(t, int64(8), w.written())

	// failed write
	mock.n, mock.err = 0, io.ErrClosedPipe
	n, err = w.Write([]byte("hello"))
	assert.Equal(t, mock.n, n)
	assert.Equal(t, mock.err, err)
	assert.Equal(t, int64(8), w.written())
}

type mockWriter struct {
	n   int
	err error
}

func (m *mockWriter) Write(_ []byte) (int, error) {
	return m.n, m.err
}
