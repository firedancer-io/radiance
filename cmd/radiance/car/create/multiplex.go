package create

import (
	"fmt"
	"sort"

	"github.com/linxGnu/grocksdb"
	"go.firedancer.io/radiance/pkg/blockstore"
	"k8s.io/klog/v2"
)

// multiWalk walks blocks in ascending order over multiple RocksDB databases.
type multiWalk struct {
	iter    *grocksdb.Iterator
	handles []dbHandle // sorted
}

// seek skips ahead to a specific slot.
// The caller must call multiWalk.next after seek.
func (m *multiWalk) seek(slot uint64) bool {
	for len(m.handles) > 0 {
		h := m.handles[0]
		if slot < h.start {
			// trying to seek to slot below lowest available
			return false
		}
		if slot <= h.stop {
			h.start = slot
			return true
		}
		m.pop()
	}
	return false
}

// len returns the number of contiguous slots that lay ahead.
func (m *multiWalk) len() (total uint64) {
	if len(m.handles) == 0 {
		return 0
	}
	start := m.handles[0].start
	for _, h := range m.handles {
		if h.start > start {
			return
		}
		stop := h.stop + 1
		total += stop - start
		start = stop
	}
	return
}

// next seeks to the next slot.
func (m *multiWalk) next() (meta *blockstore.SlotMeta, ok bool) {
	if len(m.handles) == 0 {
		return nil, false
	}
	h := m.handles[0]
	if m.iter == nil {
		// Open next database
		m.iter = h.db.DB.NewIteratorCF(grocksdb.NewDefaultReadOptions(), h.db.CfMeta)
		key := blockstore.MakeSlotKey(h.start)
		m.iter.Seek(key[:])
	}
	if !m.iter.Valid() {
		// Close current DB and go to next
		m.pop()
		return m.next() // TODO tail recursion optimization?
	}

	// Get key at current position.
	slot, ok := blockstore.ParseSlotKey(m.iter.Key().Data())
	if !ok {
		klog.Exitf("Invalid slot key: %x", m.iter.Key().Data())
	}
	if slot > h.stop {
		m.pop()
		return m.next()
	}
	h.start = slot

	// Get value at current position.
	meta, err := blockstore.ParseBincode[blockstore.SlotMeta](m.iter.Value().Data())
	if err != nil {
		// Invalid slot metas are irrecoverable.
		// The CAR generation process must stop here.
		klog.Errorf("FATAL: invalid slot meta at slot %d, aborting CAR generation: %s", slot, err)
		return nil, false
	}

	// Seek iterator to next slot in chain.
	if len(meta.NextSlots) != 1 {
		// TODO: Does NextSlots indicate the canonical chain?
		klog.Errorf("FATAL: don't know which slot comes after %d, possible values: %v", slot, meta.NextSlots)
		return nil, false
	}

	// Seek iterator to next entry.
	key := blockstore.MakeSlotKey(meta.NextSlots[0])
	m.iter.Seek(key[:])

	return meta, true
}

// get returns the entries at the current cursor.
// Caller must have made an ok call to multiWalk.next before calling this.
func (m *multiWalk) get(meta *blockstore.SlotMeta) ([]blockstore.Entries, error) {
	h := m.handles[0]
	return h.db.GetEntries(meta)
}

// pop closes the current open DB.
func (m *multiWalk) pop() {
	m.iter.Close()
	m.iter = nil
	m.handles[0].db.Close()
	m.handles = m.handles[1:]
}

func (m *multiWalk) close() {
	if m.iter != nil {
		m.iter.Close()
		m.iter = nil
	}
	for _, h := range m.handles {
		h.db.Close()
	}
	m.handles = nil
}

type dbHandle struct {
	db    *blockstore.DB
	start uint64
	stop  uint64 // inclusive
}

// sortDBs detects bounds of each DB and sorts handles.
func sortDBs(h []dbHandle) error {
	for i, db := range h {
		// Find lowest and highest available slot in DB.
		start, err := getLowestCompletedSlot(db.db)
		if err != nil {
			return err
		}
		stop, err := db.db.MaxRoot()
		if err != nil {
			return err
		}
		h[i] = dbHandle{
			start: start,
			stop:  stop,
			db:    db.db,
		}
	}
	sort.Slice(h, func(i, j int) bool {
		return h[i].start < h[j].start
	})
	return nil
}

// getLowestCompleteSlot finds the lowest slot in a RocksDB from which slots are complete onwards.
func getLowestCompletedSlot(d *blockstore.DB) (uint64, error) {
	iter := d.DB.NewIteratorCF(grocksdb.NewDefaultReadOptions(), d.CfMeta)
	defer iter.Close()
	iter.SeekToFirst()

	// The Solana validator periodically prunes old slots to keep database space bounded.
	// Therefore, the first (few) slots might have valid meta entries but missing data shreds.
	// To work around this, we simply start at the lowest meta and iterate until we find a complete entry.

	const maxTries = 32
	for i := 0; iter.Valid() && i < maxTries; i++ {
		slot, ok := blockstore.ParseSlotKey(iter.Key().Data())
		if !ok {
			return 0, fmt.Errorf(
				"getLowestCompletedSlot(%s): choked on invalid slot key: %x", d.DB.Name(), iter.Key().Data())
		}

		// RocksDB row writes are atomic, therefore meta should never be broken.
		// If we fail to decode meta, bail as early as possible, as we cannot guarantee compatibility.
		meta, err := blockstore.ParseBincode[blockstore.SlotMeta](iter.Value().Data())
		if err != nil {
			return 0, fmt.Errorf(
				"getLowestCompletedSlot(%s): choked on invalid meta for slot %d", d.DB.Name(), slot)
		}

		if _, err = d.GetEntries(meta); err == nil {
			// Success!
			return slot, nil
		}

		iter.Next()
	}

	return 0, fmt.Errorf("failed to find a valid complete slot in DB: %s", d.DB.Name())
}
