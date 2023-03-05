package sbpf

// Op classes
const (
	ClassLd = uint8(iota)
	ClassLdx
	ClassSt
	ClassStx
	ClassAlu
	ClassJmp
	Class0x06 // reserved
	ClassAlu64
)

// Size modes
const (
	SizeW = uint8(iota * 0x08)
	SizeH
	SizeB
	SizeDw
)

// Addressing modes
const (
	AddrImm = uint8(iota * 0x20)
	AddrAbs
	AddrInd
	AddrMem
	Addr0x80 // reserved
	Addr0xa0 // reserved
	AddrXadd
)

// Source modes
const (
	SrcK = uint8(iota * 0x08)
	SrcX
)

// ALU operations
const (
	AluAdd = uint8(iota * 0x10)
	AluSub
	AluMul
	AluDiv
	AluOr
	AluAnd
	AluLsh
	AluRsh
	AluNeg
	AluMod
	AluXor
	AluMov
	AluArsh
	AluEnd
	AluSdiv
)

// Jump operations
const (
	JumpAlways = uint8(iota * 0x10)
	JumpEq
	JumpGt
	JumpGe
	JumpSet
	JumpNe
	JumpSgt
	JumpSge
	JumpCall
	JumpExit
	JumpLt
	JumpLe
	JumpSlt
	JumpSle
)

