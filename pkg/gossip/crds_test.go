package gossip

import (
	"math"
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCrdsFilterSet(t *testing.T) {
	filters := NewCrdsFilterSet(55345017, 4098)
	assert.Equal(t, 16384, len(filters))
	maskBits := filters[0].MaskBits
	rightShift := 64 - maskBits
	ones := uint64(math.MaxUint64) >> maskBits
	for i, filter := range filters {
		assert.Equal(t, maskBits, filter.MaskBits)
		assert.Equal(t, uint64(i), filter.Mask>>rightShift)
		assert.Equal(t, ones, ones&filter.Mask)
	}
}

func TestCrdsFilterSet_Add(t *testing.T) {
	filters := NewCrdsFilterSet(9672788, 8196)
	hashValues := make([]Hash, 1024)
	for i := range hashValues {
		rand.Read(hashValues[i][:])
		filters.Add(hashValues[i])
	}
	for _, hashValue := range hashValues {
		numHits := uint(0)
		falsePositives := uint(0)
		for _, filter := range filters {
			if filter.TestMask(&hashValue) {
				numHits++
				assert.True(t, filter.Contains(&hashValue))
				assert.True(t, filter.Filter.Contains(&hashValue))
			} else if filter.Filter.Contains(&hashValue) {
				falsePositives++
			}
		}
		assert.Equal(t, numHits, uint(1))
		assert.Less(t, falsePositives, uint(5))
	}
}
