package replay

import (
	"bytes"

	"github.com/davecgh/go-spew/spew"
	"github.com/gagliardetto/solana-go"
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/pkg/blockstore"
	"go.firedancer.io/radiance/pkg/genesis"
	"go.firedancer.io/radiance/pkg/runtime"
	"go.firedancer.io/radiance/pkg/runtime/poh"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "replay",
	Short: "Replay historical blockchain data",
	Args:  cobra.NoArgs,
	Run:   run,
}

var flags = Cmd.Flags()

var (
	flagGenesis string
	flagDB      string
)

func init() {
	flags.StringVar(&flagGenesis, "genesis", "", "Path to genesis")
	flags.StringVar(&flagDB, "db", "", "Path to RocksDB")
}

func run(c *cobra.Command, _ []string) {
	if flagGenesis == "" {
		klog.Exit("No genesis given")
	}
	if flagDB == "" {
		klog.Exit("No database given")
	}

	// Read genesis, containing the initial set of accounts.
	genesisConfig, genesisHash, err := genesis.ReadGenesisFromFile(flagGenesis)
	if err != nil {
		klog.Exitf("Failed to read genesis: %s", err)
	}
	klog.V(2).Infof("Genesis hash: %s", solana.Hash(*genesisHash))

	// Load initial accounts into memory.
	// Obviously, an in-memory database won't cut it for later stages of replay.
	accounts := runtime.NewMemAccounts()
	genesisConfig.FillAccounts(accounts)

	// Open blockstore database.
	db, err := blockstore.OpenReadOnly(flagDB)
	if err != nil {
		klog.Exitf("Failed to open blockstore: %s", err)
	}

	// Open block iterator.
	walker, err := blockstore.NewBlockWalk([]blockstore.WalkHandle{{DB: db}})
	if err != nil {
		klog.Fatal(err)
	}
	defer walker.Close()

	// PoH delay function (SHA-256 hash chain).
	var chain poh.State
	chain.Entry.Hash = *genesisHash

replay:
	for slot := uint64(0); true; slot++ {
		klog.V(2).Infof("Slot %d", slot)
		meta, ok := walker.Next()
		if !ok {
			break
		}
		entries, err := walker.Entries(meta)
		if err != nil {
			klog.Errorf("Failed to get entries of block %d: %s", slot, err)
			break
		}
		for i, batch := range entries {
			for j, entry := range batch.Entries {
				klog.V(7).Infof("Replay slot=%d entry=%02d/%02d hash=%s txs=%d",
					slot, i, j, entry.Hash, len(entry.Txns))
				chain.Hash(entry.NumHashes)
				for _, tx := range entry.Txns {
					spew.Dump(tx)
				}
				if !bytes.Equal(entry.Hash[:], chain.Entry.Hash[:]) {
					klog.Errorf("PoH mismatch! expected %s, actual %s",
						entry.Hash, solana.Hash(chain.Entry.Hash))
					break replay
				}
			}
		}
	}
}
