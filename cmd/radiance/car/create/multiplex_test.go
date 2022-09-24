package create

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultiWalk_Len(t *testing.T) {
	mw := multiWalk{
		handles: []dbHandle{
			{start: 0, stop: 16},
			{start: 14, stop: 31},
			{start: 32, stop: 34},
			{start: 36, stop: 100},
		},
	}
	assert.Equal(t, uint64(35), mw.len())
}
