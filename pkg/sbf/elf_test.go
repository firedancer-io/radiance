package sbf

import (
	"debug/elf"
	_ "embed"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	//go:embed tests/noop.so
	soNoop []byte
)

func TestLoadProgram_Noop(t *testing.T) {
	exe, err := LoadProgram(soNoop)
	require.NoError(t, err)
	assert.Equal(t, &Executable{
		Header: elf.Header64{
			Ident: [16]byte{
				0x7f, 0x45, 0x4c, 0x46,
				0x02, 0x01, 0x01, 0x00,
				0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
			},
			Type:      uint16(elf.ET_DYN),
			Machine:   uint16(elf.EM_BPF),
			Version:   uint32(elf.EV_CURRENT),
			Entry:     4096,
			Phoff:     64,
			Shoff:     8792,
			Flags:     0,
			Ehsize:    64,
			Phentsize: 56,
			Phnum:     7,
			Shentsize: 64,
			Shnum:     13,
			Shstrndx:  11,
		},
		Load: elf.Prog64{
			Type:   uint32(elf.PT_LOAD),
			Flags:  06,
			Off:    8192,
			Vaddr:  8192,
			Paddr:  8192,
			Filesz: 208,
			Memsz:  208,
			Align:  4096,
		},
		ShShstrtab: elf.Section64{
			Name:      82,
			Type:      uint32(elf.SHT_STRTAB),
			Flags:     0,
			Addr:      0,
			Off:       8648,
			Size:      100,
			Addralign: 1,
		},
		ShSymtab: &elf.Section64{
			Name:      74,
			Type:      uint32(elf.SHT_SYMTAB),
			Flags:     0,
			Addr:      0,
			Off:       8504,
			Size:      144,
			Link:      12,
			Info:      3,
			Addralign: 8,
			Entsize:   24,
		},
		ShStrtab: &elf.Section64{
			Name:      92,
			Type:      uint32(elf.SHT_STRTAB),
			Flags:     0,
			Addr:      0,
			Off:       8748,
			Size:      39,
			Link:      0,
			Info:      0,
			Addralign: 1,
			Entsize:   0,
		},
		ShDynstr: &elf.Section64{
			Name:      25,
			Type:      uint32(elf.SHT_STRTAB),
			Flags:     uint64(elf.DF_SYMBOLIC),
			Addr:      624,
			Off:       624,
			Size:      23,
			Link:      0,
			Info:      0,
			Addralign: 1,
			Entsize:   0,
		},
	}, exe)
}
