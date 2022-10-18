package poh

import (
	"crypto/sha256"
)

type State struct {
	Entry           Entry
	HashesPerTick   uint64
	RemainingHashes uint64
	TicksPerSlot    uint64
	TickNum         uint64
}

type Entry struct {
	NumHashes uint64
	Hash      [32]byte
}

func (s *State) Record(mixin *[32]byte) Entry {
	h := sha256.New()
	h.Write(s.Entry.Hash[:])
	h.Write(mixin[:])
	h.Sum(s.Entry.Hash[:0])

	entry := s.Entry
	entry.NumHashes++
	s.Entry.NumHashes = 0
	s.RemainingHashes--

	return entry
}

func (s *State) Tick() Entry {
	s.Entry.Hash = sha256.Sum256(s.Entry.Hash[:])

	entry := s.Entry
	entry.NumHashes++
	s.Entry.NumHashes = 0
	s.RemainingHashes = s.HashesPerTick

	return entry
}

func (s *State) Hash(maxIterations uint64) bool {
	numHashes := s.RemainingHashes - 1
	if numHashes > maxIterations {
		numHashes = maxIterations
	}

	for i := uint64(0); i < numHashes; i++ {
		s.Entry.Hash = sha256.Sum256(s.Entry.Hash[:])
	}
	s.Entry.NumHashes += numHashes
	s.RemainingHashes -= numHashes
	return s.RemainingHashes == 1
}
