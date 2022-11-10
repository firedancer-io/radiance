package compactindex

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// DB is a compactindex handle.
type DB struct {
	Header
	Stream io.ReaderAt
}

// Open returns a handle to access a compactindex.
//
// The provided stream must start with the Magic byte sequence.
// Tip: Use io.NewSectionReader to create aligned substreams when dealing with a file that contains multiple indexes.
func Open(stream io.ReaderAt) (*DB, error) {
	// Read the static 32-byte header.
	// Ignore errors if the read fails after filling the buffer (e.g. EOF).
	var fileHeader [headerSize]byte
	n, readErr := stream.ReadAt(fileHeader[:], 0)
	if n < len(fileHeader) {
		// ReadAt must return non-nil error here.
		return nil, readErr
	}
	db := new(DB)
	if err := db.Header.Load(&fileHeader); err != nil {
		return nil, err
	}
	db.Stream = stream
	return db, nil
}

// FindBucket returns a handle to the bucket that might contain the given key.
func (db *DB) FindBucket(key []byte) (*Bucket, error) {
	return db.GetBucket(db.Header.BucketHash(key))
}

// GetBucket returns a handle to the bucket at the given index.
func (db *DB) GetBucket(i uint) (*Bucket, error) {
	if i > uint(db.Header.NumBuckets) {
		return nil, fmt.Errorf("out of bounds bucket index: %d >= %d", i, db.Header.NumBuckets)
	}
	headerOffset := headerSize + int64(i)*bucketOffsetSize

	// Read pointer to bucket header.
	var bucketOffsetBuf [8]byte
	n, readErr := db.Stream.ReadAt(bucketOffsetBuf[:bucketOffsetSize], headerOffset)
	if n < bucketOffsetSize {
		return nil, readErr
	}
	bucketOffset := binary.LittleEndian.Uint64(bucketOffsetBuf[:])
	// Fill bucket handle.
	bucket := &Bucket{
		BucketDescriptor: BucketDescriptor{
			Stride:      db.entryStride(),
			OffsetWidth: intWidth(db.FileSize),
		},
	}
	// Read bucket header.
	bucket.BucketHeader, readErr = readBucketHeader(db.Stream, int64(bucketOffset))
	if readErr != nil {
		return nil, readErr
	}
	bucket.Entries = io.NewSectionReader(db.Stream, int64(bucketOffset)+bucketHeaderSize, int64(bucket.NumEntries)*int64(bucket.Stride))
	return bucket, nil
}

func (db *DB) entryStride() uint8 {
	hashSize := 3 // TODO remove hardcoded constant
	offsetSize := intWidth(db.FileSize)
	return uint8(hashSize) + offsetSize
}

func readBucketHeader(rd io.ReaderAt, at int64) (hdr BucketHeader, err error) {
	var buf [bucketHeaderSize]byte
	var n int
	n, err = rd.ReadAt(buf[:], at)
	if n >= len(buf) {
		err = nil
	}
	hdr.Load(&buf)
	return
}

// Bucket is a database handle pointing to a subset of the index.
type Bucket struct {
	BucketDescriptor
	Entries *io.SectionReader
}

// maxEntriesPerBucket is the hardcoded maximum permitted number of entries per bucket.
const maxEntriesPerBucket = 1 << 24 // (16 * stride) MiB

const targetEntriesPerBucket = 10000

// Load retrieves all entries in the hashtable.
func (b *Bucket) Load(batchSize int) ([]Entry, error) {
	if batchSize <= 0 {
		batchSize = 1
	}
	// TODO bounds check
	if b.NumEntries > maxEntriesPerBucket {
		return nil, fmt.Errorf("refusing to load bucket with %d entries", b.NumEntries)
	}
	entries := make([]Entry, 0, b.NumEntries)

	stride := int(b.Stride)
	buf := make([]byte, batchSize*stride)
	off := int64(0)
	for {
		// Read another chunk.
		n, err := b.Entries.ReadAt(buf, off)
		// Decode all entries in it.
		sub := buf[:n]
		for len(sub) > stride {
			entries = append(entries, b.loadEntry(sub))
			sub = sub[stride:]
			off += int64(stride)
		}
		// Handle error.
		if errors.Is(err, io.EOF) || errors.Is(err, io.ErrUnexpectedEOF) {
			break
		} else if err != nil {
			return nil, err
		}
	}

	return entries, nil
}

var ErrNotFound = errors.New("not found")
