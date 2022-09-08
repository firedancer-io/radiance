// Copyright 2022 Solana Foundation.
// Go port by Richard Patel <me@terorie.dev>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package gossip

// Original Rust source: https://crates.io/crates/solana-bloom

import (
	"math"
	"math/rand"
)

const MaxBloomSize = 928

func NewBloom(numBits uint64, keys []uint64) *Bloom {
	bits := make([]uint64, (numBits+63)/64)
	ret := &Bloom{
		Keys: keys,
		Bits: BitVecU64{
			Bits: BitVecU64Inner{Value: &bits},
			Len:  numBits,
		},
		NumBitsSet: 0,
	}
	return ret
}

func NewBloomRandom(numItems uint64, falseRate float64, maxBits uint64) *Bloom {
	m := BloomNumBits(float64(numItems), falseRate)
	numBits := uint64(m)
	if maxBits < numBits {
		numBits = maxBits
	}
	if numBits == 0 {
		numBits = 1
	}
	numKeys := uint64(BloomNumKeys(float64(numBits), float64(numItems)))
	keys := make([]uint64, numKeys)
	for i := range keys {
		keys[i] = rand.Uint64()
	}
	return NewBloom(numBits, keys)
}

func BloomNumBits(n, p float64) float64 {
	return math.Ceil((n * math.Log(p)) / math.Log(1/math.Pow(2, math.Log(2))))
}

func BloomNumKeys(m, n float64) float64 {
	if n == 0 {
		return 0
	}
	return math.Max(1, math.Round((m/n)*math.Log(2)))
}

func BloomMaxItems(m, p, k float64) float64 {
	return math.Ceil(m / (-k / math.Log(1-math.Exp(math.Log(p)/k))))
}

func BloomMaskBits(numItems, maxItems float64) uint32 {
	return uint32(math.Max(math.Ceil(math.Log2(numItems/maxItems)), 0))
}

func (b *Bloom) Pos(key *Hash, k uint64) uint64 {
	return FNV1a(key[:], k) % b.Bits.Len
}

func (b *Bloom) Clear() {
	bits := *b.Bits.Bits.Value
	for i := range bits {
		bits[i] = 0
	}
	b.NumBitsSet = 0
}

func (b *Bloom) Add(key *Hash) {
	for _, k := range b.Keys {
		pos := b.Pos(key, k)
		if !b.Bits.Get(pos) {
			b.NumBitsSet += 1
			b.Bits.Set(pos, true)
		}
	}
}

func (b *Bloom) Contains(key *Hash) bool {
	for _, k := range b.Keys {
		if !b.Bits.Get(b.Pos(key, k)) {
			return false
		}
	}
	return true
}

func FNV1a(slice []byte, hash uint64) uint64 {
	for _, c := range slice {
		hash ^= uint64(c)
		hash *= 1099511628211
	}
	return hash
}
