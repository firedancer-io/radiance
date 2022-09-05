// Package fixtures contains test data
package fixtures

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// Load returns the fixture at the given path.
func Load(t *testing.T, strs ...string) []byte {
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
	parts := make([]string, 1, 1+len(strs))
	parts[0] = filepath.Dir(file)
	parts = append(parts, strs...)
	data, err := os.ReadFile(filepath.Join(parts...))
	require.NoError(t, err)
	return data
}
