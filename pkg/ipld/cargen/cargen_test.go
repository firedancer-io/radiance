package cargen

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/go-logr/logr/testr"
	"github.com/ipfs/go-cid"
	"github.com/ipld/go-car"
	"github.com/ipld/go-ipld-prime/codec/dagcbor"
	"github.com/ipld/go-ipld-prime/node/basicnode"
	"github.com/multiformats/go-multicodec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.firedancer.io/radiance/pkg/blockstore"
	"go.firedancer.io/radiance/pkg/shred"
	"k8s.io/klog/v2"
)

// mockBlockWalk is a mock implementation of blockstore.BlockWalkI.
type mockBlockWalk struct {
	queue   []*blockstore.SlotMeta
	entries map[*blockstore.SlotMeta][][]shred.Entry
	staged  [][]shred.Entry
	slot    uint64
	left    int64
}

func newMockBlockWalk() *mockBlockWalk {
	return &mockBlockWalk{
		queue:   nil,
		entries: make(map[*blockstore.SlotMeta][][]shred.Entry),
		staged:  nil,
	}
}

func (m *mockBlockWalk) append(meta *blockstore.SlotMeta, entries [][]shred.Entry) {
	m.queue = append(m.queue, meta)
	m.entries[meta] = entries
}

func (m *mockBlockWalk) Seek(slot uint64) (ok bool) {
	for len(m.queue) > 0 && m.queue[len(m.queue)-1].Slot < slot {
		delete(m.entries, m.queue[0])
		m.queue = m.queue[1:]
	}
	return len(m.queue) > 0
}

func (m *mockBlockWalk) SlotsAvailable() uint64 {
	if m.left <= 0 {
		return 0
	}
	return uint64(m.left)
}

func (m *mockBlockWalk) Next() (meta *blockstore.SlotMeta, ok bool) {
	if len(m.queue) == 0 {
		m.left = 0
		return nil, false
	}

	meta = m.queue[0]
	m.queue = m.queue[1:]
	m.left -= int64(meta.Slot) - int64(m.slot)

	m.staged = m.entries[meta]
	delete(m.entries, meta)
	m.slot = meta.Slot

	ok = true
	return
}

func (m *mockBlockWalk) Entries(*blockstore.SlotMeta) ([][]shred.Entry, error) {
	return m.staged, nil
}

func (m *mockBlockWalk) Close() {
	m.queue = nil
	m.entries = nil
	m.staged = nil
}

func TestGen_Empty(t *testing.T) {
	walk := newMockBlockWalk()
	dir := t.TempDir()
	worker, err := NewWorker(dir, 0, walk)
	assert.EqualError(t, err, "slot 0 not available in any DB")
	assert.Nil(t, worker)
}

