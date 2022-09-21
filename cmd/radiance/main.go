// Gossip interacts with Solana gossip networks
package main

import (
	"context"
	"flag"
	"os"
	"os/signal"

	"go.firedancer.io/radiance/cmd/radiance/blockstore"
	"go.firedancer.io/radiance/cmd/radiance/gossip"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"
)

var cmd = cobra.Command{
	Use:   "radiance",
	Short: "Solana Go playground",
}

func init() {
	klogFlags := flag.NewFlagSet("klog", flag.ExitOnError)
	klog.InitFlags(klogFlags)
	cmd.PersistentFlags().AddGoFlagSet(klogFlags)

	cmd.AddCommand(
		&blockstore.Cmd,
		&gossip.Cmd,
	)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	cobra.CheckErr(cmd.ExecuteContext(ctx))
}
