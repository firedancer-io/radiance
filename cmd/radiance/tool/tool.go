package tool

import "github.com/spf13/cobra"

var Cmd = cobra.Command{
	Use:   "tool",
	Short: "Random tools",
}

func init() {
	Cmd.AddCommand()
}
