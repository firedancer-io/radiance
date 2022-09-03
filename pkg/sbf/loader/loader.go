package loader

import (
	"bytes"
	"debug/elf"
	"fmt"
	"io"

	"github.com/certusone/radiance/pkg/sbf"
)

// Loader is based on solana_rbpf::elf_parser
type Loader struct {
	// File containing ELF
	rd       io.ReaderAt
	fileSize uint64

	// ELF data structures
	eh         elf.Header64
	phLoad     elf.Prog64
	phDynamic  *elf.Prog64
	shShstrtab elf.Section64
	shText     *elf.Section64
	shSymtab   *elf.Section64
	shStrtab   *elf.Section64
	shDynstr   *elf.Section64
	shDynamic  *elf.Section64
	dynamic    [DT_NUM]uint64
	relocsIter *tableIter[elf.Rel64]
	dynSymIter *tableIter[elf.Sym64]

	// Program section/segment mappings
	// Uses physical addressing
	rodatas   []addrRange
	text      addrRange
	progRange addrRange

	// Contains most of ELF (.text and rodata-like)
	// Non-loaded sections are zeroed
	program []byte

	// Symbols
	//funcs    map[uint32]symbol
	//syscalls map[uint32]string
}

// Bounds checks
const (
	// 64 MiB max program size.
	// Allows loader to use unchecked math when adding 32-bit offsets.
	maxFileLen = 1 << 26

	maxSectionNameLen = 16
	maxSymbolNameLen  = 1024
)

// EF_SBF_V2 is the SBFv2 ELF flag
const EF_SBF_V2 = 0x20

// DT_NUM is the number of ELF generic dynamic entry types
const DT_NUM = 35

// Hardcoded addresses.
const (
	VaddrProgram = uint64(0x1_0000_0000)
	VaddrStack   = uint64(0x2_0000_0000)
	VaddrHeap    = uint64(0x3_0000_0000)
	VaddrInput   = uint64(0x4_0000_0000)
)

// NewLoaderFromBytes creates an ELF loader from a byte slice.
func NewLoaderFromBytes(buf []byte) (*Loader, error) {
	if len(buf) > maxFileLen {
		return nil, fmt.Errorf("ELF file too large")
	}
	l := &Loader{
		rd:       bytes.NewReader(buf),
		fileSize: uint64(len(buf)),
	}
	return l, nil
}

// Load parses, loads, and relocates an SBF program.
func (l *Loader) Load() (*sbf.Program, error) {
	if err := l.parse(); err != nil {
		return nil, err
	}
	if err := l.copy(); err != nil {
		return nil, err
	}
	//if err := l.relocate(); err != nil {
	//	return nil, err
	//}
	panic("unimplemented")
}
