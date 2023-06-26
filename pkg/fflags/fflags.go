// Package fflags manages Solana protocol features.
//
// The feature mechanism coordinates the activation of (breaking)
// changes to the Solana protocol.
package fflags

import (
	"go.firedancer.io/radiance/pkg/solana"
)

// Feature is an opaque handle to a feature flag.
type Feature uint

type featureInfo struct {
	name string
	gate solana.Address
}

// seq is the sequence number of allocating feature flag IDs.
// Zero is the sentinel value.
var seq Feature

// featureMap maps feature handle numbers to feature gate addresses.
var featureMap = make(map[Feature]featureInfo)

// Register creates a new application-wide feature flag for the given
// feature gate address. Returns an opaque handle number. Panics on
// duplicate gate address. Not thread-safe -- should be only called
// from the init/main goroutine.
func Register(gate solana.Address, name string) Feature {
	seq++
	s := seq
	if _, ok := featureMap[s]; ok {
		panic("duplicate feature flag handle")
	}
	featureMap[s] = featureInfo{
		name: name,
		gate: gate,
	}
	return s
}

// Features is a set of feature flags.
type Features struct {
	buckets []uint32
}

func (s *Features) set(idx uint, v uint) {
	if idx == 0 || idx > uint(seq) {
		panic("invalid feature flag handle")
	}
	bucket := int(idx / 32)
	if bucket >= len(s.buckets) {
		s.buckets = append(s.buckets, make([]uint32, bucket-len(s.buckets)+1)...)
	}
	s.buckets[bucket] |= 1 << (idx % 32)
}

// HasFeature returns true if the given feature flag is set.
func (s *Features) HasFeature(flag Feature) bool {
	return s.buckets[uint(flag)/32]&(1<<(uint(flag)%32)) != 0
}

// WithFeature modifies s to include the given feature flag.
// Returns s to support chaining-style syntax. Panics on invalid handle.
func (s *Features) WithFeature(flag Feature) *Features {
	s.set(uint(flag), 1)
	return s
}

// WithoutFeature modifies s to exclude the given feature flag.
// Returns s to support chaining-style syntax. Panics on invalid handle.
func (s *Features) WithoutFeature(flag uint) *Features {
	s.set(uint(flag), 0)
	return s
}

// Clone creates a copy of s.
func (s *Features) Clone() *Features {
	c := new(Features)
	c.buckets = make([]uint32, len(s.buckets))
	copy(c.buckets, s.buckets)
	return c
}