func TestGen_Split(t *testing.T) {
	klog.SetLogger(testr.New(t))
	t.Cleanup(klog.ClearLogger)

	walk := newMockBlockWalk()
	walk.left = 432000
	walk.append(
		&blockstore.SlotMeta{Slot: 432000, EntryEndIndexes: []uint32{0}},
		[][]shred.Entry{
			{
				{
					Txns: []solana.Transaction{
						{
							Signatures: make([]solana.Signature, 1),
							Message: solana.Message{
								AccountKeys:     make([]solana.PublicKey, 3),
								Header:          solana.MessageHeader{},
								RecentBlockhash: solana.Hash{},
								Instructions: []solana.CompiledInstruction{
									{
										Accounts: make([]uint16, 2),
										Data:     make(solana.Base58, 64),
									},
								},
							},
						},
					},
				},
			},
		})
	walk.append(
		&blockstore.SlotMeta{Slot: 432002, EntryEndIndexes: []uint32{0, 1}},
		[][]shred.Entry{
			{
				{
					Txns: []solana.Transaction{
						{
							Signatures: make([]solana.Signature, 1),
							Message: solana.Message{
								AccountKeys:     make([]solana.PublicKey, 1),
								Header:          solana.MessageHeader{},
								RecentBlockhash: solana.Hash{},
								Instructions: []solana.CompiledInstruction{
									{
										Accounts: make([]uint16, 1),
										Data:     make(solana.Base58, 20),
									},
								},
							},
						},
					},
				},
			},
			{
				{
					Txns: []solana.Transaction{
						{
							Signatures: make([]solana.Signature, 1),
							Message: solana.Message{
								AccountKeys:     make([]solana.PublicKey, 1),
								Header:          solana.MessageHeader{},
								RecentBlockhash: solana.Hash{},
								Instructions: []solana.CompiledInstruction{
									{
										Accounts: make([]uint16, 1),
										Data:     make(solana.Base58, 20),
									},
								},
							},
						},
					},
				},
			},
		})

	dir := t.TempDir()
	worker, err := NewWorker(dir, 1, walk)
	require.NoError(t, err)
	require.NotNil(t, worker)

	worker.CARSize = 1024
	require.NoError(t, worker.Run(context.Background()))

	entries, err := os.ReadDir(dir)
	require.NoError(t, err)

	sort.Slice(entries, func(i, j int) bool {
		return strings.Compare(entries[i].Name(), entries[j].Name()) < 0
	})
	assert.Len(t, entries, 2)

	cases := map[string]struct {
		size int64
		cids []string
	}{
		"ledger-e1-s432000.car": {
			size: 534,
			cids: []string{
				"bafkreicloicbw5yk5bpkthpwgnxrjy4iwqiitbvxs52tbkxm3fqcmgvl7a", // KindTx
				"bafyreiawzny3vxq6bir23qjai6fcujyiciqz7iftnibrfzbuvtoy7fvqtq", // KindEntry
				"bafyreiaihhxk2rfo5lgilus3wp3mqglv2do3zgtfqpp5zfygrpzdbez3am", // KindBlock
			},
		},
		"ledger-e1-s432002.car": {
			size: 782,
			cids: []string{
				"bafkreihgemfmup2imylvtyp2wj3fymakbx2bfyi2rlzxxq4xv3ycgmc3da", // KindTx
				"bafyreih5wmt3ly25phouneamyxg5fs4uu3ma6kdvgwfregnmjsibkuipeq", // KindEntry
				"bafkreihgemfmup2imylvtyp2wj3fymakbx2bfyi2rlzxxq4xv3ycgmc3da", // KindTx
				"bafyreih5wmt3ly25phouneamyxg5fs4uu3ma6kdvgwfregnmjsibkuipeq", // KindEntry
				"bafyreibqnqnmnhunzwcopnc7srlvq5p5bbgex2xspb6ayzf265rcwtaiie", // KindBlock
			},
		},
	}

	for name, tc := range cases {
		filePath := filepath.Join(dir, name)
		t.Run("Parse_"+name, func(t *testing.T) {
			// match file size
			info, err := os.Stat(filePath)
			require.NoError(t, err)
			assert.Equal(t, tc.size, info.Size())

			// ensure CARs decode
			f, err := os.Open(filepath.Join(dir, name))
			require.NoError(t, err)
			defer f.Close()

			rd, err := car.NewCarReader(f)
			require.NoError(t, err)

			var cids []cid.Cid
			for {
				block, err := rd.Next()
				if errors.Is(err, io.EOF) {
					break
				}
				require.NoError(t, err)
				cid := block.Cid()
				cidType := cid.Type()
				t.Logf("CID=%s Multicodec=%#x", cid, cidType)
				switch multicodec.Code(cidType) {
				case multicodec.Raw:
					var tx solana.Transaction
					require.NoError(t, bin.UnmarshalBin(&tx, block.RawData()))
					t.Logf("  Txn: %s", &tx)
				case multicodec.DagCbor:
					decodeOpt := dagcbor.DecodeOptions{
						AllowLinks: true,
					}
					builder := basicnode.Prototype.Any.NewBuilder()
					require.NoError(t, decodeOpt.Decode(builder, bytes.NewReader(block.RawData())))
					node := builder.Build()
					t.Logf("  Entry: %s", node.Kind())
				default:
					panic("Unexpected entry")
				}
				cids = append(cids, cid)
			}

			// match CIDs
			cidStrs := make([]string, len(cids))
			for i, c := range cids {
				cidStrs[i] = c.String()
			}
			assert.Equal(t, tc.cids, cidStrs)
		})
	}
}
