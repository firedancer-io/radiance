package gossip

import "github.com/gagliardetto/solana-go"

func (p Pubkey) MarshalText() ([]byte, error) {
	return solana.PublicKey(p).MarshalText()
}

func (h Hash) MarshalText() ([]byte, error) {
	return solana.Hash(h).MarshalText()
}

func (s Signature) MarshalText() ([]byte, error) {
	return solana.Signature(s).MarshalText()
}
