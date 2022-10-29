package car

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

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
