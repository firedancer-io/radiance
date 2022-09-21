package gossip

import (
	"go.firedancer.io/radiance/cmd/radiance/gossip/ping"
	"go.firedancer.io/radiance/cmd/radiance/gossip/pull"
	"github.com/spf13/cobra"
)

var Cmd = cobra.Command{
	Use:   "gossip",
	Short: "Interact with Solana gossip networks",
}

func init() {
	Cmd.AddCommand(
		&ping.Cmd,
		&pull.Cmd,
	)
}
