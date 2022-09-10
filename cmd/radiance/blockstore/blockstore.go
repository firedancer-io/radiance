package blockstore

import (
	"github.com/spf13/cobra"
)

var Cmd = cobra.Command{
	Use:   "blockstore",
	Short: "Access blockstore database",
}

func init() {
	Cmd.AddCommand()
}
