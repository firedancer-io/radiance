//go:build rocksdb

package blockstore

import (
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/dumpshreds"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/verifydata"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/yaml"
)

var Cmd = cobra.Command{
	Use:   "blockstore",
	Short: "Access blockstore database",
}

func init() {
	Cmd.AddCommand(
		&dumpshreds.Cmd,
		&verifydata.Cmd,
		&yaml.Cmd,
	)
}
