package insbuilder

type Instruction interface {
	// returns instruction opt code
	OptCodeByte() uint8
	// returns destination register
	GetDst() uint8
	// returns source register
	GetSrc() uint8
	// returns offset bytes
	GetOff() int16
	// returns immediate value
	GetImm() int64
	// sets destination register
	SetDst(dst uint8) Instruction
	// sets source register
	SetSrc(src uint8) Instruction
	// sets offset bytes
	SetOff(offset int16) Instruction
	// sets immediate value
	SetImm(imm int64) Instruction
	// get `ebpf::Insn` struct
	GetInsn() *Insn
	// get mutable `ebpf::Insn` struct
	// GetInsnMut() *Insn // TODO: remove this method

	Push() *BpfCode
	BytesInterface
}

type BytesInterface interface {
	Bytes() []byte
}

// General implementation of `IntoBytes` for `Instruction`
func IntoBytes(i Instruction) []byte {
	/// transform immutable reference of `Instruction` into `Vec<u8>` with size of 8
	/// [ 1 byte ,      1 byte      , 2 bytes,  4 bytes  ]
	/// [ OP_CODE, SRC_REG | DST_REG, OFFSET , IMMEDIATE ]

	bytes := make([]byte, 8)
	bytes[0] = i.OptCodeByte()
	bytes[1] = i.GetSrc()<<4 | i.GetDst()

	bytes[2] = uint8(i.GetOff())
	bytes[3] = uint8(i.GetOff() >> 8)

	bytes[4] = uint8(i.GetImm())
	bytes[5] = uint8(i.GetImm() >> 8)
	bytes[6] = uint8(i.GetImm() >> 16)
	bytes[7] = uint8(i.GetImm() >> 24)
	return bytes
}

// BPF instruction stack in byte representation
type BpfCode struct {
	instructions []byte
}

// / The source of ALU and JMP instructions
type Source int32

const (
	SourceImm Source = Source(BPF_IMM)
	SourceReg Source = Source(BPF_X)
)

type OpBits int

const (
	OpBitsAdd        OpBits = OpBits(BPF_ADD)
	OpBitsSub        OpBits = OpBits(BPF_SUB)
	OpBitsMul        OpBits = OpBits(BPF_MUL)
	OpBitsDiv        OpBits = OpBits(BPF_DIV)
	OpBitsBitOr      OpBits = OpBits(BPF_OR)
	OpBitsBitAnd     OpBits = OpBits(BPF_AND)
	OpBitsLShift     OpBits = OpBits(BPF_LSH)
	OpBitsRShift     OpBits = OpBits(BPF_RSH)
	OpBitsNegate     OpBits = OpBits(BPF_NEG)
	OpBitsMod        OpBits = OpBits(BPF_MOD)
	OpBitsBitXor     OpBits = OpBits(BPF_XOR)
	OpBitsMov        OpBits = OpBits(BPF_MOV)
	OpBitsSignRShift OpBits = OpBits(BPF_ARSH)
)

// / Architecture of instructions
type Arch int

const (
	ArchX64 Arch = Arch(BPF_ALU64)
	ArchX32 Arch = Arch(BPF_ALU)
)

type Addressing int

const (
	AddressingImm Addressing = Addressing(BPF_IMM)
	AddressingAbs Addressing = Addressing(BPF_ABS)
	AddressingInd Addressing = Addressing(BPF_IND)
	AddressingMem Addressing = Addressing(BPF_MEM)
)

// / Memory size for LOAD and STORE instructions
type MemSize int

const (
	// 8-bit size
	MemSizeByte MemSize = MemSize(BPF_B)
	// 16-bit size
	MemSizeHalfWord MemSize = MemSize(BPF_H)
	// 32-bit size
	MemSizeWord MemSize = MemSize(BPF_W)
	// 64-bit size
	MemSizeDoubleWord MemSize = MemSize(BPF_DW)
)

// / Bytes endian
type Endian int

const (
	// Little endian
	EndianLittle Endian = Endian(LE)
	// Big endian
	EndianBig Endian = Endian(BE)
)

// / Conditions for JMP instructions
type Cond int

