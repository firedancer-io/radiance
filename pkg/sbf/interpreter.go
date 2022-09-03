package sbf

import (
	"fmt"
	"math"
	"math/bits"
)

// Interpreter implements the SBF core in pure Go.
type Interpreter struct {
	text  []byte
	ro    []byte
	stack []byte
	heap  []byte
	input []byte

	entry uint64

	cuMax  uint64
	cuLeft uint64

	syscalls  map[uint32]Syscall
	vmContext any
}

// NewInterpreter creates a new interpreter instance for a program execution.
//
// The caller must create a new interpreter object for every new execution.
// In other words, Run may only be called once per interpreter.
func NewInterpreter(p *Program, opts *VMOpts) *Interpreter {
	return &Interpreter{
		text:      p.Text,
		ro:        p.RO,
		stack:     make([]byte, opts.StackSize),
		heap:      make([]byte, opts.HeapSize),
		input:     opts.Input,
		entry:     p.Entrypoint,
		cuMax:     opts.MaxCU,
		cuLeft:    opts.MaxCU,
		syscalls:  opts.Syscalls,
		vmContext: opts.Context,
	}
}

// Run executes the program.
//
// This function may panic given code that doesn't pass the static verifier.
func (i *Interpreter) Run() (err error) {
	// Deliberately implementing the entire core in a single function here
	// to give the compiler more creative liberties.

	var r [11]uint64
	r[1] = VaddrInput
	// TODO frame pointer
	pc := int64(i.entry)

	// TODO step to next instruction

	for {
		// Fetch
		ins := i.getSlot(pc)
		// Execute
		pc++
		switch ins.Op() {
		case OpAdd32Imm:
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) + ins.Imm())
		case OpAdd32Reg:
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) + int32(r[ins.Src()]))
		case OpAdd64Imm:
			r[ins.Dst()] += uint64(ins.Imm())
		case OpAdd64Reg:
			r[ins.Dst()] += r[ins.Src()]
		case OpSub32Imm:
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) - ins.Imm())
		case OpSub32Reg:
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) - int32(r[ins.Src()]))
		case OpSub64Imm:
			r[ins.Dst()] -= uint64(ins.Imm())
		case OpSub64Reg:
			r[ins.Dst()] -= r[ins.Src()]
		case OpMul32Imm:
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) * ins.Imm())
		case OpMul32Reg:
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) * int32(r[ins.Src()]))
		case OpMul64Imm:
			r[ins.Dst()] *= uint64(ins.Imm())
		case OpMul64Reg:
			r[ins.Dst()] *= r[ins.Src()]
		case OpDiv32Imm:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) / ins.Uimm())
		case OpDiv32Reg:
			if src := uint32(r[ins.Src()]); src != 0 {
				r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) / src)
			} else {
				return ExcDivideByZero
			}
		case OpDiv64Imm:
			r[ins.Dst()] /= uint64(ins.Imm())
		case OpDiv64Reg:
			if src := r[ins.Src()]; src != 0 {
				r[ins.Dst()] /= src
			} else {
				return ExcDivideByZero
			}
		case OpSdiv32Imm:
			if int32(r[ins.Dst()]) == math.MinInt32 && ins.Imm() == -1 {
				return ExcDivideOverflow
			}
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) / ins.Imm())
		case OpSdiv32Reg:
			if src := int32(r[ins.Src()]); src != 0 {
				if int32(r[ins.Dst()]) == math.MinInt32 && src == -1 {
					return ExcDivideOverflow
				}
				r[ins.Dst()] = uint64(int32(r[ins.Dst()]) / src)
			} else {
				return ExcDivideByZero
			}
		case OpSdiv64Imm:
			if int64(r[ins.Dst()]) == math.MinInt64 && ins.Imm() == -1 {
				return ExcDivideOverflow
			}
			r[ins.Dst()] = uint64(int64(r[ins.Dst()]) / int64(ins.Imm()))
		case OpSdiv64Reg:
			if src := int64(r[ins.Src()]); src != 0 {
				if int64(r[ins.Dst()]) == math.MinInt64 && src == -1 {
					return ExcDivideOverflow
				}
				r[ins.Dst()] = uint64(int64(r[ins.Dst()]) / src)
			} else {
				return ExcDivideByZero
			}
		case OpOr32Imm:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) | ins.Uimm())
		case OpOr32Reg:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) | uint32(r[ins.Src()]))
		case OpOr64Imm:
			r[ins.Dst()] |= uint64(ins.Imm())
		case OpOr64Reg:
			r[ins.Dst()] |= r[ins.Src()]
		case OpAnd32Imm:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) & ins.Uimm())
		case OpAnd32Reg:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) & uint32(r[ins.Src()]))
		case OpAnd64Imm:
			r[ins.Dst()] &= uint64(ins.Imm())
		case OpAnd64Reg:
			r[ins.Dst()] &= r[ins.Src()]
		case OpLsh32Imm:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) << ins.Uimm())
		case OpLsh32Reg:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) << uint32(r[ins.Src()]&0x1f))
		case OpLsh64Imm:
			r[ins.Dst()] <<= uint64(ins.Imm())
		case OpLsh64Reg:
			r[ins.Dst()] <<= r[ins.Src()] & 0x3f
		case OpRsh32Imm:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) >> ins.Uimm())
		case OpRsh32Reg:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) >> uint32(r[ins.Src()]&0x1f))
		case OpRsh64Imm:
			r[ins.Dst()] >>= uint64(ins.Imm())
		case OpRsh64Reg:
			r[ins.Dst()] >>= r[ins.Src()] & 0x3f
		case OpNeg32:
			r[ins.Dst()] = uint64(-int32(r[ins.Dst()]))
		case OpNeg64:
			r[ins.Dst()] = uint64(-int64(r[ins.Dst()]))
		case OpMod32Imm:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) % ins.Uimm())
		case OpMod32Reg:
			if src := uint32(r[ins.Src()]); src != 0 {
				r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) % src)
			} else {
				return ExcDivideByZero
			}
		case OpMod64Imm:
			r[ins.Dst()] %= uint64(ins.Imm())
		case OpMod64Reg:
			if src := r[ins.Src()]; src != 0 {
				r[ins.Dst()] %= src
			} else {
				return ExcDivideByZero
			}
		case OpXor32Imm:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) ^ ins.Uimm())
		case OpXor32Reg:
			r[ins.Dst()] = uint64(uint32(r[ins.Dst()]) ^ uint32(r[ins.Src()]))
		case OpXor64Imm:
			r[ins.Dst()] ^= uint64(ins.Imm())
		case OpXor64Reg:
			r[ins.Dst()] ^= r[ins.Src()]
		case OpMov32Imm:
			r[ins.Dst()] = uint64(ins.Uimm())
		case OpMov32Reg:
			r[ins.Dst()] = r[ins.Src()] & math.MaxUint32
		case OpMov64Imm:
			r[ins.Dst()] = uint64(ins.Imm())
		case OpMov64Reg:
			r[ins.Dst()] = r[ins.Src()]
		case OpArsh32Imm:
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) >> ins.Uimm())
		case OpArsh32Reg:
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) >> uint32(r[ins.Src()]&0x1f))
		case OpArsh64Imm:
			r[ins.Dst()] = uint64(int64(r[ins.Dst()]) >> ins.Imm())
		case OpArsh64Reg:
			r[ins.Dst()] = uint64(int64(r[ins.Dst()]) >> (r[ins.Src()] & 0x3f))
		case OpLe:
			switch ins.Uimm() {
			case 16:
				r[ins.Dst()] &= math.MaxUint16
			case 32:
				r[ins.Dst()] &= math.MaxUint32
			case 64:
				r[ins.Dst()] &= math.MaxUint64
			default:
				panic("invalid le instruction")
			}
		case OpBe:
			switch ins.Uimm() {
			case 16:
				r[ins.Dst()] = uint64(bits.ReverseBytes16(uint16(r[ins.Dst()])))
			case 32:
				r[ins.Dst()] = uint64(bits.ReverseBytes32(uint32(r[ins.Dst()])))
			case 64:
				r[ins.Dst()] = bits.ReverseBytes64(r[ins.Dst()])
			default:
				panic("invalid be instruction")
			}
		case OpLddw:
			r[ins.Dst()] = uint64(ins.Uimm()) | (uint64(i.getSlot(pc+1).Uimm()) << 32)
			pc++
		case OpJa:
			pc += int64(ins.Off())
		case OpJeqImm:
			if r[ins.Dst()] == uint64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJeqReg:
			if r[ins.Dst()] == r[ins.Src()] {
				pc += int64(ins.Off())
			}
		case OpJgtImm:
			if r[ins.Dst()] > uint64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJgtReg:
			if r[ins.Dst()] > r[ins.Src()] {
				pc += int64(ins.Off())
			}
		case OpJgeImm:
			if r[ins.Dst()] >= uint64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJgeReg:
			if r[ins.Dst()] >= r[ins.Src()] {
				pc += int64(ins.Off())
			}
		case OpJltImm:
			if r[ins.Dst()] < uint64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJltReg:
			if r[ins.Dst()] < r[ins.Src()] {
				pc += int64(ins.Off())
			}
		case OpJleImm:
			if r[ins.Dst()] <= uint64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJleReg:
			if r[ins.Dst()] <= r[ins.Src()] {
				pc += int64(ins.Off())
			}
		case OpJsetImm:
			if r[ins.Dst()]&uint64(ins.Imm()) != 0 {
				pc += int64(ins.Off())
			}
		case OpJsetReg:
			if r[ins.Dst()]&r[ins.Src()] != 0 {
				pc += int64(ins.Off())
			}
		case OpJneImm:
			if r[ins.Dst()] != uint64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJneReg:
			if r[ins.Dst()] != r[ins.Src()] {
				pc += int64(ins.Off())
			}
		case OpJsgtImm:
			if int64(r[ins.Dst()]) > int64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJsgtReg:
			if int64(r[ins.Dst()]) > int64(r[ins.Src()]) {
				pc += int64(ins.Off())
			}
		case OpJsgeImm:
			if int64(r[ins.Dst()]) >= int64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJsgeReg:
			if int64(r[ins.Dst()]) >= int64(r[ins.Src()]) {
				pc += int64(ins.Off())
			}
		case OpJsltImm:
			if int64(r[ins.Dst()]) < int64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJsltReg:
			if int64(r[ins.Dst()]) < int64(r[ins.Src()]) {
				pc += int64(ins.Off())
			}
		case OpJsleImm:
			if int64(r[ins.Dst()]) <= int64(ins.Imm()) {
				pc += int64(ins.Off())
			}
		case OpJsleReg:
			if int64(r[ins.Dst()]) <= int64(r[ins.Src()]) {
				pc += int64(ins.Off())
			}
		case OpCall:
			// TODO use src reg hint
			if sc, ok := i.syscalls[ins.Uimm()]; ok {
				r[0], err = sc.Invoke(i, r[1], r[2], r[3], r[4], r[5])
			} else {
				panic("bpf function calls not implemented")
			}
		case OpCallx:
			panic("callx not implemented")
		case OpExit:
			return nil
		default:
			panic(fmt.Sprintf("unimplemented opcode %#02x", ins.Op()))
		}
	}
}

func (i *Interpreter) getSlot(pc int64) Slot {
	return GetSlot(i.text[pc*SlotSize:])
}

func (i *Interpreter) VMContext() any {
	return i.vmContext
}
