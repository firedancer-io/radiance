package fixtures

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// SBF returns the given SBF fixture.
func SBF(t *testing.T, name string) []byte {
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
	data, err := os.ReadFile(filepath.Join(filepath.Dir(file), "sbf", name))
	require.NoError(t, err)
	return data
}