const (
	// Absolute or unconditional
	CondAbs Cond = Cond(BPF_JA)
	// Jump if `==`
	CondEquals Cond = Cond(BPF_JEQ)
	// Jump if `>`
	CondGreater Cond = Cond(BPF_JGT)
	// Jump if `>=`
	CondGreaterEquals Cond = Cond(BPF_JGE)
	// Jump if `<`
	CondLower Cond = Cond(BPF_JLT)
	// Jump if `<=`
	CondLowerEquals Cond = Cond(BPF_JLE)
	// Jump if `src` & `dst`
	CondBitAnd Cond = Cond(BPF_JSET)
	// Jump if `!=`
	CondNotEquals Cond = Cond(BPF_JNE)
	// Jump if `>` (signed)
	CondGreaterSigned Cond = Cond(BPF_JSGT)
	// Jump if `>=` (signed)
	CondGreaterEqualsSigned Cond = Cond(BPF_JSGE)
	// Jump if `<` (signed)
	CondLowerSigned Cond = Cond(BPF_JSLT)
	// Jump if `<=` (signed)
	CondLowerEqualsSigned Cond = Cond(BPF_JSLE)
)

func NewBpfCode() *BpfCode {
	return &BpfCode{
		instructions: make([]byte, 0),
	}
}

// Bytes returns byte representation of BPF instructions
func (b *BpfCode) Bytes() []byte {
	return b.instructions
}

// impl BpfCode {
//     /// creates new empty BPF instruction stack
//     pub fn new() -> Self {
//         BpfCode {
//             instructions: vec![],
//         }
//     }

//     /// create ADD instruction
//     pub fn add(&mut self, source: Source, arch: Arch) -> Move {
//         self.mov_internal(source, arch, OpBits::Add)
//     }

// /// create ADD instruction
func (b *BpfCode) Add(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsAdd)
}

//     /// create SUB instruction
//     pub fn sub(&mut self, source: Source, arch: Arch) -> Move {
//         self.mov_internal(source, arch, OpBits::Sub)
//     }

// /// create SUB instruction
func (b *BpfCode) Sub(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsSub)
}

// /// create MUL instruction
//
//	pub fn mul(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::Mul)
//	}
//
// /// create MUL instruction
func (b *BpfCode) Mul(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsMul)
}

// /// create DIV instruction
//
//	pub fn div(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::Div)
//	}
//
// /// create DIV instruction
func (b *BpfCode) Div(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsDiv)
}

// /// create OR instruction
//
//	pub fn bit_or(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::BitOr)
//	}
//
// /// create OR instruction
func (b *BpfCode) BitOr(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsBitOr)
}

// /// create AND instruction
//
//	pub fn bit_and(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::BitAnd)
//	}
//
// /// create AND instruction
func (b *BpfCode) BitAnd(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsBitAnd)
}

// /// create LSHIFT instruction
//
//	pub fn left_shift(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::LShift)
//	}
//
// /// create LSHIFT instruction
func (b *BpfCode) LeftShift(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsLShift)
}

// /// create RSHIFT instruction
//
//	pub fn right_shift(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::RShift)
//	}
//
// /// create RSHIFT instruction
func (b *BpfCode) RightShift(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsRShift)
}

// /// create NEGATE instruction
//
//	pub fn negate(&mut self, arch: Arch) -> Move {
//	    self.mov_internal(Source::Imm, arch, OpBits::Negate)
//	}
//
// /// create NEGATE instruction
func (b *BpfCode) Negate(arch Arch) *Move {
	return b.mov_internal(SourceImm, arch, OpBitsNegate)
}

// /// create MOD instruction
//
//	pub fn modulo(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::Mod)
//	}
//
// /// create MOD instruction
func (b *BpfCode) Modulo(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsMod)
}

// /// create XOR instruction
//
//	pub fn bit_xor(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::BitXor)
//	}
//
// /// create XOR instruction
func (b *BpfCode) BitXor(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsBitXor)
}

// /// create MOV instruction
//
//	pub fn mov(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::Mov)
//	}
//
// /// create MOV instruction
func (b *BpfCode) Mov(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsMov)
}

// /// create SIGNED RSHIFT instruction
//
//	pub fn signed_right_shift(&mut self, source: Source, arch: Arch) -> Move {
//	    self.mov_internal(source, arch, OpBits::SignRShift)
//	}
//
// /// create SIGNED RSHIFT instruction
func (b *BpfCode) SignedRightShift(source Source, arch Arch) *Move {
	return b.mov_internal(source, arch, OpBitsSignRShift)
}

