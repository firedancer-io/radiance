package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.firedancer.io/radiance/pkg/sbf"
)

func TestSymbolHash_Entrypoint(t *testing.T) {
	assert.Equal(t, sbf.EntrypointHash, sbf.SymbolHash("entrypoint"))
}
