package bigtable

import (
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/cmd/radiance/bigtable/dumpblock"
)

var Cmd = cobra.Command{
	Use:   "bigtable",
	Short: "Google Bigtable tools",
}

func init() {
	Cmd.AddCommand(
		&dumpblock.Cmd,
	)
}