// #[inline]
//
//	fn mov_internal(&mut self, source: Source, arch_bits: Arch, op_bits: OpBits) -> Move {
//	    Move {
//	        bpf_code: self,
//	        src_bit: source,
//	        op_bits,
//	        arch_bits,
//	        insn: Insn::default(),
//	    }
//	}
func (b *BpfCode) mov_internal(source Source, arch_bits Arch, op_bits OpBits) *Move {
	return &Move{
		bpfCode:  b,
		srcBit:   source,
		opBits:   op_bits,
		archBits: arch_bits,
		insn:     Insn{},
	}
}

// /// create byte swap instruction
//
//	pub fn swap_bytes(&mut self, endian: Endian) -> SwapBytes {
//	    SwapBytes {
//	        bpf_code: self,
//	        endian,
//	        insn: Insn::default(),
//	    }
//	}
//
// /// create byte swap instruction
func (b *BpfCode) SwapBytes(endian Endian) *SwapBytes {
	return &SwapBytes{
		bpfCode: b,
		endian:  endian,
		insn:    Insn{},
	}
}

// /// create LOAD instruction, IMMEDIATE is the source
//
//	pub fn load(&mut self, mem_size: MemSize) -> Load {
//	    self.load_internal(mem_size, Addressing::Imm, BPF_LD)
//	}
//
// /// create LOAD instruction, IMMEDIATE is the source
func (b *BpfCode) Load(mem_size MemSize) *Load_ {
	return b.load_internal(mem_size, AddressingImm, BPF_LD)
}

// /// create ABSOLUTE LOAD instruction
//
//	pub fn load_abs(&mut self, mem_size: MemSize) -> Load {
//	    self.load_internal(mem_size, Addressing::Abs, BPF_LD)
//	}
//
// /// create ABSOLUTE LOAD instruction
func (b *BpfCode) LoadAbs(mem_size MemSize) *Load_ {
	return b.load_internal(mem_size, AddressingAbs, BPF_LD)
}

// /// create INDIRECT LOAD instruction
//
//	pub fn load_ind(&mut self, mem_size: MemSize) -> Load {
//	    self.load_internal(mem_size, Addressing::Ind, BPF_LD)
//	}
//
// /// create INDIRECT LOAD instruction
func (b *BpfCode) LoadInd(mem_size MemSize) *Load_ {
	return b.load_internal(mem_size, AddressingInd, BPF_LD)
}

// /// create LOAD instruction, MEMORY is the source
//
//	pub fn load_x(&mut self, mem_size: MemSize) -> Load {
//	    self.load_internal(mem_size, Addressing::Mem, BPF_LDX)
//	}
//
// /// create LOAD instruction, MEMORY is the source
func (b *BpfCode) LoadX(mem_size MemSize) *Load_ {
	return b.load_internal(mem_size, AddressingMem, BPF_LDX)
}

// #[inline]
//
//	fn load_internal(&mut self, mem_size: MemSize, addressing: Addressing, source: u8) -> Load {
//	    Load {
//	        bpf_code: self,
//	        addressing,
//	        mem_size,
//	        source,
//	        insn: Insn::default(),
//	    }
//	}
func (b *BpfCode) load_internal(mem_size MemSize, addressing Addressing, source uint8) *Load_ {
	return &Load_{
		bpfCode:    b,
		addressing: addressing,
		memSize:    mem_size,
		source:     source,
		insn:       Insn{},
	}
}

// /// creates STORE instruction, IMMEDIATE is the source
//
//	pub fn store(&mut self, mem_size: MemSize) -> Store {
//	    self.store_internal(mem_size, BPF_IMM)
//	}
//
// /// creates STORE instruction, IMMEDIATE is the source
func (b *BpfCode) Store(mem_size MemSize) *Store_ {
	return b.store_internal(mem_size, BPF_IMM)
}

// /// creates STORE instruction, MEMORY is the source
//
//	pub fn store_x(&mut self, mem_size: MemSize) -> Store {
//	    self.store_internal(mem_size, BPF_MEM | BPF_STX)
//	}
//
// /// creates STORE instruction, MEMORY is the source
func (b *BpfCode) StoreX(mem_size MemSize) *Store_ {
	return b.store_internal(mem_size, BPF_MEM|BPF_STX)
}

