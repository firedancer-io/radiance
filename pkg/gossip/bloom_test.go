package gossip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBloom_FilterMath(t *testing.T) {
	assert.Equal(t, uint64(480), uint64(BloomNumBits(100, 0.1)))
	assert.Equal(t, uint64(959), uint64(BloomNumBits(100, 0.01)))
	assert.Equal(t, uint64(14), uint64(BloomNumKeys(1000, 50)))
	assert.Equal(t, uint64(28), uint64(BloomNumKeys(2000, 50)))
	assert.Equal(t, uint64(55), uint64(BloomNumKeys(2000, 25)))
	assert.Equal(t, uint64(1), uint64(BloomNumKeys(20, 1000)))
}

func TestBloom_AddContains(t *testing.T) {
}
