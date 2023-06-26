package solana

import "go.firedancer.io/radiance/pkg/base58"

type Hash [32]byte
type Address [32]byte
type Signature [64]byte

func (p *Hash) String() string {
	return base58.Encode(p[:])
}

func (p *Hash) UnmarshalText(b []byte) error {
	if !base58.Decode32((*[32]byte)(p), b) {
		return base58.ErrEncode
	}
	return nil
}

func (p *Address) String() string {
	return (*Hash)(p).String()
}

func (p *Address) UnmarshalText(b []byte) error {
	return (*Hash)(p).UnmarshalText(b)
}
