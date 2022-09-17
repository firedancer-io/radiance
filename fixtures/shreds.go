package fixtures

import (
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/require"
)

func CodeShreds(t testing.TB, network string, slot uint64) [][]byte {
	return shreds(t, network, slot, 'c')
}

func DataShreds(t testing.TB, network string, slot uint64) [][]byte {
	return shreds(t, network, slot, 'd')
}

func shreds(t testing.TB, network string, slot uint64, shredType rune) [][]byte {
	dir := Path(t, "shreds", network, strconv.FormatUint(slot, 10))
	entries, err := os.ReadDir(dir)
	require.NoError(t, err, "cannot open shreds")
	sort.Slice(entries, func(i, j int) bool {
		return strings.Compare(entries[i].Name(), entries[j].Name()) < 0
	})
	var shreds [][]byte
	for _, entry := range entries {
		if !entry.Type().IsRegular() {
			continue
		}
		name := entry.Name()
		fileType, _ := utf8.DecodeRuneInString(name)
		if fileType == shredType {
			shred, err := os.ReadFile(filepath.Join(dir, entry.Name()))
			require.NoError(t, err)
			shreds = append(shreds, shred)
		}
	}
	return shreds
}
