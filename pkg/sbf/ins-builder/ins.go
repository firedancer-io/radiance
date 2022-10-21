package insbuilder

import (
	"encoding/binary"
	"fmt"

	"github.com/spaolacci/murmur3"
)

type Insn struct {
	// Instruction pointer.
	Ptr uint64

	// Operation code.
	Opc uint8

	// Destination register operand.
	Dst uint8

	// Source register operand.
	Src uint8

	// Offset operand.
	Off int16

	// Immediate value operand.
	Imm int64
}

func (i Insn) String() string {
	return fmt.Sprintf(
		"Insn { ptr: 0x%x, opc: 0x%x, dst: %v, src: %v, off: 0x%x, imm: 0x%x }",
		i.Ptr, i.Opc, i.Dst, i.Src, i.Off, i.Imm,
	)
}

func (insn *Insn) ToArray() [INSN_SIZE]uint8 {
	var bytes [8]uint8
	bytes[0] = insn.Opc
	bytes[1] = insn.Src<<4 | insn.Dst
	bytes[2] = uint8(insn.Off & 0xff)
	bytes[3] = uint8(insn.Off >> 8)
	bytes[4] = uint8(insn.Imm & 0xff)
	bytes[5] = uint8((insn.Imm & 0xff00) >> 8)
	bytes[6] = uint8((insn.Imm & 0xff0000) >> 16)
	bytes[7] = uint8((insn.Imm & 0xff000000) >> 24)
	return bytes
}

// ToSlice converts an array to a slice.
// ToSlice is the equivalent of rusts's to_vec() function.
func (insn *Insn) ToSlice() []uint8 {
	arr := insn.ToArray()
	return arr[:]
}

const (
	// SBF version flag
	EF_SBF_V2 = 0x20
	// Maximum number of instructions in an eBPF program.
	PROG_MAX_INSNS = 65_536
	// Size of an eBPF instructions, in bytes.
	INSN_SIZE = 8
	// Frame pointer register
	FRAME_PTR_REG = 10
	// Stack pointer register
	STACK_PTR_REG = 11
	// First scratch register
	FIRST_SCRATCH_REG = 6
	// Number of scratch registers
	SCRATCH_REGS = 4
	// ELF dump instruction offset
	// Instruction numbers typically start at 29 in the ELF dump, use this offset
	// when reporting so that trace aligns with the dump.
	ELF_INSN_DUMP_OFFSET = 29
	// Alignment of the memory regions in host address space in bytes
	HOST_ALIGN = 16
	// Upper half of a pointer is the region index, lower half the virtual address inside that region.
	VIRTUAL_ADDRESS_BITS = 32
)

const (
	// Memory map regions virtual addresses need to be (1 << VIRTUAL_ADDRESS_BITS) bytes apart.
	// Also the region at index 0 should be skipped to catch NULL ptr accesses.

	// Start of the program bits (text and ro segments) in the memory map
	MM_PROGRAM_START = 0x100000000
	// Start of the stack in the memory map
	MM_STACK_START = 0x200000000
	// Start of the heap in the memory map
	MM_HEAP_START = 0x300000000
	// Start of the input buffers in the memory map
	MM_INPUT_START = 0x400000000
)

const (
	// eBPF op codes.
	// See also https://www.kernel.org/doc/Documentation/networking/filter.txt

	// Three least significant bits are operation class:
	// BPF operation class: load from immediate.
	BPF_LD uint8 = 0x00
	// BPF operation class: load from register.
	BPF_LDX uint8 = 0x01
	// BPF operation class: store immediate.
	BPF_ST uint8 = 0x02
	// BPF operation class: store value from register.
	BPF_STX uint8 = 0x03
	// BPF operation class: 32 bits arithmetic operation.
	BPF_ALU uint8 = 0x04
	// BPF operation class: jump.
	BPF_JMP uint8 = 0x05
	// [ class 6 unused, reserved for future use ]
	// BPF operation class: 64 bits arithmetic operation.
	BPF_ALU64 uint8 = 0x07
)