// #[inline]
//
//	fn store_internal(&mut self, mem_size: MemSize, source: u8) -> Store {
//	    Store {
//	        bpf_code: self,
//	        mem_size,
//	        source,
//	        insn: Insn::default(),
//	    }
//	}
func (b *BpfCode) store_internal(mem_size MemSize, source uint8) *Store_ {
	return &Store_{
		bpfCode: b,
		memSize: mem_size,
		source:  source,
		insn:    Insn{},
	}
}

// /// create unconditional JMP instruction
//
//	pub fn jump_unconditional(&mut self) -> Jump {
//	    self.jump_conditional(Cond::Abs, Source::Imm)
//	}
//
// /// create unconditional JMP instruction
func (b *BpfCode) JumpUnconditional() *Jump_ {
	return b.JumpConditional(CondAbs, SourceImm)
}

// /// create conditional JMP instruction
//
//	pub fn jump_conditional(&mut self, cond: Cond, src_bit: Source) -> Jump {
//	    Jump {
//	        bpf_code: self,
//	        cond,
//	        src_bit,
//	        insn: Insn::default(),
//	    }
//	}
//
// /// create conditional JMP instruction
func (b *BpfCode) JumpConditional(cond Cond, src_bit Source) *Jump_ {
	return &Jump_{
		bpfCode: b,
		cond:    cond,
		srcBit:  src_bit,
		insn:    Insn{},
	}
}

// /// create CALL instruction
//
//	pub fn call(&mut self) -> FunctionCall {
//	    FunctionCall {
//	        bpf_code: self,
//	        insn: Insn::default(),
//	    }
//	}
//
// /// create CALL instruction
func (b *BpfCode) Call() *FunctionCall {
	return &FunctionCall{
		bpfCode: b,
		insn:    Insn{},
	}
}

// /// create EXIT instruction
//
//	pub fn exit(&mut self) -> Exit {
//	    Exit {
//	        bpf_code: self,
//	        insn: Insn::default(),
//	    }
//	}
//
// /// create EXIT instruction
func (b *BpfCode) Exit() *Exit {
	return &Exit{
		bpfCode: b,
		insn:    Insn{},
	}
}

// }

var (
	_ Instruction = (*Move)(nil)
	_ Instruction = (*SwapBytes)(nil)
	_ Instruction = (*Load_)(nil)
	_ Instruction = (*Store_)(nil)
	_ Instruction = (*Jump_)(nil)
	_ Instruction = (*FunctionCall)(nil)
	_ Instruction = (*Exit)(nil)
)

// struct to represent `MOV ALU` instructions
type Move struct {
	bpfCode  *BpfCode
	srcBit   Source
	opBits   OpBits
	archBits Arch
	insn     Insn
}

func (m *Move) OptCodeByte() byte {
	opBits := byte(m.opBits)
	srcBit := byte(m.srcBit)
	archBits := byte(m.archBits)
	return opBits | srcBit | archBits
}

func (m *Move) Push() *BpfCode {
	asm := m.Bytes()
	m.bpfCode.instructions = append(m.bpfCode.instructions, asm...)
	return m.bpfCode
}

func (m *Move) Bytes() []byte {
	i := Instruction(m)
	return IntoBytes(i)
}

func (m *Move) GetInsn() *Insn {
	return &m.insn
}

func (i *Move) GetSrc() uint8 {
	return i.insn.Src
}

func (i *Move) GetDst() uint8 {
	return i.insn.Dst
}

func (i *Move) GetOff() int16 {
	return i.insn.Off
}

func (i *Move) GetImm() int64 {
	return i.insn.Imm
}

func (i *Move) SetDst(dst uint8) Instruction {
	i.insn.Dst = dst
	return i
}

func (i *Move) SetSrc(src uint8) Instruction {
	i.insn.Src = src
	return i
}

func (i *Move) SetOff(off int16) Instruction {
	i.insn.Off = off
	return i
}

func (i *Move) SetImm(imm int64) Instruction {
	i.insn.Imm = imm
	return i
}

////////////////////////////////////////////////////////////////

// struct representation of byte swap operation
type SwapBytes struct {
	bpfCode *BpfCode
	endian  Endian
	insn    Insn
}

func (sb *SwapBytes) OptCodeByte() uint8 {
	return uint8(sb.endian)
}

