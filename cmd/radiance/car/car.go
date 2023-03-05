//go:build !lite

package car

import (
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/cmd/radiance/car/create"
	"go.firedancer.io/radiance/cmd/radiance/car/dump"
)

var Cmd = cobra.Command{
	Use:   "car",
	Short: "Manage IPLD Content-addressable ARchives",
	Long:  "https://ipld.io/specs/transport/car/",
}

func init() {
	Cmd.AddCommand(
		&create.Cmd,
		&dump.Cmd,
	)
}
