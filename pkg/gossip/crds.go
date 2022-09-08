package gossip

import (
	"crypto/ed25519"
	"encoding/binary"
	"math"
)

// CrdsBloomP is bloom filter 'p' parameter (probability)
const CrdsBloomP = 0.1

func (f *CrdsFilter) Contains(item *Hash) bool {
	if !f.TestMask(item) {
		return false
	}
	return f.Filter.Contains(item)
}

func (f *CrdsFilter) TestMask(item *Hash) bool {
	ones := uint64(math.MaxUint64) >> f.MaskBits
	bits := binary.LittleEndian.Uint64(item[:8]) | ones
	return bits == f.Mask
}

type CrdsFilterSet []CrdsFilter

func NewCrdsFilterSet(numItems, maxBytes uint64) CrdsFilterSet {
	maxBits := maxBytes * 8
	maxItems := BloomMaxItems(float64(maxBits), CrdsBloomP, 8)
	maskBits := BloomMaskBits(float64(numItems), maxItems)
	filters := make([]CrdsFilter, 1<<maskBits)
	for i := uint64(0); i < uint64(len(filters)); i++ {
		bloom := NewBloomRandom(uint64(maxItems), CrdsBloomP, maxBits)
		seed := i << (64 - maskBits)
		filters[i] = CrdsFilter{
			Filter:   *bloom,
			Mask:     seed | (math.MaxUint64 >> maskBits),
			MaskBits: maskBits,
		}
	}
	return filters
}

func (c CrdsFilterSet) Add(h Hash) {
	index := binary.LittleEndian.Uint64(h[:8])
	index >>= 64 - c[0].MaskBits
	c[index].Filter.Add(&h)
}

func (c *CrdsValue) Sign(identity ed25519.PrivateKey) error {
	// Write pubkey into data field
	pubkey := ed25519.PublicKey(c.Data.Pubkey()[:])
	copy(
		pubkey[:ed25519.PublicKeySize],
		identity.Public().(ed25519.PublicKey)[:ed25519.PublicKeySize],
	)

	msg, err := c.Data.BincodeSerialize()
	if err != nil {
		return err
	}
	sig := ed25519.Sign(identity, msg)
	copy(c.Signature[:ed25519.SignatureSize], sig[:ed25519.SignatureSize])
	return nil
}

func (c *CrdsValue) VerifySignature() bool {
	msg, err := c.Data.BincodeSerialize()
	if err != nil {
		return false
	}
	pubkey := ed25519.PublicKey(c.Data.Pubkey()[:])
	return ed25519.Verify(pubkey, msg, c.Signature[:])
}