const (
	// For load and store instructions:
	// +------------+--------+------------+
	// |   3 bits   | 2 bits |   3 bits   |
	// |    mode    |  size  | insn class |
	// +------------+--------+------------+
	// (MSB)                          (LSB)

	// Size modifiers:
	// BPF size modifier: word (4 bytes).
	BPF_W uint8 = 0x00
	// BPF size modifier: half-word (2 bytes).
	BPF_H uint8 = 0x08
	// BPF size modifier: byte (1 byte).
	BPF_B uint8 = 0x10
	// BPF size modifier: double word (8 bytes).
	BPF_DW uint8 = 0x18

	// Mode modifiers:
	// BPF mode modifier: immediate value.
	BPF_IMM uint8 = 0x00
	// BPF mode modifier: absolute load.
	BPF_ABS uint8 = 0x20
	// BPF mode modifier: indirect load.
	BPF_IND uint8 = 0x40
	// BPF mode modifier: load from / store to memory.
	BPF_MEM uint8 = 0x60
	// [ 0x80 reserved ]
	// [ 0xa0 reserved ]
	// BPF mode modifier: exclusive add.
	BPF_XADD uint8 = 0xc0
)

const (
	// For arithmetic (BPF_ALU/BPF_ALU64) and jump (BPF_JMP) instructions:
	// +----------------+--------+--------+
	// |     4 bits     |1 b.|   3 bits   |
	// | operation code | src| insn class |
	// +----------------+----+------------+
	// (MSB)                          (LSB)

	// Source modifiers:
	// BPF source operand modifier: 32-bit immediate value.
	BPF_K uint8 = 0x00
	// BPF source operand modifier: `src` register.
	BPF_X uint8 = 0x08

	// Operation codes -- BPF_ALU or BPF_ALU64 classes:
	// BPF ALU/ALU64 operation code: addition.
	BPF_ADD uint8 = 0x00
	// BPF ALU/ALU64 operation code: subtraction.
	BPF_SUB uint8 = 0x10
	// BPF ALU/ALU64 operation code: multiplication.
	BPF_MUL uint8 = 0x20
	// BPF ALU/ALU64 operation code: division.
	BPF_DIV uint8 = 0x30
	// BPF ALU/ALU64 operation code: or.
	BPF_OR uint8 = 0x40
	// BPF ALU/ALU64 operation code: and.
	BPF_AND uint8 = 0x50
	// BPF ALU/ALU64 operation code: left shift.
	BPF_LSH uint8 = 0x60
	// BPF ALU/ALU64 operation code: right shift.
	BPF_RSH uint8 = 0x70
	// BPF ALU/ALU64 operation code: negation.
	BPF_NEG uint8 = 0x80
	// BPF ALU/ALU64 operation code: modulus.
	BPF_MOD uint8 = 0x90
	// BPF ALU/ALU64 operation code: exclusive or.
	BPF_XOR uint8 = 0xa0
	// BPF ALU/ALU64 operation code: move.
	BPF_MOV uint8 = 0xb0
	// BPF ALU/ALU64 operation code: sign extending right shift.
	BPF_ARSH uint8 = 0xc0
	// BPF ALU/ALU64 operation code: endianness conversion.
	BPF_END uint8 = 0xd0
	// BPF ALU/ALU64 operation code: signed division.
	BPF_SDIV uint8 = 0xe0
)

