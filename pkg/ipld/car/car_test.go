package car

import (
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/filecoin-project/go-leb128"
	"github.com/stretchr/testify/assert"
)

func TestIdentityCIDStr(t *testing.T) {
	assert.Equal(t, "bafkqaaa", IdentityCID.String())
}

func TestLeb128Len(t *testing.T) {
	cases := []struct {
		n   uint64
		len int
	}{
		{0, 1},
		{1, 1},
		{2, 1},
		{3, 1},
		{4, 1},
		{5, 1},
		{63, 1},
		{64, 1},
		{65, 1},
		{100, 1},
		{127, 1},
		{128, 2},
		{129, 2},
		{2141192192, 5},
		{^uint64(0), 10},
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("N_%d", tc.n), func(t *testing.T) {
			assert.Equal(t, tc.len, leb128Len(tc.n))
		})
	}
}

// Sanity-check asserting that Go Uvarint and LEB128 are compatible.
func TestUvarintIsLeb128(t *testing.T) {
	cases := []uint64{
		0, 1, 2, 3, 4, 5, 63, 64,
		65, 100, 127, 128, 129,
		2141192192, ^uint64(0),
	}
	for _, tc := range cases {
		t.Run(fmt.Sprintf("N_%d", tc), func(t *testing.T) {
			var v0Buf [binary.MaxVarintLen64]byte
			v0 := v0Buf[:binary.PutUvarint(v0Buf[:], tc)]
			v1 := leb128.FromUInt64(tc)
			assert.Equal(t, v0, v1)
		})
	}
}
