package genesis

import (
	"fmt"
	"io"

	"github.com/certusone/radiance/pkg/archiveutil"
	bin "github.com/gagliardetto/binary"
)

// ReadGenesisFromArchive reads a `genesis.tar.bz2` file.
func ReadGenesisFromArchive(archive io.Reader) (*Genesis, error) {
	tar, err := archiveutil.OpenTar(archive)
	if err != nil {
		return nil, err
	}
	hdr, err := tar.Next()
	if err != nil {
		return nil, err
	}
	if hdr.Name != "genesis.bin" {
		return nil, fmt.Errorf("first file is not genesis.bin")
	}
	const maxSize = 10_000_001
	rd := io.LimitReader(tar, maxSize)
	genesisBytes, err := io.ReadAll(rd)
	if err != nil {
		return nil, err
	}
	if len(genesisBytes) == maxSize {
		return nil, fmt.Errorf("genesis.bin too large")
	}
	genesis := new(Genesis)
	dec := bin.NewBinDecoder(genesisBytes)
	err = dec.Decode(genesis)
	if err == nil && dec.HasRemaining() {
		err = fmt.Errorf("not all of genesis.bin was read (%d bytes remaining)", dec.Remaining())
	}
	return genesis, err
}
