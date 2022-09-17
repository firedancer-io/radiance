// Package fixtures contains test data
package fixtures

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/require"
)

// Path returns the absolute path for the fixture at the given relative path.
func Path(t testing.TB, strs ...string) string {
	_, file, _, ok := runtime.Caller(0)
	require.True(t, ok, "runtime.Caller failed")
	parts := make([]string, 1, 1+len(strs))
	parts[0] = filepath.Dir(file)
	parts = append(parts, strs...)
	return filepath.Join(parts...)
}

// Load returns the fixture at the given path.
func Load(t testing.TB, strs ...string) []byte {
	data, err := os.ReadFile(Path(t, strs...))
	require.NoError(t, err)
	return data
}

// Open returns a file handle for a fixture.
func Open(t testing.TB, strs ...string) *os.File {
	f, err := os.Open(Path(t, strs...))
	require.NoError(t, err)
	return f
}