const (
	// Operation codes -- BPF_JMP class:
	// BPF JMP operation code: jump.
	BPF_JA uint8 = 0x00
	// BPF JMP operation code: jump if equal.
	BPF_JEQ uint8 = 0x10
	// BPF JMP operation code: jump if greater than.
	BPF_JGT uint8 = 0x20
	// BPF JMP operation code: jump if greater or equal.
	BPF_JGE uint8 = 0x30
	// BPF JMP operation code: jump if `src` & `reg`.
	BPF_JSET uint8 = 0x40
	// BPF JMP operation code: jump if not equal.
	BPF_JNE uint8 = 0x50
	// BPF JMP operation code: jump if greater than (signed).
	BPF_JSGT uint8 = 0x60
	// BPF JMP operation code: jump if greater or equal (signed).
	BPF_JSGE uint8 = 0x70
	// BPF JMP operation code: syscall function call.
	BPF_CALL uint8 = 0x80
	// BPF JMP operation code: return from program.
	BPF_EXIT uint8 = 0x90
	// BPF JMP operation code: jump if lower than.
	BPF_JLT uint8 = 0xa0
	// BPF JMP operation code: jump if lower or equal.
	BPF_JLE uint8 = 0xb0
	// BPF JMP operation code: jump if lower than (signed).
	BPF_JSLT uint8 = 0xc0
	// BPF JMP operation code: jump if lower or equal (signed).
	BPF_JSLE uint8 = 0xd0
)

