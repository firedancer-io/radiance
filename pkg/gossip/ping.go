package gossip

import (
	"crypto/ed25519"
	"crypto/sha256"
)

// PingSize is the size of a serialized ping message.
const PingSize = 128

// NewPing creates and signs a new ping message.
//
// Panics if the provided private key is invalid.
func NewPing(token [32]byte, key ed25519.PrivateKey) (p Ping) {
	sig := ed25519.Sign(key, token[:])
	copy(p.From[:], key.Public().(ed25519.PublicKey))
	copy(p.Token[:], token[:])
	copy(p.Signature[:], sig[:])
	return p
}

// Verify checks the Ping's signature.
func (p *Ping) Verify() bool {
	return ed25519.Verify(p.From[:], p.Token[:], p.Signature[:])
}

// HashPingToken returns the pong token given a ping token.
func HashPingToken(token [32]byte) [32]byte {
	return sha256.Sum256(token[:])
}
