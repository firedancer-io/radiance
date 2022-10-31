// Package ipldgen transforms Solana ledger data into IPLD DAGs.
package ipldgen

import (
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/ipld/go-ipld-prime/datamodel"
	cidlink "github.com/ipld/go-ipld-prime/linking/cid"
	"go.firedancer.io/radiance/pkg/ipld/car"
	"go.firedancer.io/radiance/pkg/ipld/ipldsch"
	"go.firedancer.io/radiance/pkg/shred"
)

// Multicodec IDs of Solana-related IPLD blocks.
const (
	// SolanaTx is canonical serialization (ledger/wire-format) of a versioned transaction.
	// Stable and upgradable format, unlikely to change soon.
	SolanaTx = 0xc00001

	// RadianceTxList is a list of transactions.
	// DAG-CBOR, non-standard.
	RadianceTxList = 0xc00102

	// RadianceEntry is a representation of a Solana entry, with extra data.
	// DAG-CBOR, non-standard.
	RadianceEntry = 0xc00103

	// RadianceBlock describes a slot and points to entries.
	// DAG-CBOR, non-standard.
	RadianceBlock = 0xc00104
)

// CIDLen is the serialized size in bytes of a Solana/Radiance CID.
const CIDLen = 39

// TargetBlockSize is the target size of variable-length IPLD blocks (e.g. link lists).
// Don't set this to the IPFS max block size, as we might overrun by a few kB.
const TargetBlockSize = 1 << 19

// lengthPrefixSize is the max practical size of an array length prefix.
const lengthPrefixSize = 4

type BlockAssembler struct {
	writer    car.OutStream
	slot      uint64
	entries   []cidlink.Link
	shredding []shredding
}

type shredding struct {
	entryEndIdx uint
	shredEndIdx uint
}

func NewBlockAssembler(writer car.OutStream, slot uint64) *BlockAssembler {
	return &BlockAssembler{
		writer: writer,
		slot:   slot,
	}
}

type EntryPos struct {
	Slot       uint64
	EntryIndex int
	Batch      int
	BatchIndex int
	LastShred  int
}

// WriteEntry appends a ledger entry to the CAR.
func (b *BlockAssembler) WriteEntry(entry shred.Entry, pos EntryPos) error {
	txList, err := NewTxListAssembler(b.writer).Assemble(entry.Txns)
	if err != nil {
		return err
	}
	builder := ipldsch.Type.Entry.NewBuilder()
	entryMap, err := builder.BeginMap(8)
	if err != nil {
		return err
	}

	if pos.LastShred > 0 {
		b.shredding = append(b.shredding, shredding{
			entryEndIdx: uint(pos.EntryIndex),
			shredEndIdx: uint(pos.LastShred),
		})
	}

	var nodeAsm datamodel.NodeAssembler

	nodeAsm, err = entryMap.AssembleEntry("numHashes")
	if err != nil {
		return err
	}
	if err = nodeAsm.AssignInt(int64(entry.NumHashes)); err != nil {
		return err
	}

	nodeAsm, err = entryMap.AssembleEntry("hash")
	if err != nil {
		return err
	}
	if err = nodeAsm.AssignBytes(entry.Hash[:]); err != nil {
		return err
	}

	nodeAsm, err = entryMap.AssembleEntry("txs")
	if err != nil {
		return err
	}
	if err = nodeAsm.AssignNode(txList); err != nil {
		return err
	}

	if err = entryMap.Finish(); err != nil {
		return err
	}
	node := builder.Build()
	block, err := car.NewBlockFromCBOR(node, RadianceEntry)
	if err != nil {
		return err
	}
	b.entries = append(b.entries, cidlink.Link{Cid: block.Cid})
	return b.writer.WriteBlock(block)
}