func (m *SwapBytes) Push() *BpfCode {
	asm := m.Bytes()
	m.bpfCode.instructions = append(m.bpfCode.instructions, asm...)
	return m.bpfCode
}

func (m *SwapBytes) Bytes() []byte {
	i := Instruction(m)
	return IntoBytes(i)
}

func (m *SwapBytes) GetInsn() *Insn {
	return &m.insn
}

func (i *SwapBytes) GetSrc() uint8 {
	return i.insn.Src
}

func (i *SwapBytes) GetDst() uint8 {
	return i.insn.Dst
}

func (i *SwapBytes) GetOff() int16 {
	return i.insn.Off
}

func (i *SwapBytes) GetImm() int64 {
	return i.insn.Imm
}

func (i *SwapBytes) SetDst(dst uint8) Instruction {
	i.insn.Dst = dst
	return i
}

func (i *SwapBytes) SetSrc(src uint8) Instruction {
	i.insn.Src = src
	return i
}

func (i *SwapBytes) SetOff(off int16) Instruction {
	i.insn.Off = off
	return i
}

func (i *SwapBytes) SetImm(imm int64) Instruction {
	i.insn.Imm = imm
	return i
}

////////////////////////////////////////////////////////////////

// / struct representation of LOAD instructions
type Load_ struct {
	bpfCode    *BpfCode
	addressing Addressing
	memSize    MemSize
	source     uint8
	insn       Insn
}

func (l *Load_) OptCodeByte() uint8 {
	size := uint8(l.memSize)
	addressing := uint8(l.addressing)
	return addressing | size | l.source
}

func (m *Load_) Push() *BpfCode {
	asm := m.Bytes()
	m.bpfCode.instructions = append(m.bpfCode.instructions, asm...)
	return m.bpfCode
}

func (m *Load_) Bytes() []byte {
	i := Instruction(m)
	return IntoBytes(i)
}

func (m *Load_) GetInsn() *Insn {
	return &m.insn
}

func (i *Load_) GetSrc() uint8 {
	return i.insn.Src
}

func (i *Load_) GetDst() uint8 {
	return i.insn.Dst
}

func (i *Load_) GetOff() int16 {
	return i.insn.Off
}

func (i *Load_) GetImm() int64 {
	return i.insn.Imm
}

func (i *Load_) SetDst(dst uint8) Instruction {
	i.insn.Dst = dst
	return i
}

func (i *Load_) SetSrc(src uint8) Instruction {
	i.insn.Src = src
	return i
}

func (i *Load_) SetOff(off int16) Instruction {
	i.insn.Off = off
	return i
}

func (i *Load_) SetImm(imm int64) Instruction {
	i.insn.Imm = imm
	return i
}

// //////////////////////////////////////////////////////////////

// / struct representation of STORE instructions
type Store_ struct {
	bpfCode *BpfCode
	memSize MemSize
	source  byte
	insn    Insn
}

func (s *Store_) OptCodeByte() byte {
	size := uint8(s.memSize)
	return BPF_MEM | BPF_ST | size | s.source
}

func (m *Store_) Push() *BpfCode {
	asm := m.Bytes()
	m.bpfCode.instructions = append(m.bpfCode.instructions, asm...)
	return m.bpfCode
}

func (m *Store_) Bytes() []byte {
	i := Instruction(m)
	return IntoBytes(i)
}

func (m *Store_) GetInsn() *Insn {
	return &m.insn
}

func (i *Store_) GetSrc() uint8 {
	return i.insn.Src
}

func (i *Store_) GetDst() uint8 {
	return i.insn.Dst
}

func (i *Store_) GetOff() int16 {
	return i.insn.Off
}

func (i *Store_) GetImm() int64 {
	return i.insn.Imm
}

func (i *Store_) SetDst(dst uint8) Instruction {
	i.insn.Dst = dst
	return i
}

func (i *Store_) SetSrc(src uint8) Instruction {
	i.insn.Src = src
	return i
}

func (i *Store_) SetOff(off int16) Instruction {
	i.insn.Off = off
	return i
}

func (i *Store_) SetImm(imm int64) Instruction {
	i.insn.Imm = imm
	return i
}

// //////////////////////////////////////////////////////////////