// Opcodes
const (
	OpLddw = ClassLd | AddrImm | SizeDw

	OpLdxb  = ClassLdx | AddrMem | SizeB
	OpLdxh  = ClassLdx | AddrMem | SizeH
	OpLdxw  = ClassLdx | AddrMem | SizeW
	OpLdxdw = ClassLdx | AddrMem | SizeDw
	OpStb   = ClassSt | AddrMem | SizeB
	OpSth   = ClassSt | AddrMem | SizeH
	OpStw   = ClassSt | AddrMem | SizeW
	OpStdw  = ClassSt | AddrMem | SizeDw
	OpStxb  = ClassStx | AddrMem | SizeB
	OpStxh  = ClassStx | AddrMem | SizeH
	OpStxw  = ClassStx | AddrMem | SizeW
	OpStxdw = ClassStx | AddrMem | SizeDw

	OpAdd32Imm  = ClassAlu | SrcK | AluAdd
	OpAdd32Reg  = ClassAlu | SrcX | AluAdd
	OpSub32Imm  = ClassAlu | SrcK | AluSub
	OpSub32Reg  = ClassAlu | SrcX | AluSub
	OpMul32Imm  = ClassAlu | SrcK | AluMul
	OpMul32Reg  = ClassAlu | SrcX | AluMul
	OpDiv32Imm  = ClassAlu | SrcK | AluDiv
	OpDiv32Reg  = ClassAlu | SrcX | AluDiv
	OpOr32Imm   = ClassAlu | SrcK | AluOr
	OpOr32Reg   = ClassAlu | SrcX | AluOr
	OpAnd32Imm  = ClassAlu | SrcK | AluAnd
	OpAnd32Reg  = ClassAlu | SrcX | AluAnd
	OpLsh32Imm  = ClassAlu | SrcK | AluLsh
	OpLsh32Reg  = ClassAlu | SrcX | AluLsh
	OpRsh32Imm  = ClassAlu | SrcK | AluRsh
	OpRsh32Reg  = ClassAlu | SrcX | AluRsh
	OpNeg32     = ClassAlu | AluNeg
	OpMod32Imm  = ClassAlu | SrcK | AluMod
	OpMod32Reg  = ClassAlu | SrcX | AluMod
	OpXor32Imm  = ClassAlu | SrcK | AluXor
	OpXor32Reg  = ClassAlu | SrcX | AluXor
	OpMov32Imm  = ClassAlu | SrcK | AluMov
	OpMov32Reg  = ClassAlu | SrcX | AluMov
	OpArsh32Imm = ClassAlu | SrcK | AluArsh
	OpArsh32Reg = ClassAlu | SrcX | AluArsh
	OpSdiv32Imm = ClassAlu | SrcK | AluSdiv
	OpSdiv32Reg = ClassAlu | SrcX | AluSdiv
	OpLe        = ClassAlu | SrcK | AluEnd
	OpBe        = ClassAlu | SrcX | AluEnd

	OpAdd64Imm  = ClassAlu64 | SrcK | AluAdd
	OpAdd64Reg  = ClassAlu64 | SrcX | AluAdd
	OpSub64Imm  = ClassAlu64 | SrcK | AluSub
	OpSub64Reg  = ClassAlu64 | SrcX | AluSub
	OpMul64Imm  = ClassAlu64 | SrcK | AluMul
	OpMul64Reg  = ClassAlu64 | SrcX | AluMul
	OpDiv64Imm  = ClassAlu64 | SrcK | AluDiv
	OpDiv64Reg  = ClassAlu64 | SrcX | AluDiv
	OpOr64Imm   = ClassAlu64 | SrcK | AluOr
	OpOr64Reg   = ClassAlu64 | SrcX | AluOr
	OpAnd64Imm  = ClassAlu64 | SrcK | AluAnd
	OpAnd64Reg  = ClassAlu64 | SrcX | AluAnd
	OpLsh64Imm  = ClassAlu64 | SrcK | AluLsh
	OpLsh64Reg  = ClassAlu64 | SrcX | AluLsh
	OpRsh64Imm  = ClassAlu64 | SrcK | AluRsh
	OpRsh64Reg  = ClassAlu64 | SrcX | AluRsh
	OpNeg64     = ClassAlu64 | AluNeg
	OpMod64Imm  = ClassAlu64 | SrcK | AluMod
	OpMod64Reg  = ClassAlu64 | SrcX | AluMod
	OpXor64Imm  = ClassAlu64 | SrcK | AluXor
	OpXor64Reg  = ClassAlu64 | SrcX | AluXor
	OpMov64Imm  = ClassAlu64 | SrcK | AluMov
	OpMov64Reg  = ClassAlu64 | SrcX | AluMov
	OpArsh64Imm = ClassAlu64 | SrcK | AluArsh
	OpArsh64Reg = ClassAlu64 | SrcX | AluArsh
	OpSdiv64Imm = ClassAlu64 | SrcK | AluSdiv
	OpSdiv64Reg = ClassAlu64 | SrcX | AluSdiv

	OpJa      = ClassJmp | JumpAlways
	OpJeqImm  = ClassJmp | SrcK | JumpEq
	OpJeqReg  = ClassJmp | SrcX | JumpEq
	OpJgtImm  = ClassJmp | SrcK | JumpGt
	OpJgtReg  = ClassJmp | SrcX | JumpGt
	OpJgeImm  = ClassJmp | SrcK | JumpGe
	OpJgeReg  = ClassJmp | SrcX | JumpGe
	OpJltImm  = ClassJmp | SrcK | JumpLt
	OpJltReg  = ClassJmp | SrcX | JumpLt
	OpJleImm  = ClassJmp | SrcK | JumpLe
	OpJleReg  = ClassJmp | SrcX | JumpLe
	OpJsetImm = ClassJmp | SrcK | JumpSet
	OpJsetReg = ClassJmp | SrcX | JumpSet
	OpJneImm  = ClassJmp | SrcK | JumpNe
	OpJneReg  = ClassJmp | SrcX | JumpNe
	OpJsgtImm = ClassJmp | SrcK | JumpSgt
	OpJsgtReg = ClassJmp | SrcX | JumpSgt
	OpJsgeImm = ClassJmp | SrcK | JumpSge
	OpJsgeReg = ClassJmp | SrcX | JumpSge
	OpJsltImm = ClassJmp | SrcK | JumpSlt
	OpJsltReg = ClassJmp | SrcX | JumpSlt
	OpJsleImm = ClassJmp | SrcK | JumpSle
	OpJsleReg = ClassJmp | SrcX | JumpSle

	OpCall  = ClassJmp | SrcK | JumpCall
	OpCallx = ClassJmp | SrcX | JumpCall
	OpExit  = ClassJmp | JumpExit
)
