package sbf

import (
	"fmt"
	"math"
	"math/bits"
	"unsafe"
)

// Interpreter implements the SBF core in pure Go.
type Interpreter struct {
	text  []byte
	ro    []byte
	stack []byte
	heap  []byte
	input []byte

	entry uint64

	cuMax uint64

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
	cuLeft := int64(i.cuMax)

	// TODO step to next instruction

mainLoop:
	for {
		// Fetch
		ins := i.getSlot(pc)
		// Execute
		switch ins.Op() {
		case OpLdxb:
			vma := uint64(int64(r[ins.Src()]) + int64(ins.Off()))
			var v uint8
			v, err = i.Read8(vma)
			r[ins.Dst()] = uint64(v)
		case OpLdxh:
			vma := uint64(int64(r[ins.Src()]) + int64(ins.Off()))
			var v uint16
			v, err = i.Read16(vma)
			r[ins.Dst()] = uint64(v)
		case OpLdxw:
			vma := uint64(int64(r[ins.Src()]) + int64(ins.Off()))
			var v uint32
			v, err = i.Read32(vma)
			r[ins.Dst()] = uint64(v)
		case OpLdxdw:
			vma := uint64(int64(r[ins.Src()]) + int64(ins.Off()))
			r[ins.Dst()], err = i.Read64(vma)
		case OpStb:
			vma := uint64(int64(r[ins.Dst()]) + int64(ins.Off()))
			err = i.Write8(vma, uint8(ins.Uimm()))
		case OpSth:
			vma := uint64(int64(r[ins.Dst()]) + int64(ins.Off()))
			err = i.Write16(vma, uint16(ins.Uimm()))
		case OpStw:
			vma := uint64(int64(r[ins.Dst()]) + int64(ins.Off()))
			err = i.Write32(vma, ins.Uimm())
		case OpStdw:
			vma := uint64(int64(r[ins.Dst()]) + int64(ins.Off()))
			err = i.Write64(vma, uint64(ins.Imm()))
		case OpStxb:
			vma := uint64(int64(r[ins.Dst()]) + int64(ins.Off()))
			err = i.Write8(vma, uint8(r[ins.Src()]))
		case OpStxh:
			vma := uint64(int64(r[ins.Dst()]) + int64(ins.Off()))
			err = i.Write16(vma, uint16(r[ins.Src()]))
		case OpStxw:
			vma := uint64(int64(r[ins.Dst()]) + int64(ins.Off()))
			err = i.Write32(vma, uint32(r[ins.Src()]))
		case OpStxdw:
			vma := uint64(int64(r[ins.Dst()]) + int64(ins.Off()))
			err = i.Write64(vma, r[ins.Src()])
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
				err = ExcDivideByZero
			}
		case OpSdiv32Imm:
			if int32(r[ins.Dst()]) == math.MinInt32 && ins.Imm() == -1 {
				err = ExcDivideOverflow
			}
			r[ins.Dst()] = uint64(int32(r[ins.Dst()]) / ins.Imm())
		case OpSdiv32Reg:
			if src := int32(r[ins.Src()]); src != 0 {
				if int32(r[ins.Dst()]) == math.MinInt32 && src == -1 {
					err = ExcDivideOverflow
				}
				r[ins.Dst()] = uint64(int32(r[ins.Dst()]) / src)
			} else {
				err = ExcDivideByZero
			}
		case OpSdiv64Imm:
			if int64(r[ins.Dst()]) == math.MinInt64 && ins.Imm() == -1 {
				err = ExcDivideOverflow
			}
			r[ins.Dst()] = uint64(int64(r[ins.Dst()]) / int64(ins.Imm()))
		case OpSdiv64Reg:
			if src := int64(r[ins.Src()]); src != 0 {
				if int64(r[ins.Dst()]) == math.MinInt64 && src == -1 {
					err = ExcDivideOverflow
				}
				r[ins.Dst()] = uint64(int64(r[ins.Dst()]) / src)
			} else {
				err = ExcDivideByZero
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
				err = ExcDivideByZero
			}
		case OpMod64Imm:
			r[ins.Dst()] %= uint64(ins.Imm())
		case OpMod64Reg:
			if src := r[ins.Src()]; src != 0 {
				r[ins.Dst()] %= src
			} else {
				err = ExcDivideByZero
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
				r[0], cuLeft, err = sc.Invoke(i, r[1], r[2], r[3], r[4], r[5], cuLeft)
			} else {
				panic("bpf function calls not implemented")
			}
		case OpCallx:
			panic("callx not implemented")
		case OpExit:
			// TODO implement function returns
			break mainLoop
		default:
			panic(fmt.Sprintf("unimplemented opcode %#02x", ins.Op()))
		}
		// Post execute
		if err != nil {
			// TODO return CPU exception error type here
			return err
		}
		pc++
	}

	return nil
}

