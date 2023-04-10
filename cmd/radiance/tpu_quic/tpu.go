package tpu_quic

import (
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/cmd/radiance/tpu_quic/ping"
)

var Cmd = cobra.Command{
	Use:   "tpu-quic",
	Short: "TPU/QUIC tools",
}

func init() {
	Cmd.AddCommand(
		&ping.Cmd,
	)
}
