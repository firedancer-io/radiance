package dump

import (
	"errors"
	"fmt"
	"io"
	"os"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/ipld/go-car"
	"github.com/spf13/cobra"
	"go.firedancer.io/radiance/pkg/ipld/ipldgen"
	"k8s.io/klog/v2"
)

var Cmd = cobra.Command{
	Use:   "dump <car>",
	Short: "Dump the contents of a CAR file",
	Args:  cobra.ExactArgs(1),
}

func init() {
	Cmd.Run = run
}

func run(_ *cobra.Command, args []string) {
	f, err := os.Open(args[0])
	if err != nil {
		klog.Exit(err.Error())
	}
	defer f.Close()

	rd, err := car.NewCarReader(f)
	if err != nil {
		klog.Exitf("Failed to open CAR: %s", err)
	}
	for {
		block, err := rd.Next()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			klog.Exitf("Failed to read block: %s", err)
		}
		if block.Cid().Type() == ipldgen.SolanaTx {
			var tx solana.Transaction
			if err := bin.UnmarshalBin(&tx, block.RawData()); err != nil {
				klog.Errorf("Invalid CID %s: %s", block.Cid(), err)
				continue
			} else if len(tx.Signatures) == 0 {
				klog.Errorf("Invalid CID %s: tx has zero signatures", block.Cid())
				continue
			}
			fmt.Println(tx.Signatures[0].String())
		}
	}
}
