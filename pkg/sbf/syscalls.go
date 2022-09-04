package sbf

import (
	"encoding/binary"

	"github.com/spaolacci/murmur3"
)

const (
	// EntrypointHash equals SymbolHash("entrypoint")
	EntrypointHash = uint32(0x71e3cf81)
)

// SymbolHash returns the murmur3 32-bit hash of a symbol name.
func SymbolHash(s string) uint32 {
	return murmur3.Sum32([]byte(s))
}

// PCHash returns the murmur3 32-bit hash of a program counter.
//
// Used by VM for non-syscall functions
func PCHash(addr uint64) uint32 {
	// TODO this is kinda pointless …
	var key [8]byte
	binary.LittleEndian.PutUint64(key[:], addr)
	return murmur3.Sum32(key[:])
}

type SyscallRegistry map[uint32]Syscall

func NewSyscallRegistry() SyscallRegistry {
	return make(SyscallRegistry)
}

func (s SyscallRegistry) Register(name string, syscall Syscall) (hash uint32, ok bool) {
	hash = SymbolHash(name)
	if _, exist := s[hash]; exist {
		return 0, false // collision or duplicate
	}
	s[hash] = syscall
	ok = true
	return
}
