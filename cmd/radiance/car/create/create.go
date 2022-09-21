package create

import (
	"go.firedancer.io/radiance/pkg/blockstore"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "create <rocksdb...>",
	Short: "Create CAR from RocksDB",
	Args:  cobra.MaximumNArgs(1),
}

func run(_ *cobra.Command, args []string) {
	rocksDB := args[0]

	db, err := blockstore.OpenReadOnly(rocksDB)
	if err != nil {
		klog.Exitf("Failed to open blockstore: %s", err)
	}
	defer db.Close()
}
