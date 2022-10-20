package poh

import (
	"crypto/rand"
	"crypto/sha256"
	"testing"

	sha256_simd "github.com/minio/sha256-simd"
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
