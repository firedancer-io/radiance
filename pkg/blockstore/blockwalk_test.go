//go:build rocksdb

package blockstore

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiWalk_Len(t *testing.T) {
	mw := BlockWalk{
		handles: []WalkHandle{
			{Start: 0, Stop: 16},
			{Start: 14, Stop: 31},
			{Start: 32, Stop: 34},
			{Start: 36, Stop: 100},
		},
	}
	assert.Equal(t, uint64(35), mw.SlotsAvailable())
}
