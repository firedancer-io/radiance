//go:build !lite

package tarblocks

import (
	"archive/tar"
	"bufio"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/pkg/blockstore"
	"go.firedancer.io/radiance/pkg/shred"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "tar-blocks <rocksdb> <out>",
	Short: "Export rooted blocks from RocksDB to TAR stream",
	Long: `Creates a TAR stream of rooted blocks in serialized-batches binary format.

Each block/<slot>.bin file is a bincode Vec<Vec<Entry>> serialization of the block data.
Currently outputs only rooted slots in replayable order.
This may change in the future.

Use "-" as <out> to write the TAR stream to stdout.`,
	Args: cobra.ExactArgs(2),
}

func init() {
	Cmd.Run = run
}

func run(_ *cobra.Command, args []string) {
	rocksDB := args[0]
	outPath := args[1]

	db, err := blockstore.OpenReadOnly(rocksDB)
	if err != nil {
		klog.Exitf("Failed to open blockstore: %s", err)
	}
	defer db.Close()

	walker, err := blockstore.NewBlockWalk([]blockstore.WalkHandle{{DB: db}}, 2)
	if err != nil {
		klog.Fatal(err)
	}
	defer walker.Close()

	var rawOutStream io.Writer
	if outPath == "-" {
		if isatty.IsTerminal(os.Stdout.Fd()) {
			klog.Exit("Refusing to write binary data to terminal")
		}
		rawOutStream = os.Stdout
	} else {
		f, err := os.Create(outPath)
		if err != nil {
			klog.Exitf("Failed to create output file: %s", err)
		}
		defer f.Close()
		rawOutStream = f
	}

	outStream := bufio.NewWriter(rawOutStream)
	defer outStream.Flush()

	outTar := tar.NewWriter(outStream)
	defer outTar.Close()

	outTar.WriteHeader(&tar.Header{
		Typeflag: tar.TypeDir,
		Name:     "block/",
		Mode:     0755,
		ModTime:  time.Now(),
	})

	for {
		meta, ok := walker.Next()
		if !ok {
			break
		}

		shreds, err := walker.Current().GetAllDataShreds(meta.Slot, 2)
		if err != nil {
			klog.Warningf("Failed to get shreds for slot %d: %s", meta.Slot, err)
			break
		}

		block := shred.Concat(shreds)
		hdr := tar.Header{
			Typeflag: tar.TypeReg,
			Name:     fmt.Sprintf("block/%d.bin", meta.Slot),
			Size:     int64(len(block)),
			Mode:     0644,
			ModTime:  time.Now(),
		}
		if err := outTar.WriteHeader(&hdr); err != nil {
			klog.Warningf("Failed to write header for slot %d: %s", meta.Slot, err)
			break
		}
		if _, err := outTar.Write(block); err != nil {
			klog.Warningf("Failed to write block %d: %s", meta.Slot, err)
			break
		}

		klog.V(7).Infof("Dmuped %s", hdr.Name)
	}

	klog.Info("Done")
}
