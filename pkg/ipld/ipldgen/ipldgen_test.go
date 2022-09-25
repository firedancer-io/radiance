package ipldgen

import (
	"fmt"
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCIDLen(t *testing.T) {
	// Check whether codecs actually result in a CID sized CIDLen.
	// This is important for our allocation strategies during merklerization.
	codecs := []uint64{
		SolanaTx,
		RadianceTxList,
		RadianceEntry,
		RadianceBlock,
		RadianceLedger,
	}
	for _, codec := range codecs {
		t.Run(fmt.Sprintf("Codec_%#x", codec), func(t *testing.T) {
			builder := cid.V1Builder{
				Codec:  codec,
				MhType: multihash.SHA2_256,
			}
			id, err := builder.Sum(nil)
			require.NoError(t, err)
			assert.Equal(t, id.ByteLen(), CIDLen)
		})
	}
}
