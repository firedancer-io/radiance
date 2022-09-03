// Package sbf implements the Solana Bytecode Format.
package sbf

import "encoding/binary"

// Hardcoded addresses.
const (
	VaddrProgram = uint64(0x1_0000_0000)
	VaddrStack   = uint64(0x2_0000_0000)
	VaddrHeap    = uint64(0x3_0000_0000)
	VaddrInput   = uint64(0x4_0000_0000)
)

const (
	// SlotSize is the size of one instruction slot.
	SlotSize = 8
	// MinInsSize is the size of the shortest possible instruction
	MinInsSize = SlotSize
	// MaxInsSize is the size of the longest possible instruction (lddw)
	MaxInsSize = 2 * SlotSize
)

func IsLongIns(op uint8) bool {
	return op == OpLddw
}

// Slot holds the content of one instruction slot.
type Slot uint64

// GetSlot reads an instruction slot from memory.
func GetSlot(buf []byte) Slot {
	return Slot(binary.LittleEndian.Uint64(buf))
}

// Op returns the opcode field.
func (s Slot) Op() uint8 {
	return uint8(s >> 56)
}

// Dst returns the destination register field.
func (s Slot) Dst() uint8 {
	return uint8(s>>52) & 0xF
}

// Src returns the source register field.
func (s Slot) Src() uint8 {
	return uint8(s>>48) & 0xF
}

// Off returns the offset field.
func (s Slot) Off() uint16 {
	return uint16(s >> 32)
}

// Imm returns the immediate field.
func (s Slot) Imm() int32 {
	return int32(uint32(s))
}