const (
	// Op codes
	// (Following operation names are not “official”, but may be proper to rbpf; Linux kernel only
	// combines above flags and does not attribute a name per operation.)

	// BPF opcode: `ldabsb src, dst, imm`.
	LD_ABS_B uint8 = BPF_LD | BPF_ABS | BPF_B
	// BPF opcode: `ldabsh src, dst, imm`.
	LD_ABS_H uint8 = BPF_LD | BPF_ABS | BPF_H
	// BPF opcode: `ldabsw src, dst, imm`.
	LD_ABS_W uint8 = BPF_LD | BPF_ABS | BPF_W
	// BPF opcode: `ldabsdw src, dst, imm`.
	LD_ABS_DW uint8 = BPF_LD | BPF_ABS | BPF_DW
	// BPF opcode: `ldindb src, dst, imm`.
	LD_IND_B uint8 = BPF_LD | BPF_IND | BPF_B
	// BPF opcode: `ldindh src, dst, imm`.
	LD_IND_H uint8 = BPF_LD | BPF_IND | BPF_H
	// BPF opcode: `ldindw src, dst, imm`.
	LD_IND_W uint8 = BPF_LD | BPF_IND | BPF_W
	// BPF opcode: `ldinddw src, dst, imm`.
	LD_IND_DW uint8 = BPF_LD | BPF_IND | BPF_DW

	// BPF opcode: `lddw dst, imm` // `dst = imm`.
	LD_DW_IMM uint8 = BPF_LD | BPF_IMM | BPF_DW
	// BPF opcode: `ldxb dst, [src + off]` // `dst = (src + off) as u8`.
	LD_B_REG uint8 = BPF_LDX | BPF_MEM | BPF_B
	// BPF opcode: `ldxh dst, [src + off]` // `dst = (src + off) as u16`.
	LD_H_REG uint8 = BPF_LDX | BPF_MEM | BPF_H
	// BPF opcode: `ldxw dst, [src + off]` // `dst = (src + off) as u32`.
	LD_W_REG uint8 = BPF_LDX | BPF_MEM | BPF_W
	// BPF opcode: `ldxdw dst, [src + off]` // `dst = (src + off) as u64`.
	LD_DW_REG uint8 = BPF_LDX | BPF_MEM | BPF_DW
	// BPF opcode: `stb [dst + off], imm` // `(dst + offset) as u8 = imm`.
	ST_B_IMM uint8 = BPF_ST | BPF_MEM | BPF_B
	// BPF opcode: `sth [dst + off], imm` // `(dst + offset) as u16 = imm`.
	ST_H_IMM uint8 = BPF_ST | BPF_MEM | BPF_H
	// BPF opcode: `stw [dst + off], imm` // `(dst + offset) as u32 = imm`.
	ST_W_IMM uint8 = BPF_ST | BPF_MEM | BPF_W
	// BPF opcode: `stdw [dst + off], imm` // `(dst + offset) as u64 = imm`.
	ST_DW_IMM uint8 = BPF_ST | BPF_MEM | BPF_DW
	// BPF opcode: `stxb [dst + off], src` // `(dst + offset) as u8 = src`.
	ST_B_REG uint8 = BPF_STX | BPF_MEM | BPF_B
	// BPF opcode: `stxh [dst + off], src` // `(dst + offset) as u16 = src`.
	ST_H_REG uint8 = BPF_STX | BPF_MEM | BPF_H
	// BPF opcode: `stxw [dst + off], src` // `(dst + offset) as u32 = src`.
	ST_W_REG uint8 = BPF_STX | BPF_MEM | BPF_W
	// BPF opcode: `stxdw [dst + off], src` // `(dst + offset) as u64 = src`.
	ST_DW_REG uint8 = BPF_STX | BPF_MEM | BPF_DW

	// BPF opcode: `stxxaddw [dst + off], src`.
	ST_W_XADD uint8 = BPF_STX | BPF_XADD | BPF_W
	// BPF opcode: `stxxadddw [dst + off], src`.
	ST_DW_XADD uint8 = BPF_STX | BPF_XADD | BPF_DW

	// BPF opcode: `add32 dst, imm` // `dst += imm`.
	ADD32_IMM uint8 = BPF_ALU | BPF_K | BPF_ADD
	// BPF opcode: `add32 dst, src` // `dst += src`.
	ADD32_REG uint8 = BPF_ALU | BPF_X | BPF_ADD
	// BPF opcode: `sub32 dst, imm` // `dst -= imm`.
	SUB32_IMM uint8 = BPF_ALU | BPF_K | BPF_SUB
	// BPF opcode: `sub32 dst, src` // `dst -= src`.
	SUB32_REG uint8 = BPF_ALU | BPF_X | BPF_SUB
	// BPF opcode: `mul32 dst, imm` // `dst *= imm`.
	MUL32_IMM uint8 = BPF_ALU | BPF_K | BPF_MUL
	// BPF opcode: `mul32 dst, src` // `dst *= src`.
	MUL32_REG uint8 = BPF_ALU | BPF_X | BPF_MUL
	// BPF opcode: `div32 dst, imm` // `dst /= imm`.
	DIV32_IMM uint8 = BPF_ALU | BPF_K | BPF_DIV
	// BPF opcode: `div32 dst, src` // `dst /= src`.
	DIV32_REG uint8 = BPF_ALU | BPF_X | BPF_DIV
	// BPF opcode: `or32 dst, imm` // `dst |= imm`.
	OR32_IMM uint8 = BPF_ALU | BPF_K | BPF_OR
	// BPF opcode: `or32 dst, src` // `dst |= src`.
	OR32_REG uint8 = BPF_ALU | BPF_X | BPF_OR
	// BPF opcode: `and32 dst, imm` // `dst &= imm`.
	AND32_IMM uint8 = BPF_ALU | BPF_K | BPF_AND
	// BPF opcode: `and32 dst, src` // `dst &= src`.
	AND32_REG uint8 = BPF_ALU | BPF_X | BPF_AND
	// BPF opcode: `lsh32 dst, imm` // `dst <<= imm`.
	LSH32_IMM uint8 = BPF_ALU | BPF_K | BPF_LSH
	// BPF opcode: `lsh32 dst, src` // `dst <<= src`.
	LSH32_REG uint8 = BPF_ALU | BPF_X | BPF_LSH
	// BPF opcode: `rsh32 dst, imm` // `dst >>= imm`.
	RSH32_IMM uint8 = BPF_ALU | BPF_K | BPF_RSH
	// BPF opcode: `rsh32 dst, src` // `dst >>= src`.
	RSH32_REG uint8 = BPF_ALU | BPF_X | BPF_RSH
	// BPF opcode: `neg32 dst` // `dst = -dst`.
	NEG32 uint8 = BPF_ALU | BPF_NEG
	// BPF opcode: `mod32 dst, imm` // `dst %= imm`.
	MOD32_IMM uint8 = BPF_ALU | BPF_K | BPF_MOD
	// BPF opcode: `mod32 dst, src` // `dst %= src`.
	MOD32_REG uint8 = BPF_ALU | BPF_X | BPF_MOD
	// BPF opcode: `xor32 dst, imm` // `dst ^= imm`.
	XOR32_IMM uint8 = BPF_ALU | BPF_K | BPF_XOR
	// BPF opcode: `xor32 dst, src` // `dst ^= src`.
	XOR32_REG uint8 = BPF_ALU | BPF_X | BPF_XOR
	// BPF opcode: `mov32 dst, imm` // `dst = imm`.
	MOV32_IMM uint8 = BPF_ALU | BPF_K | BPF_MOV
	// BPF opcode: `mov32 dst, src` // `dst = src`.
	MOV32_REG uint8 = BPF_ALU | BPF_X | BPF_MOV
	// BPF opcode: `arsh32 dst, imm` // `dst >>= imm (arithmetic)`.
	///
	// <https://en.wikipedia.org/wiki/Arithmetic_shift>
	ARSH32_IMM uint8 = BPF_ALU | BPF_K | BPF_ARSH
	// BPF opcode: `arsh32 dst, src` // `dst >>= src (arithmetic)`.
	///
	// <https://en.wikipedia.org/wiki/Arithmetic_shift>
	ARSH32_REG uint8 = BPF_ALU | BPF_X | BPF_ARSH
	// BPF opcode: `sdiv32 dst, imm` // `dst s/= imm`.
	SDIV32_IMM uint8 = BPF_ALU | BPF_K | BPF_SDIV
	// BPF opcode: `sdiv32 dst, src` // `dst s/= src`.
	SDIV32_REG uint8 = BPF_ALU | BPF_X | BPF_SDIV

	// BPF opcode: `le dst` // `dst = htole<imm>(dst), with imm in {16, 32, 64}`.
	LE uint8 = BPF_ALU | BPF_K | BPF_END
	// BPF opcode: `be dst` // `dst = htobe<imm>(dst), with imm in {16, 32, 64}`.
	BE uint8 = BPF_ALU | BPF_X | BPF_END

	// BPF opcode: `add64 dst, imm` // `dst += imm`.
	ADD64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_ADD
	// BPF opcode: `add64 dst, src` // `dst += src`.
	ADD64_REG uint8 = BPF_ALU64 | BPF_X | BPF_ADD
	// BPF opcode: `sub64 dst, imm` // `dst -= imm`.
	SUB64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_SUB
	// BPF opcode: `sub64 dst, src` // `dst -= src`.
	SUB64_REG uint8 = BPF_ALU64 | BPF_X | BPF_SUB
	// BPF opcode: `div64 dst, imm` // `dst /= imm`.
	MUL64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_MUL
	// BPF opcode: `div64 dst, src` // `dst /= src`.
	MUL64_REG uint8 = BPF_ALU64 | BPF_X | BPF_MUL
	// BPF opcode: `div64 dst, imm` // `dst /= imm`.
	DIV64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_DIV
	// BPF opcode: `div64 dst, src` // `dst /= src`.
	DIV64_REG uint8 = BPF_ALU64 | BPF_X | BPF_DIV
	// BPF opcode: `or64 dst, imm` // `dst |= imm`.
	OR64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_OR
	// BPF opcode: `or64 dst, src` // `dst |= src`.
	OR64_REG uint8 = BPF_ALU64 | BPF_X | BPF_OR
	// BPF opcode: `and64 dst, imm` // `dst &= imm`.
	AND64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_AND
	// BPF opcode: `and64 dst, src` // `dst &= src`.
	AND64_REG uint8 = BPF_ALU64 | BPF_X | BPF_AND
	// BPF opcode: `lsh64 dst, imm` // `dst <<= imm`.
	LSH64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_LSH
	// BPF opcode: `lsh64 dst, src` // `dst <<= src`.
	LSH64_REG uint8 = BPF_ALU64 | BPF_X | BPF_LSH
	// BPF opcode: `rsh64 dst, imm` // `dst >>= imm`.
	RSH64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_RSH
	// BPF opcode: `rsh64 dst, src` // `dst >>= src`.
	RSH64_REG uint8 = BPF_ALU64 | BPF_X | BPF_RSH
	// BPF opcode: `neg64 dst, imm` // `dst = -dst`.
	NEG64 uint8 = BPF_ALU64 | BPF_NEG
	// BPF opcode: `mod64 dst, imm` // `dst %= imm`.
	MOD64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_MOD
	// BPF opcode: `mod64 dst, src` // `dst %= src`.
	MOD64_REG uint8 = BPF_ALU64 | BPF_X | BPF_MOD
	// BPF opcode: `xor64 dst, imm` // `dst ^= imm`.
	XOR64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_XOR
	// BPF opcode: `xor64 dst, src` // `dst ^= src`.
	XOR64_REG uint8 = BPF_ALU64 | BPF_X | BPF_XOR
	// BPF opcode: `mov64 dst, imm` // `dst = imm`.
	MOV64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_MOV
	// BPF opcode: `mov64 dst, src` // `dst = src`.
	MOV64_REG uint8 = BPF_ALU64 | BPF_X | BPF_MOV
	// BPF opcode: `arsh64 dst, imm` // `dst >>= imm (arithmetic)`.
	///
	// <https://en.wikipedia.org/wiki/Arithmetic_shift>
	ARSH64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_ARSH
	// BPF opcode: `arsh64 dst, src` // `dst >>= src (arithmetic)`.
	///
	// <https://en.wikipedia.org/wiki/Arithmetic_shift>
	ARSH64_REG uint8 = BPF_ALU64 | BPF_X | BPF_ARSH
	// BPF opcode: `sdiv64 dst, imm` // `dst s/= imm`.
	SDIV64_IMM uint8 = BPF_ALU64 | BPF_K | BPF_SDIV
	// BPF opcode: `sdiv64 dst, src` // `dst s/= src`.
	SDIV64_REG uint8 = BPF_ALU64 | BPF_X | BPF_SDIV

	// BPF opcode: `ja +off` // `PC += off`.
	JA uint8 = BPF_JMP | BPF_JA
	// BPF opcode: `jeq dst, imm, +off` // `PC += off if dst == imm`.
	JEQ_IMM uint8 = BPF_JMP | BPF_K | BPF_JEQ
	// BPF opcode: `jeq dst, src, +off` // `PC += off if dst == src`.
	JEQ_REG uint8 = BPF_JMP | BPF_X | BPF_JEQ
	// BPF opcode: `jgt dst, imm, +off` // `PC += off if dst > imm`.
	JGT_IMM uint8 = BPF_JMP | BPF_K | BPF_JGT
	// BPF opcode: `jgt dst, src, +off` // `PC += off if dst > src`.
	JGT_REG uint8 = BPF_JMP | BPF_X | BPF_JGT
	// BPF opcode: `jge dst, imm, +off` // `PC += off if dst >= imm`.
	JGE_IMM uint8 = BPF_JMP | BPF_K | BPF_JGE
	// BPF opcode: `jge dst, src, +off` // `PC += off if dst >= src`.
	JGE_REG uint8 = BPF_JMP | BPF_X | BPF_JGE
	// BPF opcode: `jlt dst, imm, +off` // `PC += off if dst < imm`.
	JLT_IMM uint8 = BPF_JMP | BPF_K | BPF_JLT
	// BPF opcode: `jlt dst, src, +off` // `PC += off if dst < src`.
	JLT_REG uint8 = BPF_JMP | BPF_X | BPF_JLT
	// BPF opcode: `jle dst, imm, +off` // `PC += off if dst <= imm`.
	JLE_IMM uint8 = BPF_JMP | BPF_K | BPF_JLE
	// BPF opcode: `jle dst, src, +off` // `PC += off if dst <= src`.
	JLE_REG uint8 = BPF_JMP | BPF_X | BPF_JLE
	// BPF opcode: `jset dst, imm, +off` // `PC += off if dst & imm`.
	JSET_IMM uint8 = BPF_JMP | BPF_K | BPF_JSET
	// BPF opcode: `jset dst, src, +off` // `PC += off if dst & src`.
	JSET_REG uint8 = BPF_JMP | BPF_X | BPF_JSET
	// BPF opcode: `jne dst, imm, +off` // `PC += off if dst != imm`.
	JNE_IMM uint8 = BPF_JMP | BPF_K | BPF_JNE
	// BPF opcode: `jne dst, src, +off` // `PC += off if dst != src`.
	JNE_REG uint8 = BPF_JMP | BPF_X | BPF_JNE
	// BPF opcode: `jsgt dst, imm, +off` // `PC += off if dst > imm (signed)`.
	JSGT_IMM uint8 = BPF_JMP | BPF_K | BPF_JSGT
	// BPF opcode: `jsgt dst, src, +off` // `PC += off if dst > src (signed)`.
	JSGT_REG uint8 = BPF_JMP | BPF_X | BPF_JSGT
	// BPF opcode: `jsge dst, imm, +off` // `PC += off if dst >= imm (signed)`.
	JSGE_IMM uint8 = BPF_JMP | BPF_K | BPF_JSGE
	// BPF opcode: `jsge dst, src, +off` // `PC += off if dst >= src (signed)`.
	JSGE_REG uint8 = BPF_JMP | BPF_X | BPF_JSGE
	// BPF opcode: `jslt dst, imm, +off` // `PC += off if dst < imm (signed)`.
	JSLT_IMM uint8 = BPF_JMP | BPF_K | BPF_JSLT
	// BPF opcode: `jslt dst, src, +off` // `PC += off if dst < src (signed)`.
	JSLT_REG uint8 = BPF_JMP | BPF_X | BPF_JSLT
	// BPF opcode: `jsle dst, imm, +off` // `PC += off if dst <= imm (signed)`.
	JSLE_IMM uint8 = BPF_JMP | BPF_K | BPF_JSLE
	// BPF opcode: `jsle dst, src, +off` // `PC += off if dst <= src (signed)`.
	JSLE_REG uint8 = BPF_JMP | BPF_X | BPF_JSLE

	// BPF opcode: `call imm` // syscall function call to syscall with key `imm`.
	CALL_IMM uint8 = BPF_JMP | BPF_CALL
	// BPF opcode: tail call.
	CALL_REG uint8 = BPF_JMP | BPF_X | BPF_CALL
	// BPF opcode: `exit` // `return r0`.
	EXIT uint8 = BPF_JMP | BPF_EXIT

	// Used in JIT
	// Mask to extract the operation class from an operation code.
	BPF_CLS_MASK = 0x07
	// Mask to extract the arithmetic operation code from an instruction operation code.
	BPF_ALU_OP_MASK = 0xf0
)

