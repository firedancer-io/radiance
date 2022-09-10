package yaml

import (
	"fmt"
	"os"
	"strings"

	"github.com/certusone/radiance/cmd/radiance/blockstore/util"
	"github.com/certusone/radiance/pkg/blockstore"
	"github.com/linxGnu/grocksdb"
	"github.com/segmentio/textio"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "yaml <rocksdb>",
	Short: "Dump blockstore content to YAML",
	Args:  cobra.ExactArgs(1),
}

var flags = Cmd.Flags()

var (
	flagSlots   = flags.String("slots", "", "Slots to dump")
	flagEntries = flags.Bool("entries", false, "Also dump slot entries")
	flagShreds  = flags.Bool("shreds", false, "Also dump shreds")
)

func init() {
	Cmd.Run = run
}

func run(_ *cobra.Command, args []string) {
	rocksDB := args[0]

	printColumnFamilies(rocksDB)

	db, err := blockstore.OpenReadOnly(rocksDB)
	if err != nil {
		klog.Exitf("Failed to open blockstore: %s", err)
	}
	defer db.Close()

	printRoot(db)

	slots, ok := util.ParseInts(*flagSlots)
	if !ok {
		klog.Exitf("Invalid slots specifier: %s", *flagSlots)
	}
	if len(slots) > 0 {
		dumpSlots(db, slots)
	}

	if klog.Stats.Error.Lines() > 0 {
		os.Exit(1)
	}
}

func printColumnFamilies(dbPath string) {
	dbOpts := grocksdb.NewDefaultOptions()
	names, err := grocksdb.ListColumnFamilies(dbOpts, dbPath)
	if err != nil {
		klog.Error("Failed to list column families: %s", err)
		return
	}
	fmt.Println("column_families:")
	for _, name := range names {
		fmt.Println("  - " + name)
	}
}

func printRoot(db *blockstore.DB) {
	root, err := db.MaxRoot()
	if err != nil {
		klog.Error("Failed to get root: ", err)
		return
	}
	fmt.Println("root:", root)
}

func dumpSlots(db *blockstore.DB, slots util.Ints) {
	fmt.Println("slots:")
	slots.Iter(func(slot uint64) bool {
		dumpSlot(db, slot)
		return true
	})
}

func dumpSlot(db *blockstore.DB, slot uint64) {
	slotMeta, err := db.GetSlotMeta(slot)
	if err != nil {
		klog.Errorf("Failed to get slot %d: %s", slot, err)
		return
	}

	fmt.Printf("  %d:\n", slot)
	printSlotMeta(slotMeta)
	if *flagShreds {
		dumpDataShreds(db, slot)
	}
	if *flagEntries {
		dumpDataEntries(db, slotMeta)
	}
}

func printSlotMeta(slotMeta *blockstore.SlotMeta) {
	enc := newYAMLPrinter(2)
	defer enc.Close()
	if err := enc.Encode(slotMeta); err != nil {
		panic(err.Error())
	}
}

func dumpDataShreds(db *blockstore.DB, slot uint64) {
	shreds, err := db.GetAllDataShreds(slot)
	if err != nil {
		klog.Errorf("Failed to get data shreds of slot %d: %s", slot, err)
		return
	}

	fmt.Println("    data_shreds:")

	enc := newYAMLPrinter(3)
	defer enc.Close()
	if err := enc.Encode(shreds); err != nil {
		panic(err.Error())
	}
}

func dumpDataEntries(db *blockstore.DB, meta *blockstore.SlotMeta) {
	entries, err := db.GetEntries(meta)
	if err != nil {
		klog.Errorf("Failed to recover entries of slot %d: %s", meta.Slot, err)
		return
	}

	fmt.Println("    entries:")

	enc := newYAMLPrinter(3)
	defer enc.Close()
	if err := enc.Encode(entries); err != nil {
		panic(err.Error())
	}
}

func newYAMLPrinter(level int) *yaml.Encoder {
	enc := yaml.NewEncoder(textio.NewPrefixWriter(os.Stdout, strings.Repeat("  ", level)))
	enc.SetIndent(2)
	return enc
}
