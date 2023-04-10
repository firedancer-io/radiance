package tpu_udp

import (
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/cmd/radiance/tpu_udp/pcap"
	"go.firedancer.io/radiance/cmd/radiance/tpu_udp/proxy"
	"go.firedancer.io/radiance/cmd/radiance/tpu_udp/sniff"
)

var Cmd = cobra.Command{
	Use:   "tpu-udp",
	Short: "TPU/UDP tools",
}

func init() {
	Cmd.AddCommand(
		&pcap.Cmd,
		&proxy.Cmd,
		&sniff.Cmd,
	)
}
