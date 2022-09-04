package sbf

import "fmt"

type Verifier struct {
	Program *Program
}

func NewVerifier(p *Program) *Verifier {
	return &Verifier{Program: p}
}

func (v *Verifier) Verify() error {
	text := v.Program.Text
	if len(text)%SlotSize != 0 {
		return fmt.Errorf("odd .text size")
	}
	if len(text) == 0 {
		return fmt.Errorf("empty text")
	}

	for pc := uint64(0); (pc+1)*SlotSize <= uint64(len(text)); pc++ {
		insBytes := text[pc*SlotSize:]
		ins := GetSlot(insBytes)

		if ins.Src() > 10 {
			return fmt.Errorf("invalid src register")
		}
		switch ins.Op() {
		case OpLdxb, OpLdxh, OpLdxw, OpLdxdw:
		case OpAdd32Imm, OpAdd32Reg, OpAdd64Imm, OpAdd64Reg:
		case OpSub32Imm, OpSub32Reg, OpSub64Imm, OpSub64Reg:
		case OpMul32Imm, OpMul32Reg, OpMul64Imm, OpMul64Reg:
		case OpOr32Imm, OpOr32Reg, OpOr64Imm, OpOr64Reg:
		case OpAnd32Imm, OpAnd32Reg, OpAnd64Imm, OpAnd64Reg:
		case OpLsh32Reg, OpLsh64Reg:
		case OpRsh32Reg, OpRsh64Reg:
		case OpNeg32, OpNeg64:
		case OpXor32Imm, OpXor32Reg, OpXor64Imm, OpXor64Reg:
		case OpMov32Imm, OpMov32Reg, OpMov64Imm, OpMov64Reg:
		case OpDiv32Reg, OpDiv64Reg:
		case OpMod32Reg, OpMod64Reg:
		case OpSdiv32Reg, OpSdiv64Reg:
		case OpCall, OpExit:
			// nothing
		case OpStb, OpSth, OpStw, OpStdw,
			OpStxb, OpStxh, OpStxw, OpStxdw:
			if ins.Dst() > 10 {
				return fmt.Errorf("invalid dst register")
			}
			continue
		case OpLsh32Imm, OpRsh32Imm, OpArsh32Imm:
			if ins.Uimm() > 31 {
				return fmt.Errorf("32-bit shift out of bounds")
			}
		case OpLsh64Imm, OpRsh64Imm, OpArsh64Imm:
			if ins.Uimm() > 63 {
				return fmt.Errorf("64-bit shift out of bounds")
			}
		case OpLe, OpBe:
			switch ins.Uimm() {
			case 16, 32, 64:
				// ok
			default:
				return fmt.Errorf("invalid bit size for endianness conversion")
			}
		case OpSdiv32Imm, OpSdiv64Imm:
			fallthrough
		case OpDiv32Imm, OpDiv64Imm, OpMod32Imm, OpMod64Imm:
			if ins.Imm() == 0 {
				return ExcDivideByZero
			}
		case OpJa,
			OpJeqImm, OpJeqReg,
			OpJgtImm, OpJgtReg,
			OpJgeImm, OpJgeReg,
			OpJltImm, OpJltReg,
			OpJleImm, OpJleReg,
			OpJsetImm, OpJsetReg,
			OpJneImm, OpJneReg,
			OpJsgtImm, OpJsgtReg,
			OpJsgeImm, OpJsgeReg,
			OpJsltImm, OpJsltReg,
			OpJsleImm, OpJsleReg:
			dst := int64(pc) + int64(ins.Off()) + 1
			if dst < 0 || (dst*SlotSize) >= int64(len(text)) {
				return fmt.Errorf("jump out of code")
			}
			dstIns := GetSlot(text[dst*SlotSize:])
			if dstIns.Op() == 0 {
				return fmt.Errorf("jump into middle of instruction")
			}
		case OpCallx:
			if uimm := ins.Uimm(); uimm >= 10 {
				return fmt.Errorf("invalid callx register")
			}
		case OpLddw:
			if len(insBytes) < 2*SlotSize {
				return fmt.Errorf("incomplete lddw instruction")
			}
			if insBytes[8] != 0 {
				return fmt.Errorf("malformed lddw instruction")
			}
			pc++
		default:
			return fmt.Errorf("unknown opcode %#02x", ins.Op())
		}

		if ins.Dst() > 9 {
			return fmt.Errorf("invalid dst register")
		}
	}

	return nil
}