// /// Get the instruction at `idx` of an eBPF program. `idx` is the index (number) of the
// /// instruction (not a byte offset). The first instruction has index 0.
// ///
// /// # Panics
// ///
// /// Panics if it is not possible to get the instruction (if idx is too high, or last instruction is
// /// incomplete).
// ///
// /// # Examples
// ///
// /// ```
// /// use solana_rbpf::ebpf;
// ///
// /// let prog = &[
// ///     0xb7, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
// ///     0x95, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00
// ///     ];
// /// let insn = ebpf::get_insn(prog, 1);
// /// assert_eq!(insn.opc, 0x95);
// /// ```
// ///
// /// The example below will panic, since the last instruction is not complete and cannot be loaded.
// ///
// /// ```rust,should_panic
// /// use solana_rbpf::ebpf;
// ///
// /// let prog = &[
// ///     0xb7, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
// ///     0x95, 0x00, 0x00, 0x00, 0x00, 0x00              // two bytes missing
// ///     ];
// /// let insn = ebpf::get_insn(prog, 1);
// /// ```
//
//	pub fn get_insn(prog: &[u8], pc: usize) -> Insn {
//	    // This guard should not be needed in most cases, since the verifier already checks the program
//	    // size, and indexes should be fine in the interpreter/JIT. But this function is publicly
//	    // available and user can call it with any `pc`, so we have to check anyway.
//	    debug_assert!(
//	        (pc + 1) * INSN_SIZE <= prog.len(),
//	        "cannot reach instruction at index {:?} in program containing {:?} bytes",
//	        pc,
//	        prog.len()
//	    );
//	    get_insn_unchecked(prog, pc)
//	}
func GetInsn(prog []uint8, pc uint64) Insn {
	// This guard should not be needed in most cases, since the verifier already checks the program
	// size, and indexes should be fine in the interpreter/JIT. But this function is publicly
	// available and user can call it with any `pc`, so we have to check anyway.
	if (pc+1)*INSN_SIZE > uint64(len(prog)) {
		panic(fmt.Sprintf(
			"cannot reach instruction at index %d in program containing %d bytes",
			pc, len(prog),
		))
	}
	return GetInsnUnchecked(prog, pc)
}

