package blockstore

import (
	"go.firedancer.io/radiance/cmd/radiance/blockstore/dumpshreds"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/verifydata"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/yaml"
	"github.com/spf13/cobra"
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
