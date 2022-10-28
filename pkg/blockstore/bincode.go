package blockstore

import (
	bin "github.com/gagliardetto/binary"
)

func ParseBincode[T any](data []byte) (*T, error) {
	dec := bin.NewBinDecoder(data)
	val := new(T)
	err := dec.Decode(val)
	return val, err
}
