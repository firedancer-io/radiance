//go:build rocksdb

package blockstore

import (
	"errors"
	"fmt"

	"github.com/linxGnu/grocksdb"
)

// DB is RocksDB wrapper
type DB struct {
	DB *grocksdb.DB

	CfDefault   *grocksdb.ColumnFamilyHandle
	CfMeta      *grocksdb.ColumnFamilyHandle
	CfRoot      *grocksdb.ColumnFamilyHandle
	CfDataShred *grocksdb.ColumnFamilyHandle
	CfCodeShred *grocksdb.ColumnFamilyHandle
}

// OpenReadOnly attaches to a blockstore in read-only mode.
//
// Attaching to running validators is supported.
// The DB handle will be a point-in-time view at the time of attaching.
func OpenReadOnly(path string) (*DB, error) {
	return open(path, "")
}

// OpenSecondary attaches to a blockstore in secondary mode.
//
// Only read operations are allowed.
// Unlike OpenReadOnly, allows the user to catch up the DB using (*grocksdb.DB).TryCatchUpWithPrimary.
//
// `secondaryPath` points to a directory where the secondary instance stores its info log.
func OpenSecondary(path string, secondaryPath string) (*DB, error) {
	return open(path, secondaryPath)
}

func open(path string, secondaryPath string) (*DB, error) {
	// List all available column families
	dbOpts := grocksdb.NewDefaultOptions()
	allCfNames, err := grocksdb.ListColumnFamilies(dbOpts, path)
	if err != nil {
		return nil, err
	}
	db := new(DB)

	// Create list of known column families
	cfNames := make([]string, 0, len(allCfNames))
	cfOptList := make([]*grocksdb.Options, 0, len(allCfNames))
	var cfHandles []*grocksdb.ColumnFamilyHandle
	handleSlots := make([]**grocksdb.ColumnFamilyHandle, 0, len(allCfNames))
	for _, cfName := range allCfNames {
		handle, cfOpts := getCfOpts(db, cfName)
		if cfOpts == nil {
			continue
		}
		cfNames = append(cfNames, cfName)
		cfOptList = append(cfOptList, cfOpts)
		handleSlots = append(handleSlots, handle)
	}

	var openFn func() (*grocksdb.DB, []*grocksdb.ColumnFamilyHandle, error)
	if secondaryPath != "" {
		openFn = func() (*grocksdb.DB, []*grocksdb.ColumnFamilyHandle, error) {
			return grocksdb.OpenDbAsSecondaryColumnFamilies(
				dbOpts,
				path,
				secondaryPath,
				cfNames,
				cfOptList,
			)
		}
	} else {
		openFn = func() (*grocksdb.DB, []*grocksdb.ColumnFamilyHandle, error) {
			return grocksdb.OpenDbForReadOnlyColumnFamilies(
				dbOpts,
				path,
				cfNames,
				cfOptList,
				/*errorIfWalExists*/ false,
			)
		}
	}

	// Open database
	db.DB, cfHandles, err = openFn()
	if err != nil {
		return nil, err
	}
	if len(cfHandles) != len(cfNames) {
		// This should never happen
		return nil, fmt.Errorf("expected %d handles, got %d", len(cfNames), len(cfHandles))
	}

	// Write handles into DB object
	for i, slot := range handleSlots {
		*slot = cfHandles[i]
	}

	if db.CfMeta == nil {
		return nil, errors.New("missing column family " + CfMeta)
	}
	if db.CfRoot == nil {
		return nil, errors.New("missing column family " + CfRoot)
	}
	if db.CfDataShred == nil {
		return nil, errors.New("missing column family " + CfDataShred)
	}
	if db.CfCodeShred == nil {
		return nil, errors.New("missing column family " + CfCodeShred)
	}

	return db, nil
}

func getCfOpts(db *DB, name string) (**grocksdb.ColumnFamilyHandle, *grocksdb.Options) {
	switch name {
	case CfDefault:
		return &db.CfDefault, grocksdb.NewDefaultOptions()
	case CfMeta:
		return &db.CfMeta, grocksdb.NewDefaultOptions()
	case CfRoot:
		return &db.CfRoot, grocksdb.NewDefaultOptions()
	case CfDataShred:
		return &db.CfDataShred, grocksdb.NewDefaultOptions()
	case CfCodeShred:
		return &db.CfCodeShred, grocksdb.NewDefaultOptions()
	default:
		return nil, nil
	}
}

func (d *DB) Close() {
	d.DB.Close()
}
