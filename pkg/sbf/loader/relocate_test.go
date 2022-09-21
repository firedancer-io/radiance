package loader

import (
	"testing"

	"go.firedancer.io/radiance/pkg/sbf"
	"github.com/stretchr/testify/assert"
)

func TestSymbolHash_Entrypoint(t *testing.T) {
	assert.Equal(t, sbf.EntrypointHash, sbf.SymbolHash("entrypoint"))
}
