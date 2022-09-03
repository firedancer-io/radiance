package loader

import (
	"math"
	"math/bits"
)

func clampAddUint64(x uint64, y uint64) uint64 {
	z, carry := bits.Add64(x, y, 0)
	if carry != 0 {
		return math.MaxUint64
	}
	return z
}

type addrRange struct {
	min, max uint64
}

func newAddrRange() addrRange {
	return addrRange{min: math.MaxUint64, max: 0}
}

func (a addrRange) len() uint64 {
	if a.min >= a.max {
		return 0
	}
	return a.max - a.min
}

func (a addrRange) contains(addr uint64) bool {
	if a.len() == 0 {
		return false
	}
	return a.min <= addr && addr < a.max
}

func (a addrRange) containsRange(b addrRange) bool {
	if a.len() == 0 || b.len() == 0 {
		return false
	}
	return a.min <= b.min && a.max >= b.max
}

func (a *addrRange) extendToFit(x uint64) {
	if x < a.min {
		a.min = x
	}
	if x > a.max {
		a.max = x
	}
}

func (a *addrRange) insert(b addrRange) {
	if b.len() == 0 {
		return
	}
	if b.min < a.min {
		a.min = b.min
	}
	if b.max > a.max {
		a.max = b.max
	}
}
