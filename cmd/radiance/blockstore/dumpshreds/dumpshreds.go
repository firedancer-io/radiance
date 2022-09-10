package dumpshreds

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"

	"github.com/certusone/radiance/cmd/radiance/blockstore/util"
	"github.com/certusone/radiance/pkg/blockstore"
	"github.com/linxGnu/grocksdb"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "dump-shreds <rocksdb> <out> <slots>",
	Short: "Dump shreds to file system",
	Args:  cobra.ExactArgs(3),
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

	if err := os.Mkdir(outPath, 0755); err != nil && !errors.Is(err, fs.ErrExist) {
		klog.Exit(err)
	}

	slots, ok := util.ParseInts(args[2])
	if !ok {
		klog.Exit("Invalid slots parameter")
	}
	slots.Iter(func(slot uint64) bool {
		err := dumpSlot(db, outPath, slot)
		if err != nil {
			klog.Warning("Failed to dump slot %d: %s", slot, err)
		}
		return true
	})
}

func dumpSlot(db *blockstore.DB, outPath string, slot uint64) error {
	slotPath := filepath.Join(outPath, strconv.FormatUint(slot, 10))
	if err := os.Mkdir(slotPath, 0755); err != nil && !errors.Is(err, fs.ErrExist) {
		return err
	}

	if err := dumpShreds(db, slotPath, slot, db.CfCodeShred, "c"); err != nil {
		return err
	}
	if err := dumpShreds(db, slotPath, slot, db.CfDataShred, "d"); err != nil {
		return err
	}
	return nil
}

func dumpShreds(
	db *blockstore.DB, slotPath string, slot uint64,
	cf *grocksdb.ColumnFamilyHandle,
	namePrefix string,
) error {
	iter := db.DB.NewIteratorCF(grocksdb.NewDefaultReadOptions(), cf)
	defer iter.Close()
	prefix := blockstore.MakeShredKey(slot, 0)
	iter.Seek(prefix[:])
	for {
		curSlot, curIndex, ok := blockstore.ParseShredKey(iter.Key().Data())
		if !ok || curSlot != slot {
			break
		}
		p := filepath.Join(slotPath, fmt.Sprintf("%s%03d", namePrefix, curIndex))
		if err := os.WriteFile(p, iter.Value().Data(), 0644); err != nil {
			return err
		}
		iter.Next()
	}
	return nil
}
