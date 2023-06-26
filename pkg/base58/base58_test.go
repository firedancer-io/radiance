package base58

import (
	"encoding/hex"
	"testing"
)

var testVector32 = []struct {
	hex string
	b58 string
}{
	{
		hex: "0000000000000000000000000000000000000000000000000000000000000000",
		b58: "11111111111111111111111111111111",
	},
	{
		hex: "0000000000000000000000000000000000000000000000000000000000000001",
		b58: "11111111111111111111111111111112",
	},
	{
		hex: "0000000000000000000000000000000000000000000000000000000000000101",
		b58: "1111111111111111111111111111115S",
	},
	{
		hex: "ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff",
		b58: "JEKNVnkbo3jma5nREBBJCDoXFVeKkD56V3xKrvRmWxFG",
	},
	{
		hex: "fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffe",
		b58: "JEKNVnkbo3jma5nREBBJCDoXFVeKkD56V3xKrvRmWxFF",
	},
}

func TestEncode32(t *testing.T) {
	for _, test := range testVector32 {
		var in [32]byte
		hex.Decode(in[:], []byte(test.hex))

		var out [44]byte
		outLen := Encode32(&out, in)

		outStr := string(out[:outLen])
		if outStr != test.b58 {
			t.Errorf("Encode32(%s) = %s, want %s", test.hex, outStr, test.b58)
		}
	}
}

func TestDecode32(t *testing.T) {
	for _, test := range testVector32 {
		var out [32]byte
		if !Decode32(&out, []byte(test.b58)) {
			t.Errorf("Decode32(%s) failed", test.b58)
			continue
		}

		outStr := hex.EncodeToString(out[:])
		if outStr != test.hex {
			t.Errorf("Decode32(%s) = %s, want %s", test.b58, outStr, test.hex)
		}
	}
}
