package car

import (
	"go.firedancer.io/radiance/cmd/radiance/car/create"
	"github.com/spf13/cobra"
)

var Cmd = cobra.Command{
	Use:   "car",
	Short: "Manage IPLD Content-addressable ARchives",
	Long:  "https://ipld.io/specs/transport/car/",
}

func init() {
	Cmd.AddCommand(
		&create.Cmd,
	)
}