// Finish appends block metadata to the CAR and returns the root CID.
func (b *BlockAssembler) Finish() (link cidlink.Link, err error) {
	builder := ipldsch.Type.Block.NewBuilder()
	entryMap, err := builder.BeginMap(2)
	if err != nil {
		return link, err
	}

	var nodeAsm datamodel.NodeAssembler

	nodeAsm, err = entryMap.AssembleEntry("slot")
	if err != nil {
		return link, err
	}
	if err = nodeAsm.AssignInt(int64(b.slot)); err != nil {
		return link, err
	}

	nodeAsm, err = entryMap.AssembleEntry("shredding")
	if err != nil {
		return link, err
	}
	list, err := nodeAsm.BeginList(int64(len(b.shredding)))
	if err != nil {
		return link, err
	}
	for _, s := range b.shredding {
		tuple, err := list.AssembleValue().BeginMap(2)
		if err != nil {
			return link, err
		}
		entry, err := tuple.AssembleEntry("entryEndIdx")
		if err != nil {
			return link, err
		}
		if err = entry.AssignInt(int64(s.entryEndIdx)); err != nil {
			return link, err
		}
		entry, err = tuple.AssembleEntry("shredEndIdx")
		if err != nil {
			return link, err
		}
		if err = entry.AssignInt(int64(s.shredEndIdx)); err != nil {
			return link, err
		}
		if err = tuple.Finish(); err != nil {
			return link, err
		}
	}
	if err = list.Finish(); err != nil {
		return link, err
	}

	nodeAsm, err = entryMap.AssembleEntry("entries")
	if err != nil {
		return link, err
	}
	list, err = nodeAsm.BeginList(int64(len(b.entries)))
	if err != nil {
		return link, err
	}
	for _, entry := range b.entries {
		if err = list.AssembleValue().AssignLink(entry); err != nil {
			return link, err
		}
	}
	if err = list.Finish(); err != nil {
		return link, err
	}

	if err = entryMap.Finish(); err != nil {
		return link, err
	}
	node := builder.Build()
	block, err := car.NewBlockFromCBOR(node, RadianceBlock)
	if err != nil {
		return link, err
	}
	if err = b.writer.WriteBlock(block); err != nil {
		return link, err
	}

	return cidlink.Link{Cid: block.Cid}, nil
}

// TxListAssembler produces a Merkle tree of transactions with wide branches.
type TxListAssembler struct {
	writer car.OutStream
	links  []cidlink.Link
}

func NewTxListAssembler(writer car.OutStream) TxListAssembler {
	return TxListAssembler{writer: writer}
}

// Assemble produces a transaction list DAG and returns the root node.
func (t TxListAssembler) Assemble(txs []solana.Transaction) (datamodel.Node, error) {
	for _, tx := range txs {
		if err := t.writeTx(tx); err != nil {
			return nil, err
		}
	}
	return t.finish()
}

// writeTx writes out SolanaTx to the CAR and appends CID to memory.
func (t *TxListAssembler) writeTx(tx solana.Transaction) error {
	buf, err := bin.MarshalBin(tx)
	if err != nil {
		panic("failed to marshal tx: " + err.Error())
	}
	leaf := car.NewBlockFromRaw(buf, SolanaTx)
	if err := t.writer.WriteBlock(leaf); err != nil {
		return err
	}
	t.links = append(t.links, cidlink.Link{Cid: leaf.Cid})
	return nil
}

// finish recursively writes out RadianceTx into a tree structure until the root fits.
func (t *TxListAssembler) finish() (datamodel.Node, error) {
	node, err := t.pack()
	if err != nil {
		return nil, err
	}
	// Terminator: Reached root, stop merklerizing.
	if len(t.links) == 0 {
		return node, nil
	}

	// Create left link.
	block, err := car.NewBlockFromCBOR(node, RadianceTxList)
	if err != nil {
		return nil, err
	}
	var links []cidlink.Link
	links = append(links, cidlink.Link{Cid: block.Cid})
	if err := t.writer.WriteBlock(block); err != nil {
		return nil, err
	}
	// Create right links.
	for len(t.links) > 0 {
		node, err = t.pack()
		if err != nil {
			return nil, err
		}
		block, err = car.NewBlockFromCBOR(node, RadianceTxList)
		if err != nil {
			return nil, err
		}
		if err = t.writer.WriteBlock(block); err != nil {
			return nil, err
		}
		links = append(links, cidlink.Link{Cid: block.Cid})
	}

	// Move up layer.
	t.links = links
	return t.finish()
}

// pack moves as many CIDs as possible into a node.
func (t *TxListAssembler) pack() (node datamodel.Node, err error) {
	builder := ipldsch.Type.TransactionList.NewBuilder()
	list, err := builder.BeginList(0)
	if err != nil {
		return nil, err
	}

	// Pack nodes until we fill TargetBlockSize.
	left := TargetBlockSize - lengthPrefixSize
	for ; len(t.links) > 0 && left >= CIDLen; left -= CIDLen {
		link := t.links[0]
		t.links = t.links[1:]
		if err := list.AssembleValue().AssignLink(link); err != nil {
			return nil, err
		}
	}

	if err := list.Finish(); err != nil {
		return nil, err
	}
	node = builder.Build()
	return node, nil
}
