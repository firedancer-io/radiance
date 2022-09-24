package create

import (
	"os"
	"strconv"

	"go.firedancer.io/radiance/pkg/blockstore"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "create <out.car> <epoch>",
	Short: "Create CAR archives from blockstore",
	Long: "Extracts Solana ledger data from blockstore (RocksDB) databases,\n" +
		"and outputs IPLD CARs (content-addressable archives).\n" +
		"\n" +
		"Produces at least one CAR per epoch.\n" +
		"CAR archive contents are deterministic.",
	Args: cobra.ExactArgs(2),
}

// TODO: Actually write the CAR!
//       |
//       Our plan is to transform epochs of Solana history (432000 slots) into batches of CAR files.
//       The CAR output must be byte-by-byte deterministic with regard to Solana's authenticated ledger content.
//       In other words, regardless of which node operator runs this tool, they should always get the same CAR file.
//       |
//       github.com/ipld/go-car/v2 is not helpful because it wants to traverse a complete IPLD link system.
//       However we create an IPLD link system (Merkle-DAG) on the fly in a single pass as we read the chain.
//       CARv1 is simple enough that we can roll a custom block writer, so no big deal. Vec<(len, cid, data)>
//       We only need to reserve sufficient space for the CARv1 header at the beginning of the file.
//       Of course, the root CID is not known yet, so we leave a placeholder hash value and fill it in on completion.
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
//       |
//       Open question: Do we want to use CARv2 indexes? Probably yes.
//       Will complicate our bespoke CAR writing approach though and make it less maintainable.
//       Maybe we can construct indexes in-memory using go-car while we are writing CARv1 and then append then once done.

// TODO: there is a number of things [above] that are conceptually incorrect -- @ribasushi

var flags = Cmd.Flags()

var (
	flagDBs = flags.StringArray("db", nil, "Path to RocksDB (can be specified multiple times)")
)

func init() {
	Cmd.Run = run
}

func run(c *cobra.Command, args []string) {
	outPath := args[0]
	epochStr := args[1]
	epoch, err := strconv.ParseUint(epochStr, 10, 32)
	if err != nil {
		klog.Exitf("Invalid epoch arg: %s", epochStr)
	}

	// TODO mainnet history later on requires multiple CAR files per epoch
	f, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		klog.Exit(err)
	}
	defer f.Close()

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

	// Seek to epoch start and make sure we have all data
	const epochLen = 432000
	start := epoch * epochLen
	stop := start + epochLen
	mw.seek(start)
	klog.Infof("Starting at slot %d", start)
	slotsAvailable := mw.len()
	if slotsAvailable < epochLen {
		klog.Exitf("Need %d slots but got %d", epochLen, slotsAvailable)
	}

	ctx := c.Context()

	for ctx.Err() == nil {
		meta, ok := mw.next()
		if !ok {
			break
		}
		if meta.Slot > stop {
			break
		}
		entries, err := mw.get(meta)
		if err != nil {
			klog.Exitf("FATAL: Failed to get entry at slot %d: %s", meta.Slot, err)
		}
		_ = entries
		panic("unimplemented")
	}
}
