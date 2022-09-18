package create

import (
	"go.firedancer.io/radiance/pkg/blockstore"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "create <epoch...>",
	Short: "Create CAR archives from blockstore",
	Long: "Extracts Solana ledger data from blockstore (RocksDB) databases,\n" +
		"and outputs IPLD CARs (content-addressable archives).\n" +
		"\n" +
		"Produces at least one CAR per epoch.\n" +
		"CAR archive contents are deterministic.",
}

var flags = Cmd.Flags()

var (
	flagDBs = flags.StringArray("db", nil, "Path to RocksDB (can be specified multiple times)")
)

func init() {
	Cmd.Run = run
}

func run(_ *cobra.Command, _ []string) {
	// Open blockstores
	dbPaths := *flagDBs
	handles := make([]dbHandle, len(*flagDBs))
	for i := range handles {
		var err error
		handles[i].db, err = blockstore.OpenReadOnly(dbPaths[i])
		if err != nil {
			klog.Exitf("Failed to open blockstore at %s: %s", dbPaths[i], err)
		}
	}

	// Sort blockstores
	mw := multiWalk{handles: handles}
	defer mw.close()
	if err := sortDBs(mw.handles); err != nil {
		klog.Exitf("Failed to open all DBs: %s", err)
	}

	for {
		// TODO context handling
		meta, ok := mw.next()
		if !ok {
			break
		}
		entries, err := mw.get(meta)
		if err != nil {
			klog.Exitf("FATAL: Failed to get entry at slot %d: %s", meta.Slot, err)
		}
		_ = entries
	}
}