func (i *Interpreter) getSlot(pc int64) Slot {
	return GetSlot(i.text[pc*SlotSize:])
}

func (i *Interpreter) VMContext() any {
	return i.vmContext
}

func (i *Interpreter) Translate(addr uint64, size uint32, write bool) (unsafe.Pointer, error) {
	// TODO exhaustive testing against rbpf
	// TODO review generated asm for performance

	hi, lo := addr>>32, addr&math.MaxUint32
	switch hi {
	case VaddrProgram >> 32:
		if write {
			return nil, NewExcBadAccess(addr, size, write, "write to program")
		}
		if lo+uint64(size) >= uint64(len(i.ro)) {
			return nil, NewExcBadAccess(addr, size, write, "out-of-bounds program read")
		}
		return unsafe.Pointer(&i.ro[lo]), nil
	case VaddrStack >> 32:
		panic("todo implement stack access check")
	case VaddrHeap >> 32:
		panic("todo implement heap access check")
	case VaddrInput >> 32:
		if lo+uint64(size) >= uint64(len(i.input)) {
			return nil, NewExcBadAccess(addr, size, write, "out-of-bounds input read")
		}
		return unsafe.Pointer(&i.input[lo]), nil
	default:
		return nil, NewExcBadAccess(addr, size, write, "unmapped region")
	}
}

func (i *Interpreter) Read(addr uint64, p []byte) error {
	ptr, err := i.Translate(addr, uint32(len(p)), false)
	if err != nil {
		return err
	}
	mem := unsafe.Slice((*uint8)(ptr), len(p))
	copy(p, mem)
	return nil
}

func (i *Interpreter) Read8(addr uint64) (uint8, error) {
	ptr, err := i.Translate(addr, 1, false)
	if err != nil {
		return 0, err
	}
	return *(*uint8)(ptr), nil
}

// TODO is it safe and portable to deref unaligned integer types?

func (i *Interpreter) Read16(addr uint64) (uint16, error) {
	ptr, err := i.Translate(addr, 2, false)
	if err != nil {
		return 0, err
	}
	return *(*uint16)(ptr), nil
}

func (i *Interpreter) Read32(addr uint64) (uint32, error) {
	ptr, err := i.Translate(addr, 4, false)
	if err != nil {
		return 0, err
	}
	return *(*uint32)(ptr), nil
}

func (i *Interpreter) Read64(addr uint64) (uint64, error) {
	ptr, err := i.Translate(addr, 8, false)
	if err != nil {
		return 0, err
	}
	return *(*uint64)(ptr), nil
}

func (i *Interpreter) Write(addr uint64, p []byte) error {
	ptr, err := i.Translate(addr, uint32(len(p)), true)
	if err != nil {
		return err
	}
	mem := unsafe.Slice((*uint8)(ptr), len(p))
	copy(mem, p)
	return nil
}

func (i *Interpreter) Write8(addr uint64, x uint8) error {
	ptr, err := i.Translate(addr, 1, true)
	if err != nil {
		return err
	}
	*(*uint8)(ptr) = x
	return nil
}

func (i *Interpreter) Write16(addr uint64, x uint16) error {
	ptr, err := i.Translate(addr, 2, true)
	if err != nil {
		return err
	}
	*(*uint16)(ptr) = x
	return nil
}

func (i *Interpreter) Write32(addr uint64, x uint32) error {
	ptr, err := i.Translate(addr, 4, true)
	if err != nil {
		return err
	}
	*(*uint32)(ptr) = x
	return nil
}

func (i *Interpreter) Write64(addr uint64, x uint64) error {
	ptr, err := i.Translate(addr, 8, false)
	if err != nil {
		return err
	}
	*(*uint64)(ptr) = x
	return nil
}
