package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSymbolHash_Entrypoint(t *testing.T) {
	assert.Equal(t, EntrypointHash, SymbolHash("entrypoint"))
}