// /// Same as `get_insn` except not checked
//
//	pub fn get_insn_unchecked(prog: &[u8], pc: usize) -> Insn {
//	    Insn {
//	        ptr: pc,
//	        opc: prog[INSN_SIZE * pc],
//	        dst: prog[INSN_SIZE * pc + 1] & 0x0f,
//	        src: (prog[INSN_SIZE * pc + 1] & 0xf0) >> 4,
//	        off: LittleEndian::read_i16(&prog[(INSN_SIZE * pc + 2)..]),
//	        imm: LittleEndian::read_i32(&prog[(INSN_SIZE * pc + 4)..]) as i64,
//	    }
//	}
func GetInsnUnchecked(prog []uint8, pc uint64) Insn {
	return Insn{
		Ptr: pc,
		Opc: prog[INSN_SIZE*pc],
		Dst: prog[INSN_SIZE*pc+1] & 0x0f,
		Src: (prog[INSN_SIZE*pc+1] & 0xf0) >> 4,
		Off: int16(binary.LittleEndian.Uint16(prog[INSN_SIZE*pc+2:])),
		Imm: int64(binary.LittleEndian.Uint32(prog[INSN_SIZE*pc+4:])),
	}
}

// /// Merge the two halves of a LD_DW_IMM instruction
//
//	pub fn augment_lddw_unchecked(prog: &[u8], insn: &mut Insn) {
//	    let more_significant_half = LittleEndian::read_i32(&prog[((insn.ptr + 1) * INSN_SIZE + 4)..]);
//	    insn.imm = ((insn.imm as u64 & 0xffffffff) | ((more_significant_half as u64) << 32)) as i64;
//	}
func AugmentLddwUnchecked(prog []uint8, insn *Insn) {
	more_significant_half := binary.LittleEndian.Uint32(prog[(insn.Ptr+1)*INSN_SIZE+4:])
	insn.Imm = int64((uint64(insn.Imm&0xffffffff) | (uint64(more_significant_half) << 32)))
}

// /// Hash a symbol name
// ///
// /// This function is used by both the relocator and the VM to translate symbol names
// /// into a 32 bit id used to identify a syscall function.  The 32 bit id is used in the
// /// eBPF `call` instruction's imm field.
//
//	pub fn hash_symbol_name(name: &[u8]) -> u32 {
//	    let mut hasher = Murmur3Hasher::default();
//	    Hash::hash_slice(name, &mut hasher);
//	    hasher.finish()
//	}
func HashSymbolName(name []uint8) uint32 {
	hasher := murmur3.New32()
	hasher.Write(name)
	return hasher.Sum32()
}
