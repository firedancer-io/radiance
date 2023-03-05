package loader

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.firedancer.io/radiance/pkg/sbpf"
)

func TestSymbolHash_Entrypoint(t *testing.T) {
	assert.Equal(t, sbpf.EntrypointHash, sbpf.SymbolHash("entrypoint"))
}