// / struct representation of JMP instructions
type Jump_ struct {
	bpfCode *BpfCode
	cond    Cond
	srcBit  Source
	insn    Insn
}

func (j *Jump_) OptCodeByte() uint8 {
	cmp := uint8(j.cond)
	srcBit := uint8(j.srcBit)
	return cmp | srcBit | BPF_JMP
}

func (m *Jump_) Push() *BpfCode {
	asm := m.Bytes()
	m.bpfCode.instructions = append(m.bpfCode.instructions, asm...)
	return m.bpfCode
}

func (m *Jump_) Bytes() []byte {
	i := Instruction(m)
	return IntoBytes(i)
}

func (m *Jump_) GetInsn() *Insn {
	return &m.insn
}

func (i *Jump_) GetSrc() uint8 {
	return i.insn.Src
}

func (i *Jump_) GetDst() uint8 {
	return i.insn.Dst
}

func (i *Jump_) GetOff() int16 {
	return i.insn.Off
}

func (i *Jump_) GetImm() int64 {
	return i.insn.Imm
}

func (i *Jump_) SetDst(dst uint8) Instruction {
	i.insn.Dst = dst
	return i
}

func (i *Jump_) SetSrc(src uint8) Instruction {
	i.insn.Src = src
	return i
}

func (i *Jump_) SetOff(off int16) Instruction {
	i.insn.Off = off
	return i
}

func (i *Jump_) SetImm(imm int64) Instruction {
	i.insn.Imm = imm
	return i
}

// //////////////////////////////////////////////////////////////

// / struct representation of CALL instruction
type FunctionCall struct {
	bpfCode *BpfCode
	insn    Insn
}

func (e *FunctionCall) OptCodeByte() uint8 {
	return BPF_CALL | BPF_JMP
}

func (m *FunctionCall) Push() *BpfCode {
	asm := m.Bytes()
	m.bpfCode.instructions = append(m.bpfCode.instructions, asm...)
	return m.bpfCode
}

func (m *FunctionCall) Bytes() []byte {
	i := Instruction(m)
	return IntoBytes(i)
}

func (m *FunctionCall) GetInsn() *Insn {
	return &m.insn
}

func (i *FunctionCall) GetSrc() uint8 {
	return i.insn.Src
}

func (i *FunctionCall) GetDst() uint8 {
	return i.insn.Dst
}

func (i *FunctionCall) GetOff() int16 {
	return i.insn.Off
}

func (i *FunctionCall) GetImm() int64 {
	return i.insn.Imm
}

func (i *FunctionCall) SetDst(dst uint8) Instruction {
	i.insn.Dst = dst
	return i
}

func (i *FunctionCall) SetSrc(src uint8) Instruction {
	i.insn.Src = src
	return i
}

func (i *FunctionCall) SetOff(off int16) Instruction {
	i.insn.Off = off
	return i
}

func (i *FunctionCall) SetImm(imm int64) Instruction {
	i.insn.Imm = imm
	return i
}

// //////////////////////////////////////////////////////////////

// / struct representation of EXIT instruction
type Exit struct {
	bpfCode *BpfCode
	insn    Insn
}

func (e *Exit) OptCodeByte() uint8 {
	return BPF_EXIT | BPF_JMP
}

func (m *Exit) Push() *BpfCode {
	asm := m.Bytes()
	m.bpfCode.instructions = append(m.bpfCode.instructions, asm...)
	return m.bpfCode
}

func (m *Exit) Bytes() []byte {
	i := Instruction(m)
	return IntoBytes(i)
}

func (m *Exit) GetInsn() *Insn {
	return &m.insn
}

func (i *Exit) GetSrc() uint8 {
	return i.insn.Src
}

func (i *Exit) GetDst() uint8 {
	return i.insn.Dst
}

func (i *Exit) GetOff() int16 {
	return i.insn.Off
}

func (i *Exit) GetImm() int64 {
	return i.insn.Imm
}

func (i *Exit) SetDst(dst uint8) Instruction {
	i.insn.Dst = dst
	return i
}

func (i *Exit) SetSrc(src uint8) Instruction {
	i.insn.Src = src
	return i
}

func (i *Exit) SetOff(off int16) Instruction {
	i.insn.Off = off
	return i
}

func (i *Exit) SetImm(imm int64) Instruction {
	i.insn.Imm = imm
	return i
}
