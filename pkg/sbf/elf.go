package sbf

import (
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/bits"
)

// TODO Fuzz
// TODO Differential fuzz against rbpf

type Executable struct {
	Header elf.Header64
	Load   elf.Prog64
	ShStr  elf.Section64
}

func LoadProgram(buf []byte) (*Executable, error) {
	l := loader{
		rd:       bytes.NewReader(buf),
		fileSize: uint64(len(buf)),
	}
	return l.load()
}

type loader struct {
	rd       io.ReaderAt
	fileSize uint64
	elf      *Executable
}

func (l *loader) load() (*Executable, error) {
	l.elf = new(Executable)
	if err := l.readHeader(); err != nil {
		return nil, err
	}
	if err := l.validateHeader(); err != nil {
		return nil, err
	}
	if err := l.loadProgramHeaderTable(); err != nil {
		return nil, err
	}
	if err := l.validateSectionHeaderTable(); err != nil {
		return nil, err
	}
	// TODO load section name section header
	// TODO parse sections
	// TODO parse dynamic segment
	return l.elf, nil
}

const (
	ehsize    = 0x40
	phentsize = 0x38
	shentsize = 0x40
)

func (l *loader) readHeader() error {
	var hdrBuf [ehsize]byte
	if _, err := io.ReadFull(io.NewSectionReader(l.rd, 0, ehsize), hdrBuf[:]); err != nil {
		return err
	}
	return binary.Read(bytes.NewReader(hdrBuf[:]), binary.LittleEndian, &l.elf.Header)
}

func (l *loader) validateHeader() error {
	eh := &l.elf.Header
	ident := &eh.Ident

	if string(ident[:elf.EI_CLASS]) != elf.ELFMAG {
		return fmt.Errorf("not an ELF file")
	}

	if elf.Class(ident[elf.EI_CLASS]) != elf.ELFCLASS64 ||
		elf.Data(ident[elf.EI_DATA]) != elf.ELFDATA2LSB ||
		elf.Version(ident[elf.EI_VERSION]) != elf.EV_CURRENT ||
		elf.OSABI(ident[elf.EI_OSABI]) != elf.ELFOSABI_NONE {
		return fmt.Errorf("incompatible binary")
	}
	// note: EI_PAD and EI_ABIVERSION are ignored

	if elf.Machine(eh.Machine) != elf.EM_BPF ||
		elf.Type(eh.Type) != elf.ET_DYN ||
		eh.Version != uint32(elf.EV_CURRENT) ||
		eh.Ehsize != ehsize ||
		eh.Phentsize != phentsize ||
		eh.Shentsize != shentsize ||
		eh.Shstrndx >= eh.Shnum {
		return fmt.Errorf("invalid ELF file")
	}

	if eh.Phoff < ehsize {
		return fmt.Errorf("program header overlaps with file header")
	}
	if eh.Shoff < ehsize {
		return fmt.Errorf("section header overlaps with file header")
	}
	if isOverlap(eh.Phoff, uint64(eh.Phnum)*phentsize, eh.Shoff, uint64(eh.Shnum)*shentsize) {
		return fmt.Errorf("program and section header overlap")
	}

	return nil
}

// scan the program header table and remember the last PT_LOAD segment
func (l *loader) loadProgramHeaderTable() error {
	eh := &l.elf.Header
	phoff := eh.Phoff

	for i := uint16(0); i < eh.Phnum; i++ {
		if phoff+phentsize > math.MaxInt64 {
			return io.ErrUnexpectedEOF
		}
		rd := io.NewSectionReader(l.rd, int64(phoff), phentsize)

		var ph elf.Prog64
		if err := binary.Read(rd, binary.LittleEndian, &ph); err != nil {
			return err
		}
		phoff += phentsize

		if elf.ProgType(ph.Type) != elf.PT_LOAD {
			continue
		}

		// vaddr must be ascending
		if ph.Vaddr < l.elf.Load.Vaddr {
			return fmt.Errorf("invalid program header")
		}

		segmentEnd, overflow := bits.Add64(ph.Off, ph.Filesz, 0)
		if segmentEnd > l.fileSize || overflow > 0 {
			return fmt.Errorf("segment out of bounds")
		}

		l.elf.Load = ph
	}

	return nil
}

func (l *loader) validateSectionHeaderTable() error {
	eh := &l.elf.Header
	shoff := eh.Shoff

	offset := uint64(0)
	for i := uint16(0); i < eh.Shnum; i++ {
		if shoff+shentsize > math.MaxInt64 {
			return io.ErrUnexpectedEOF
		}
		rd := io.NewSectionReader(l.rd, int64(shoff), shentsize)

		var sh elf.Section64
		if err := binary.Read(rd, binary.LittleEndian, &sh); err != nil {
			return err
		}
		shoff += shentsize

		if i == 0 {
			if elf.SectionType(sh.Type) != elf.SHT_NULL {
				return fmt.Errorf("section 0 is not SHT_NULL")
			}
			continue
		}
		if elf.SectionType(sh.Type) == elf.SHT_NOBITS {
			continue
		}

		shend, overflow := bits.Add64(sh.Off, sh.Size, 0)
		if overflow != 0 {
			return fmt.Errorf("integer overflow in section %d", i)
		}
		if sh.Off < ehsize {
			return fmt.Errorf("section %d overlaps with file header", i)
		}
		if isOverlap(eh.Phoff, uint64(eh.Phnum)*phentsize, sh.Off, sh.Size) {
			return fmt.Errorf("section %d overlaps with program header", i)
		}
		if isOverlap(eh.Shoff, uint64(eh.Shnum)*shentsize, sh.Off, sh.Size) {
			return fmt.Errorf("section %d overlaps with section header", i)
		}
		if eh.Shoff < offset {
			return fmt.Errorf("sections not in order")
		}
		if shend > l.fileSize {
			return fmt.Errorf("section %d out of bounds", i)
		}
		offset = shend
	}

	return nil
}

func isOverlap(startA uint64, sizeA uint64, startB uint64, sizeB uint64) bool {
	if startA > startB {
		startA, sizeA, startB, sizeB = startB, sizeB, startA, sizeA
	}
	endA, endB := startA+sizeA, startB+sizeB
	if endA < startA || endB < startB {
		panic("isOverlap: integer overflow")
	}
	return sizeA != 0 && sizeB != 0 && (startA == startB || endA > endB)
}
