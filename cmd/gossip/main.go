// Gossip interacts with Solana gossip networks
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"github.com/certusone/radiance/cmd/gossip/ping"
	"github.com/certusone/radiance/cmd/gossip/pull"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var cmd = cobra.Command{
	Use:   "gossip",
	Short: "Interact with Solana gossip networks",
}

func init() {
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	cmd.PersistentFlags().AddGoFlagSet(klogFlags)

	cmd.AddCommand(
		&ping.Cmd,
		&pull.Cmd,
	)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cobra.CheckErr(cmd.ExecuteContext(ctx))
}
