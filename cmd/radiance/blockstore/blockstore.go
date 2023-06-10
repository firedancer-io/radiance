//go:build !lite

package blockstore

import (
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/compact"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/dumpbatches"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/dumpshreds"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/statdatarate"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/statentries"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/verifydata"
	"go.firedancer.io/radiance/cmd/radiance/blockstore/yaml"
)

var Cmd = cobra.Command{
	Use:   "blockstore",
	Short: "Access blockstore database",
}

func init() {
	Cmd.AddCommand(
		&compact.Cmd,
		&dumpshreds.Cmd,
		&dumpbatches.Cmd,
		&statdatarate.Cmd,
		&statentries.Cmd,
		&verifydata.Cmd,
		&yaml.Cmd,
	)
}
