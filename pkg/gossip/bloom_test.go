package gossip

import (
	"crypto/sha256"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewBloomRandom(t *testing.T) {
	cases := []struct {
		name      string
		numItems  uint64
		falseRate float64
		maxBits   uint64

		wantKeys int
		wantBits uint64
	}{
		{
			name:     "Empty",
			numItems: 0, falseRate: 0.1, maxBits: 100,
			wantKeys: 0, wantBits: 1,
		},
		{
			name:     "Random",
			numItems: 10, falseRate: 0.1, maxBits: 100,
			wantKeys: 3, wantBits: 48,
		},
		{
			name:     "Random",
			numItems: 100, falseRate: 0.1, maxBits: 100,
			wantKeys: 1, wantBits: 100,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			bloom := NewBloomRandom(tc.numItems, tc.falseRate, tc.maxBits)
			require.NotNil(t, bloom)
			assert.Equal(t, tc.wantKeys, len(bloom.Keys))
			assert.Equal(t, tc.wantBits, bloom.Bits.Len)
		})
	}
}

func TestBloom_FilterMath(t *testing.T) {
	assert.Equal(t, uint64(480), uint64(BloomNumBits(100, 0.1)))
	assert.Equal(t, uint64(959), uint64(BloomNumBits(100, 0.01)))
	assert.Equal(t, uint64(14), uint64(BloomNumKeys(1000, 50)))
	assert.Equal(t, uint64(28), uint64(BloomNumKeys(2000, 50)))
	assert.Equal(t, uint64(55), uint64(BloomNumKeys(2000, 25)))
	assert.Equal(t, uint64(1), uint64(BloomNumKeys(20, 1000)))
}

func TestBloom_AddContains(t *testing.T) {
	bloom := NewBloomRandom(100, 0.1, 100)
	require.NotNil(t, bloom)
	// known keys to avoid false positives in the test
	bloom.Keys = []uint64{0, 1, 2, 3}

	var key [32]byte

	key = sha256.Sum256([]byte("hello"))
	assert.False(t, bloom.Contains(&key))
	bloom.Add(&key)
	assert.True(t, bloom.Contains(&key))

	key = sha256.Sum256([]byte("world"))
	assert.False(t, bloom.Contains(&key))
	bloom.Add(&key)
	assert.True(t, bloom.Contains(&key))
}

func TestBloom_Randomness(t *testing.T) {
	b1 := NewBloomRandom(10, 0.1, 100)
	b2 := NewBloomRandom(10, 0.1, 100)
	require.NotNil(t, b1)
	require.NotNil(t, b2)

	sort.Slice(b1.Keys, func(i, j int) bool { return b1.Keys[i] < b1.Keys[j] })
	sort.Slice(b2.Keys, func(i, j int) bool { return b2.Keys[i] < b2.Keys[j] })

	assert.NotEqual(t, b1.Keys, b2.Keys)
}
