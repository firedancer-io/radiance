package compactindex

import (
	"context"
	"encoding/binary"
	"errors"
	"io"
	"io/fs"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {
	const numKeys = uint(1000000)
	const keySize = uint(16)
	const maxOffset = uint64(1000000)
	const queries = int(100000)

	// Create new builder session.
	builder, err := NewBuilder("", numKeys, maxOffset)
	require.NoError(t, err)
	require.NotEmpty(t, builder.buckets)

	// Ensure we cleaned up after ourselves.
	defer func() {
		_, statErr := os.Stat(builder.dir)
		assert.Truef(t, errors.Is(statErr, fs.ErrNotExist), "Delete failed: %v", statErr)
	}()
	defer builder.Close()

	// Insert items to temp buckets.
	preInsert := time.Now()
	key := make([]byte, keySize)
	for i := uint(0); i < numKeys; i++ {
		binary.LittleEndian.PutUint64(key, uint64(i))
		err := builder.Insert(key, uint64(rand.Int63n(int64(maxOffset))))
		require.NoError(t, err)
	}
	t.Logf("Inserted %d keys in %s", numKeys, time.Since(preInsert))

	// Create file for final index.
	targetFile, err := os.CreateTemp("", "compactindex-final-")
	require.NoError(t, err)
	defer os.Remove(targetFile.Name())

	// Seal to final index.
	preSeal := time.Now()
	sealErr := builder.Seal(context.TODO(), targetFile)
	require.NoError(t, sealErr, "Seal failed")
	t.Logf("Sealed in %s", time.Since(preSeal))

	// Print some stats.
	targetStat, err := targetFile.Stat()
	require.NoError(t, err)
	t.Logf("Index size: %d", targetStat.Size())
	t.Logf("Bytes per entry: %f", float64(targetStat.Size())/float64(numKeys))
	t.Logf("Indexing speed: %f/s", float64(numKeys)/time.Since(preInsert).Seconds())

	// Open index.
	_, seekErr := targetFile.Seek(0, io.SeekStart)
	require.NoError(t, seekErr)
	db, err := Open(targetFile)
	require.NoError(t, err, "Failed to open generated index")

	// Run query benchmark.
	preQuery := time.Now()
	for i := queries; i != 0; i-- {
		keyN := uint64(rand.Int63n(int64(numKeys)))
		binary.LittleEndian.PutUint64(key, keyN)

		bucket, err := db.FindBucket(key)
		require.NoError(t, err)

		hash := EntryHash64(bucket.HashDomain, key) & 0xFFFFFF // TODO fix mask
		entries, err := bucket.Load(1024)
		require.NoError(t, err)

		// XXX use external binary search here
		entry := SearchSortedEntries(entries, hash)
		require.NotNil(t, entry)
	}
	t.Logf("Queried %d items", queries)
	t.Logf("Query speed: %f/s", float64(queries)/time.Since(preQuery).Seconds())
}
