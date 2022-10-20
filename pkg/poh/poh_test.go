package poh

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"testing"

	sha256_simd "github.com/minio/sha256-simd"
	"github.com/stretchr/testify/assert"
)

func BenchmarkHashchain_Stdlib(b *testing.B) {
	var state [32]byte
	rand.Read(state[:])
	for i := 0; i < b.N; i++ {
		state = sha256.Sum256(state[:])
	}
}

func BenchmarkHashchain_MinioSimd(b *testing.B) {
	var state [32]byte
	rand.Read(state[:])
	for i := 0; i < b.N; i++ {
		state = sha256_simd.Sum256(state[:])
	}
}

func TestState(t *testing.T) {
	type step struct {
		append uint
		mixin  string
	}
	cases := map[string]struct {
		pre   string
		post  string
		steps []step
	}{
		"None": {
			pre:   "0000000000000000000000000000000000000000000000000000000000000000",
			post:  "0000000000000000000000000000000000000000000000000000000000000000",
			steps: nil,
		},
		"Mainnet_0": {
			pre:  "45296998a6f8e2a784db5d9f95e18fc23f70441a1039446801089879b08c7ef0",
			post: "3973e330c29b831f3fcb0e49374ed8d0388f410a23e4ebf23328505036efbd03",
			steps: []step{
				{append: 800000},
			},
		},
		"Mainnet_1": {
			pre:  "3973e330c29b831f3fcb0e49374ed8d0388f410a23e4ebf23328505036efbd03",
			post: "8ee20607dcf1d9393cf5a2f2c9f7babe167dbdd267491b513c73d2cbf87413f5",
			steps: []step{
				{append: 14612},
				{mixin: "c95f2f13a9a77f32b1437976c4cffe3029298a49bf37007f8e45d793a520f30b"},
				{append: 210347},
				{mixin: "1aaeeb36611f484d984683a3db9269f2292dd9bb81bdab82b28c45625d9abd59"},
				{append: 428775},
				{mixin: "db31e861b310f44954403e345b6beeb3ded34084b90694bccaa2345306d366e1"},
				{append: 146263},
			},
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			state := State(mustHashFromHex(tc.pre))
			for _, step := range tc.steps {
				if step.append != 0 {
					state.Hash(step.append)
				} else if step.mixin != "" {
					mixin := mustHashFromHex(step.mixin)
					state.Record(&mixin)
				}
			}
			post := State(mustHashFromHex(tc.post))
			assert.Equal(t, post, state)
		})
	}
}

func mustHashFromHex(s string) (out [32]byte) {
	n, err := hex.Decode(out[:], []byte(s))
	if err != nil {
		panic(err.Error())
	} else if n != len(out) {
		panic("odd hex len")
	}
	return
}
