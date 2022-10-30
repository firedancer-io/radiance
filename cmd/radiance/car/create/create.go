//go:build rocksdb

package create

import (
	"path/filepath"
	"strconv"
	"time"

	"github.com/pkg/profile"
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/pkg/ipld/cargen"
	"k8s.io/klog/v2"

	"go.firedancer.io/radiance/pkg/blockstore"
)

var Cmd = cobra.Command{
	Use:   "create <epoch>",
	Short: "Create CAR files from blockstore",
	Long: "Extracts Solana ledger data from blockstore (RocksDB) databases,\n" +
		"and outputs IPLD CARs (content-addressable archives).\n" +
		"\n" +
		"Produces at least one CAR per epoch.\n" +
		"CAR archive contents are deterministic.",
	Args: cobra.ExactArgs(1),
}

var flags = Cmd.Flags()

var (
	flagOut = flags.StringP("out", "o", "", "Output directory")
	flagDBs = flags.StringArray("db", nil, "Path to RocksDB (can be specified multiple times)")

	flagProfilePath = flags.String("profile-path", "", "Path to record profile to (empty to use temp dir)")
	flagProfileDur  = flags.Duration("profile-duration", 0, "Record Go pprof (0 to disable)")
)

func init() {
	Cmd.Run = run
}

func run(c *cobra.Command, args []string) {
	start := time.Now()

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

	// Invoke pprof profiler.
	if *flagProfileDur > 0 {
		go func() {
			p := profile.Start(profile.ProfilePath(*flagProfilePath))
			defer p.Stop()
			time.Sleep(*flagProfileDur)
		}()
	}

	ctx := c.Context()
	if err = w.Run(ctx); err != nil {
		klog.Exitf("FATAL: %s", err)
	}
	klog.Info("DONE")

	timeTaken := time.Since(start)
	klog.Infof("Time taken: %s", timeTaken)
}
