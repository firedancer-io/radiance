package blockstore

import (
	"github.com/certusone/radiance/cmd/radiance/blockstore/dumpshreds"
	"github.com/certusone/radiance/cmd/radiance/blockstore/verifydata"
	"github.com/certusone/radiance/cmd/radiance/blockstore/yaml"
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
