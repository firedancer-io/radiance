package blockstore

// BlockWalkI abstracts iterators over block data.
//
// The main (and only) implementation in this package is BlockWalk.
type BlockWalkI interface {
	Seek(slot uint64) (ok bool)
	SlotsAvailable() (total uint64)
	Next() (meta *SlotMeta, ok bool)
	Entries(meta *SlotMeta) ([]Entries, error)
	Close()
}
