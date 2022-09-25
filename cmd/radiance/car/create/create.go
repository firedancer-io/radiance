package create

import (
	"path/filepath"
	"strconv"

	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/pkg/ipld/cargen"
	"k8s.io/klog/v2"

	"go.firedancer.io/radiance/pkg/blockstore"
)

var Cmd = cobra.Command{
	Use:   "create <epoch>",
	Short: "Create CAR archives from blockstore",
	Long: "Extracts Solana ledger data from blockstore (RocksDB) databases,\n" +
		"and outputs IPLD CARs (content-addressable archives).\n" +
		"\n" +
		"Produces at least one CAR per epoch.\n" +
		"CAR archive contents are deterministic.",
	Args: cobra.ExactArgs(1),
}

// TODO: Actually write the CAR!
//       |
//       Our plan is to transform epochs of Solana history (432000 slots) into batches of CAR files.
//       The CAR output must be byte-by-byte deterministic with regard to Solana's authenticated ledger content.
//       In other words, regardless of which node operator runs this tool, they should always get the same CAR file.
//       |
//       The procedure needs to respect Filecoin's 32GB sector size and split data across multiple CARs if needed.
//       We use Solana blocks as an atomic unit that is never split across CARs.
//       This allows us to assign a slot range to each CAR for the reader's convenience, at negligible alignment cost.
//       |
//       Transforming a single epoch, which takes about a day on mainnet, should take a few hours to transform into CAR.
//       Because of epoch alignment, the CAR generation process can be trivially parallelized by launching multiple instances.
//       In theory, the ledger data extraction process for even a single CAR can be parallelized, at questionable gains.
//       We can synchronize multiple RocksDB iterators that jump over each other block-by-block.
//       CAR writing cannot be parallelized because of strict ordering requirements (determinism).

// TODO: there is a number of things [above] that are conceptually incorrect -- @ribasushi

var flags = Cmd.Flags()

var (
	flagOut = flags.StringP("out", "o", "", "Output directory")
	flagDBs = flags.StringArray("db", nil, "Path to RocksDB (can be specified multiple times)")
)

func init() {
	Cmd.Run = run
}

func run(c *cobra.Command, args []string) {
	outPath := filepath.Clean(*flagOut)
	epochStr := args[0]
	epoch, err := strconv.ParseUint(epochStr, 10, 32)
	if err != nil {
		klog.Exitf("Invalid epoch arg: %s", epochStr)
	}

	// Open blockstores
	dbPaths := *flagDBs
	handles := make([]blockstore.WalkHandle, len(*flagDBs))
	for i := range handles {
		var err error
		handles[i].DB, err = blockstore.OpenReadOnly(dbPaths[i])
		if err != nil {
			klog.Exitf("Failed to open blockstore at %s: %s", dbPaths[i], err)
		}
	}

	// Create new walker object
	walker, err := blockstore.NewBlockWalk(handles)
	if err != nil {
		klog.Exitf("Failed to create multi-DB iterator: %s", err)
	}
	defer walker.Close()

	// Create new cargen worker.
	w, err := cargen.NewWorker(outPath, epoch, walker)
	if err != nil {
		klog.Exitf("Failed to init cargen: %s", err)
	}

	ctx := c.Context()
	if err = w.Run(ctx); err != nil {
		klog.Exitf("FATAL: %s", err)
	}
	klog.Info("DONE")
}
