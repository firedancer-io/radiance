package create

import (
	"os"
	"strconv"

	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/pkg/ipld/car"
	"go.firedancer.io/radiance/pkg/ipld/ipldgen"
	"k8s.io/klog/v2"

	"go.firedancer.io/radiance/pkg/blockstore"
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
	if !mw.seek(start) {
		klog.Exitf("Slot %d not available in any DB", start)
	}
	// TODO: This is not robust; if the DB starts in the middle of the epoch, the first slots are going to be skipped.
	klog.Infof("Starting at slot %d", start)
	slotsAvailable := mw.len()
	if slotsAvailable < epochLen {
		klog.Exitf("Need slots [%d:%d] (epoch %d) but only have up to %d",
			start, stop, epoch, start+slotsAvailable)
	}

	// TODO mainnet history later on requires multiple CAR files per epoch
	f, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0666)
	if err != nil {
		klog.Exit(err)
	}
	defer f.Close()

	carOut, err := car.NewWriter(f)
	if err != nil {
		klog.Exitf("Cannot create CARv1 file: %s", err)
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

		asm := ipldgen.NewBlockAssembler(carOut, meta.Slot)

		entryNum := 0
		klog.V(3).Infof("Slot %d", meta.Slot)
		for i, batch := range entries {
			klog.V(6).Infof("Slot %d batch %d", meta.Slot, i)

			for j, entry := range batch.Entries {
				pos := ipldgen.EntryPos{
					Slot:       meta.Slot,
					EntryIndex: entryNum,
					Batch:      i,
					BatchIndex: j,
					LastShred:  -1,
				}
				if j == len(batch.Entries)-1 {
					// We map "last shred of batch" to each "last entry of batch"
					// so we can reconstruct the shred/entry-batch assignments.
					pos.LastShred = int(batch.Shreds[len(batch.Shreds)-1].CommonHeader().Index)
				}

				if err := asm.WriteEntry(entry, pos); err != nil {
					klog.Exitf("Failed to write slot %d shred %d (batch %d index %d): %s",
						meta.Slot, entryNum, i, j, err)
				}

				entryNum++
			}
		}

		// TODO roll up into ledger entries
		if _, err := asm.Finish(); err != nil {
			klog.Exitf("Failed to write block: %s", err)
		}
	}
}
