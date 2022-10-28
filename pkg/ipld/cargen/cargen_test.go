package cargen

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.firedancer.io/radiance/pkg/blockstore"
)

// mockBlockWalk is a mock implementation of blockstore.BlockWalkI.
type mockBlockWalk struct {
	queue   []*blockstore.SlotMeta
	entries map[*blockstore.SlotMeta][]blockstore.Entries
	staged  []blockstore.Entries
}

func newMockBlockWalk() *mockBlockWalk {
	return &mockBlockWalk{
		queue:   nil,
		entries: make(map[*blockstore.SlotMeta][]blockstore.Entries),
		staged:  nil,
	}
}

func (m *mockBlockWalk) append(meta *blockstore.SlotMeta, entries []blockstore.Entries) {
	m.queue = append(m.queue, meta)
	m.entries[meta] = entries
}

func (m *mockBlockWalk) Seek(slot uint64) (ok bool) {
	for len(m.queue) > 0 && m.queue[len(m.queue)-1].Slot < slot {
		delete(m.entries, m.queue[0])
		m.queue = m.queue[1:]
	}
	return false
}

func (m *mockBlockWalk) SlotsAvailable() uint64 {
	return uint64(len(m.queue))
}

func (m *mockBlockWalk) Next() (meta *blockstore.SlotMeta, ok bool) {
	if len(m.queue) == 0 {
		return nil, false
	}

	meta = m.queue[0]
	m.queue = m.queue[1:]

	m.staged = m.entries[meta]
	delete(m.entries, meta)

	ok = true
	return
}

func (m *mockBlockWalk) Entries(*blockstore.SlotMeta) ([]blockstore.Entries, error) {
	return m.staged, nil
}

func (m *mockBlockWalk) Close() {
	m.queue = nil
	m.entries = nil
	m.staged = nil
}

func TestGen_Empty(t *testing.T) {
	walk := newMockBlockWalk()
	dir := t.TempDir()
	worker, err := NewWorker(dir, 0, walk)
	assert.EqualError(t, err, "slot 0 not available in any DB")
	assert.Nil(t, worker)
}
