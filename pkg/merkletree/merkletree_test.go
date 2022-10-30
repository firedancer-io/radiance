package merkletree

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashNodes(t *testing.T) {
	cases := map[string]struct {
		leaves [][]byte
		root   string
	}{
		"Empty": {
			leaves: nil,
			root:   "",
		},
		"One": {
			leaves: [][]byte{[]byte("test")},
			root:   "dbebd10e61bc8c28591273feafbbef95d544f874693301d8f7f8e54c6e30058e",
		},
		"Many": {
			leaves: [][]byte{
				// trent's mom is cool
				[]byte("my"), []byte("very"), []byte("eager"), []byte("mother"), []byte("just"),
				[]byte("served"), []byte("us"), []byte("nine"), []byte("pizzas"),
				[]byte("make"), []byte("prime"),
			},
			root: "b40c847546fdceea166f927fc46c5ca33c3638236a36275c1346d3dffb84e1bc",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			actual := HashNodes(tc.leaves)
			actualRoot := actual.GetRoot()
			for _, node := range actual.Nodes {
				t.Log(hex.EncodeToString(node[:]))
			}
			if tc.root == "" {
				assert.Nil(t, actualRoot)
			} else {
				assert.Equal(t, tc.root, hex.EncodeToString(actualRoot[:]))
			}
		})
	}
}

func TestHashNodes_One(t *testing.T) {
	input := []byte("test")
	actual := HashNodes([][]byte{input})
	expected := HashLeaf(input)
	assert.Equal(t, expected, *actual.GetRoot())
}
