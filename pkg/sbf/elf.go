package sbf

import (
	"bufio"
	"bytes"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"math/bits"
	"strings"
)

// TODO Fuzz
// TODO Differential fuzz against rbpf

type Executable struct {
	Header     elf.Header64
	Load       elf.Prog64
	ShShstrtab elf.Section64
	ShSymtab   *elf.Section64
	ShStrtab   *elf.Section64
	ShDynstr   *elf.Section64
}

// Bounds checks
const (
	// 64 MiB max program size.
	// Allows loader to use unchecked math when adding 32-bit offsets.
	maxFileLen = 1 << 26

	maxSectionNameLen = 16
	maxSymbolNameLen  = 1024
)

func LoadProgram(buf []byte) (*Executable, error) {
	if len(buf) > maxFileLen {
		return nil, fmt.Errorf("ELF file too large")
	}
	l := loader{
		rd:       bytes.NewReader(buf),
		fileSize: uint64(len(buf)),
	}
	return l.load()
}

const EF_SBF_V2 = 0x20

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
	if err := l.readSectionHeaderTable(); err != nil {
		return nil, err
	}
	if err := l.parseSections(); err != nil {
		return nil, err
	}
	// TODO parse dynamic segment
	return l.elf, nil
}

const (
	ehsize    = 0x40
	phentsize = 0x38
	shentsize = 0x40
)

func (l *loader) newPhTableIter() *tableIter[elf.Prog64] {
	eh := &l.elf.Header
	return newTableIterator[elf.Prog64](l, eh.Phoff, eh.Phnum, phentsize)
}

func (l *loader) newShTableIter() *tableIter[elf.Section64] {
	eh := &l.elf.Header
	return newTableIterator[elf.Section64](l, eh.Shoff, eh.Shnum, shentsize)
}

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
		elf.OSABI(ident[elf.EI_OSABI]) != elf.ELFOSABI_NONE ||
		elf.Machine(eh.Machine) != elf.EM_BPF ||
		elf.Type(eh.Type) != elf.ET_DYN {
		return fmt.Errorf("incompatible binary")
	}
	// note: EI_PAD and EI_ABIVERSION are ignored

	if eh.Version != uint32(elf.EV_CURRENT) ||
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
	iter := l.newPhTableIter()
	for iter.Next() && iter.Err() == nil {
		ph := iter.Item()

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
	return iter.Err()
}

// reads and validates the section header table.
// remembers the section header table.
func (l *loader) readSectionHeaderTable() error {
	eh := &l.elf.Header
	iter := l.newShTableIter()
	sectionDataOff := uint64(0)

	if !iter.Next() {
		return fmt.Errorf("missing section 0")
	}
	if elf.SectionType(iter.Item().Type) != elf.SHT_NULL {
		return fmt.Errorf("section 0 is not SHT_NULL")
	}

	for iter.Next() && iter.Err() == nil {
		i, sh := iter.Index(), iter.Item()
		if elf.SectionType(sh.Type) == elf.SHT_NOBITS {
			continue
		}

		// Ensure section data is not overlapping with ELF headers
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

		// More checks
		if eh.Shoff < sectionDataOff {
			return fmt.Errorf("sections not in order")
		}
		if shend > l.fileSize {
			return fmt.Errorf("section %d out of bounds", i)
		}

		// Remember section header string table.
		if eh.Shstrndx != uint16(elf.SHN_UNDEF) && eh.Shstrndx == i {
			l.elf.ShShstrtab = sh
		}

		sectionDataOff = shend
	}
	// TODO validate offset and size (?)
	if elf.SectionType(l.elf.ShShstrtab.Type) != elf.SHT_STRTAB {
		return fmt.Errorf("invalid .shstrtab")
	}
	return iter.Err()
}

func (l *loader) getString(strtab *elf.Section64, stroff uint32, maxLen uint16) (string, error) {
	if elf.SectionType(strtab.Type) != elf.SHT_STRTAB {
		return "", fmt.Errorf("invalid strtab")
	}
	offset := strtab.Off + uint64(stroff)
	if offset > l.fileSize || offset+uint64(maxLen) > l.fileSize {
		return "", io.ErrUnexpectedEOF
	}
	rd := bufio.NewReader(io.NewSectionReader(l.rd, int64(offset), int64(maxLen)))
	var builder strings.Builder
	for {
		b, err := rd.ReadByte()
		if err != nil {
			return "", err
		}
		if b == 0 {
			break
		}
		builder.WriteByte(b)
	}
	return builder.String(), nil
}

// Iterate sections and remember special sections by name.
func (l *loader) parseSections() error {
	shShstrtab := &l.elf.ShShstrtab
	iter := l.newShTableIter()
	for iter.Next() && iter.Err() == nil {
		sh := iter.Item()
		sectionName, err := l.getString(shShstrtab, sh.Name, maxSectionNameLen)
		if err != nil {
			return fmt.Errorf("getString: %w", err)
		}

		// Remember special section or error if it already exists.
		setSection := func(shPtr **elf.Section64) error {
			if *shPtr != nil {
				return fmt.Errorf("duplicate section: %s", sectionName)
			}
			*shPtr = new(elf.Section64)
			**shPtr = sh
			return nil
		}
		switch sectionName {
		case ".symtab":
			err = setSection(&l.elf.ShSymtab)
		case ".strtab":
			err = setSection(&l.elf.ShStrtab)
		case ".dynstr":
			err = setSection(&l.elf.ShDynstr)
		}
		if err != nil {
			return err
		}
	}
	return iter.Err()
}

// tableIter is a memory-efficient iterator over densely packed tables of statically sized items.
// Such as the ELF program header and section header tables.
type tableIter[T any] struct {
	l        *loader
	off      uint64
	i        uint16 // one ahead
	count    uint16
	elemSize uint16
	elem     T
	err      error
}

// newTableIterator creates a new tableIter at `off` for `count` elements of `elemSize` len.
func newTableIterator[T any](l *loader, off uint64, count uint16, elemSize uint16) *tableIter[T] {
	return &tableIter[T]{
		l:        l,
		off:      off,
		count:    count,
		elemSize: elemSize,
	}
}

// Next reads one element.
//
// Returns true on success, false if table end has been reached or error occurred.
// The caller should abort iteration on error.
func (it *tableIter[T]) Next() (ok bool) {
	ok, it.err = it.getNext()
	if ok && it.err != nil {
		panic("unreachable")
	}
	return
}

// Index returns the current table index.
func (it *tableIter[T]) Index() uint16 {
	return it.i - 1
}

// Err returns the current error.
func (it *tableIter[T]) Err() error {
	return it.err
}

// Item returns the current element read.
//
// Next must be called before.
func (it *tableIter[T]) Item() T {
	return it.elem
}

func (it *tableIter[T]) getNext() (bool, error) {
	if it.i >= it.count {
		return false, nil
	}
	if it.off >= math.MaxInt64 || it.off+uint64(it.elemSize) > math.MaxInt64 {
		return false, io.ErrUnexpectedEOF
	}

	rd := io.NewSectionReader(it.l.rd, int64(it.off), int64(it.elemSize))
	if err := binary.Read(rd, binary.LittleEndian, &it.elem); err != nil {
		return false, err
	}

	it.off += uint64(it.elemSize)
	it.i++
	return true, nil
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
