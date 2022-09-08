package gossip

import "github.com/gagliardetto/solana-go"

func (p Pubkey) MarshalText() ([]byte, error) {
	return solana.PublicKey(p).MarshalText()
}

func (h Hash) MarshalText() ([]byte, error) {
	return []byte(solana.Hash(h).String()), nil
}

func (s Signature) MarshalText() ([]byte, error) {
	return []byte(solana.Signature(s).String()), nil
}
