//go:build rocksdb

package blockstore

import (
	"fmt"

	"github.com/linxGnu/grocksdb"
)

// MaxRoot returns the last known root slot.
func (d *DB) MaxRoot() (uint64, error) {
	opts := grocksdb.NewDefaultReadOptions()
	iter := d.DB.NewIteratorCF(opts, d.CfRoot)
	defer iter.Close()
	iter.SeekToLast()
	if !iter.Valid() {
		return 0, ErrNotFound
	}
	slot, ok := ParseSlotKey(iter.Key().Data())
	if !ok {
		return 0, fmt.Errorf("invalid key in root cf")
	}
	return slot, nil
}

// GetSlotMeta returns the shredding metadata of a given slot.
func (d *DB) GetSlotMeta(slot uint64) (*SlotMeta, error) {
	key := MakeSlotKey(slot)
	return GetBincode[SlotMeta](d.DB, d.CfMeta, key[:])
}
