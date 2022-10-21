package insbuilder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIns(t *testing.T) {
	{
		/// let prog: &[u8] = &[
		///     0xb7, 0x12, 0x56, 0x34, 0xde, 0xbc, 0x9a, 0x78,
		///     ];
		/// let insn = ebpf::Insn {
		///     ptr: 0x00,
		///     opc: 0xb7,
		///     dst: 2,
		///     src: 1,
		///     off: 0x3456,
		///     imm: 0x789abcde
		/// };
		/// assert_eq!(insn.to_array(), prog);

		prog := [INSN_SIZE]uint8{0xb7, 0x12, 0x56, 0x34, 0xde, 0xbc, 0x9a, 0x78}

		insn := Insn{
			Ptr: 0x00,
			Opc: 0xb7,
			Dst: 2,
			Src: 1,
			Off: 0x3456,
			Imm: 0x789abcde,
		}

		require.Equal(t, prog, insn.ToArray())
	}
}

func TestCallImmediate(t *testing.T) {
	program := NewBpfCode()
	program.Call().SetImm(0x11_22_33_44).Push()

	require.Equal(t, program.Bytes(), []byte{0x85, 0x00, 0x00, 0x00, 0x44, 0x33, 0x22, 0x11})
}

func TestExitOperation(t *testing.T) {
	program := NewBpfCode()
	program.Exit().Push()

	require.Equal(t, []byte{0x95, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstEqualsSrc(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondEquals, SourceReg).SetDst(0x01).SetSrc(0x02).Push()

	require.Equal(t, []byte{0x1d, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstGreaterThanSrc(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondGreater, SourceReg).
		SetDst(0x03).
		SetSrc(0x02).
		Push()

	require.Equal(t, program.Bytes(), []byte{0x2d, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestJumpOnDstGreaterOrEqualsToSrc(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondGreaterEquals, SourceReg).SetDst(0x04).SetSrc(0x01).Push()

	require.Equal(t, program.Bytes(), []byte{0x3d, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestJumpOnDstLowerThanSrc(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondLower, SourceReg).SetDst(0x03).SetSrc(0x02).Push()
	require.Equal(t, []byte{0xad, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstLowerOrEqualsToSrc(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondLowerEquals, SourceReg).
		SetDst(0x04).
		SetSrc(0x01).
		Push()

	require.Equal(t, []byte{
		0xbd, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, program.Bytes())
}

func TestJumpOnDstBitAndWithSrcNotEqualZero(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondBitAnd, SourceReg).SetDst(0x05).SetSrc(0x02).Push()

	require.Equal(t, []byte{
		0x4d, 0x25, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, program.Bytes())
}

func TestJumpOnDstNotEqualsSrc(t *testing.T) {
	program := NewBpfCode()
	program.
		JumpConditional(CondNotEquals, SourceReg).
		SetDst(0x03).
		SetSrc(0x05).
		Push()

	require.Equal(t, []byte{0x5d, 0x53, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstGreaterThanSrcSigned(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondGreaterSigned, SourceReg).
		SetDst(0x04).
		SetSrc(0x01).
		Push()

	require.Equal(t, program.Bytes(), []byte{0x6d, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestJumpOnDstGreaterOrEqualsSrcSigned(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondGreaterEqualsSigned, SourceReg).SetDst(0x01).SetSrc(0x03).Push()

	require.Equal(t, []byte{0x7d, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstLowerThanSrcSigned(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondLowerSigned, SourceReg).
		SetDst(0x04).
		SetSrc(0x01).
		Push()

	require.Equal(t, []byte{0xcd, 0x14, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstLowerOrEqualsSrcSigned(t *testing.T) {
	program := NewBpfCode()
	program.
		JumpConditional(CondLowerEqualsSigned, SourceReg).
		SetDst(0x01).
		SetSrc(0x03).
		Push()

	require.Equal(t, []byte{
		0xdd, 0x31, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, program.Bytes())
}

func TestJumpToLabel(t *testing.T) {
	program := NewBpfCode()
	program.JumpUnconditional().SetOff(0x00_11).Push()

	require.Equal(t, program.Bytes(), []byte{0x05, 0x00, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestJumpOnDstEqualsConst(t *testing.T) {
	program := NewBpfCode()
	program.
		JumpConditional(CondEquals, SourceImm).
		SetDst(0x01).
		SetImm(0x00_11_22_33).
		Push()

	require.Equal(t, program.Bytes(), []byte{0x15, 0x01, 0x00, 0x00, 0x33, 0x22, 0x11, 0x00})
}

func TestJumpOnDstGreaterThanConst(t *testing.T) {
	program := NewBpfCode()
	program.
		JumpConditional(CondGreater, SourceImm).
		SetDst(0x02).
		SetImm(0x00110011).
		Push()

	require.Equal(t, []byte{
		0x25, 0x02, 0x00, 0x00, 0x11, 0x00, 0x11, 0x00,
	}, program.Bytes())
}

func TestJumpOnDstGreaterOrEqualsToConst(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondGreaterEquals, SourceImm).
		SetDst(0x04).
		SetImm(0x00_22_11_00).
		Push()

	require.Equal(
		t,
		[]byte{0x35, 0x04, 0x00, 0x00, 0x00, 0x11, 0x22, 0x00},
		program.Bytes(),
	)
}

func TestJumpOnDstLowerThanConst(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondLower, SourceImm).SetDst(0x02).SetImm(0x00_11_00_11).Push()

	require.Equal(t, []byte{0xa5, 0x02, 0x00, 0x00, 0x11, 0x00, 0x11, 0x00}, program.Bytes())
}

func TestJumpOnDstLowerOrEqualsToConst(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondLowerEquals, SourceImm).SetDst(0x04).SetImm(0x00_22_11_00).Push()

	require.Equal(t, []byte{0xb5, 0x04, 0x00, 0x00, 0x00, 0x11, 0x22, 0x00}, program.Bytes())
}

func Test_jump_on_dst_bit_and_with_const_not_equal_zero(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondBitAnd, SourceImm).
		SetDst(0x05).
		Push()

	require.Equal(t, []byte{0x45, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstNotEqualsConst(t *testing.T) {
	program := NewBpfCode()
	program.
		JumpConditional(CondNotEquals, SourceImm).
		SetDst(0x03).
		Push()

	require.Equal(t, []byte{0x55, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstGreaterThanConstSigned(t *testing.T) {
	var program BpfCode
	program.JumpConditional(CondGreaterSigned, SourceImm).SetDst(0x04).Push()

	require.Equal(
		t,
		[]byte{0x65, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
		program.Bytes(),
	)
}

func TestJumpOnDstGreaterOrEqualsSrcSigned2(t *testing.T) {
	program := NewBpfCode()
	program.
		JumpConditional(CondGreaterEqualsSigned, SourceImm).
		SetDst(0x01).
		Push()

	require.Equal(t, []byte{0x75, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestJumpOnDstLowerThanConstSigned(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondLowerSigned, SourceImm).SetDst(0x04).Push()

	require.Equal(t, []byte{
		0xc5, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, program.Bytes())
}

func TestJumpOnDstLowerOrEqualsSrcSigned2(t *testing.T) {
	program := NewBpfCode()
	program.JumpConditional(CondLowerEqualsSigned, SourceImm).SetDst(0x01).Push()

	require.Equal(t, []byte{0xd5, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestStoreWordFromDstIntoImmediateAddress(t *testing.T) {
	program := NewBpfCode()
	program.Store(MemSizeWord).SetDst(0x01).SetOff(0x0011).SetImm(0x11223344).Push()

	expected := []byte{0x62, 0x01, 0x11, 0x00, 0x44, 0x33, 0x22, 0x11}
	require.Equal(t, expected, program.Bytes())
}

func TestStoreHalfWordFromDstIntoImmediateAddress(t *testing.T) {
	program := NewBpfCode()
	program.
		Store(MemSizeHalfWord).SetDst(0x02).SetOff(0x11_22).Push()

	require.Equal(t, []byte{0x6a, 0x02, 0x22, 0x11, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestStoreByteFromDstIntoImmediateAddress(t *testing.T) {
	program := BpfCode{}
	program.Store(MemSizeByte).Push()

	require.Equal(t, program.Bytes(), []byte{0x72, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestStoreDoubleWordFromDstIntoImmediateAddress(t *testing.T) {
	program := NewBpfCode()
	program.Store(MemSizeDoubleWord).Push()

	require.Equal(t, []byte{
		0x7a, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, program.Bytes())
}

func TestStoreWordFromDstIntoSrcAddress(t *testing.T) {
	program := NewBpfCode()
	program.
		StoreX(MemSizeWord).
		SetDst(0x01).
		SetSrc(0x02).
		Push()

	require.Equal(t, []byte{0x63, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestStoreHalfWordFromDstIntoSrcAddress(t *testing.T) {
	program := NewBpfCode()
	program.StoreX(MemSizeHalfWord).Push()

	require.Equal(t, []byte{
		0x6b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, program.Bytes())
}

func TestStoreByteFromDstIntoSrcAddress(t *testing.T) {
	program := NewBpfCode()
	program.StoreX(MemSizeByte).Push()

	require.Equal(t, []byte{0x73, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestStoreDoubleWordFromDstIntoSrcAddress(t *testing.T) { // line 1103
	program := NewBpfCode()
	program.StoreX(MemSizeDoubleWord).Push()

	require.Equal(t, program.Bytes(), []byte{0x7b, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestLoadWordFromSetSrcWithOffset(t *testing.T) {
	program := NewBpfCode()
	program.
		LoadX(MemSizeWord).
		SetDst(0x01).
		SetSrc(0x02).
		SetOff(0x00_02).
		Push()

	require.Equal(t, []byte{0x61, 0x21, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestLoadHalfWordFromSetSrcWithOffset(t *testing.T) {
	program := NewBpfCode()
	program.LoadX(MemSizeHalfWord).SetDst(0x02).SetSrc(0x01).SetOff(0x1122).Push()

	require.Equal(t, []byte{0x69, 0x12, 0x22, 0x11, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestLoadByteFromSetSrcWithOffset(t *testing.T) {
	program := NewBpfCode()
	program.LoadX(MemSizeByte).SetDst(0x01).SetSrc(0x04).SetOff(0x00_11).Push()

	require.Equal(t, []byte{0x71, 0x41, 0x11, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestLoadDoubleWordFromSetSrcWithOffset(t *testing.T) {
	program := NewBpfCode()
	program.LoadX(MemSizeDoubleWord).SetDst(0x04).SetSrc(0x05).SetOff(0x4455).Push()
	require.Equal(t, []byte{0x79, 0x54, 0x55, 0x44, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestLoadDoubleWord(t *testing.T) {
	program := NewBpfCode()
	program.
		Load(MemSizeDoubleWord).
		SetDst(0x01).
		SetImm(0x00_01_02_03).
		Push()

	require.Equal(t, program.Bytes(), []byte{0x18, 0x01, 0x00, 0x00, 0x03, 0x02, 0x01, 0x00})
}

func TestLoadAbsWord(t *testing.T) {
	program := NewBpfCode()
	program.LoadAbs(MemSizeWord).Push()

	require.Equal(t, []byte{0x20, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestLoadAbsHalfWord(t *testing.T) {
	program := NewBpfCode()
	program.LoadAbs(MemSizeHalfWord).SetDst(0x05).Push()

	require.Equal(t, []byte{0x28, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestLoadAbsByte(t *testing.T) {
	program := NewBpfCode()
	program.LoadAbs(MemSizeByte).SetDst(0x01).Push()

	require.Equal(t, program.Bytes(), []byte{0x30, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00})
}

func TestLoadAbsDoubleWord(t *testing.T) {
	program := NewBpfCode()
	program.LoadAbs(MemSizeDoubleWord).SetDst(0x01).SetImm(0x01_02_03_04).Push()

	require.Equal(t, []byte{0x38, 0x01, 0x00, 0x00, 0x04, 0x03, 0x02, 0x01}, program.Bytes())
}

func TestLoadIndirectWord(t *testing.T) {
	program := NewBpfCode()
	program.LoadInd(MemSizeWord).Push()

	require.Equal(t, []byte{0x40, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestLoadIndirectHalfWord(t *testing.T) {
	program := NewBpfCode()
	program.LoadInd(MemSizeHalfWord).Push()

	require.Equal(t, []byte{
		0x48, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, program.Bytes())
}

func TestLoadIndirectByte(t *testing.T) {
	program := NewBpfCode()
	program.LoadInd(MemSizeByte).Push()

	require.Equal(t, []byte{0x50, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestLoadIndirectDoubleWord(t *testing.T) { // line 1286
	program := NewBpfCode()
	program.LoadInd(MemSizeDoubleWord).Push()

	require.Equal(t, []byte{0x58, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestConvertHostToLittleEndian16bits(t *testing.T) {
	program := NewBpfCode()
	program.SwapBytes(EndianLittle).SetDst(0x01).SetImm(0x00_00_00_10).Push()

	require.Equal(t, []byte{0xd4, 0x01, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestConvertHostToLittleEndian32bits(t *testing.T) {
	program := NewBpfCode()
	program.SwapBytes(EndianLittle).SetDst(0x02).SetImm(0x00_00_00_20).Push()

	require.Equal(t, []byte{0xd4, 0x02, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestConvertHostToLittleEndian64bit(t *testing.T) {
	program := NewBpfCode()
	program.SwapBytes(EndianLittle).SetDst(0x03).SetImm(0x00_00_00_40).Push()

	require.Equal(t, []byte{0xd4, 0x03, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestConvertHostToBigEndian16bits(t *testing.T) {
	program := NewBpfCode()
	program.SwapBytes(EndianBig).SetDst(0x01).SetImm(0x00_00_00_10).Push()

	require.Equal(t, []byte{0xdc, 0x01, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestConvertHostToBigEndian32bits(t *testing.T) {
	program := NewBpfCode()
	program.SwapBytes(EndianBig).SetDst(0x02).SetImm(0x00_00_00_20).Push()

	require.Equal(t, []byte{0xdc, 0x02, 0x00, 0x00, 0x20, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestConvertHostToBigEndian64bit(t *testing.T) {
	program := NewBpfCode()
	program.SwapBytes(EndianBig).SetDst(0x03).SetImm(0x00_00_00_40).Push()

	require.Equal(t, []byte{0xdc, 0x03, 0x00, 0x00, 0x40, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveAndAddConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Add(SourceImm, ArchX64).SetDst(0x02).SetImm(0x01_02_03_04).Push()

	require.Equal(t, []byte{0x07, 0x02, 0x00, 0x00, 0x04, 0x03, 0x02, 0x01}, program.Bytes())
}

func TestMoveSubConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Sub(SourceImm, ArchX64).SetDst(0x04).SetImm(0x00_01_02_03).Push()

	require.Equal(t, []byte{0x17, 0x04, 0x00, 0x00, 0x03, 0x02, 0x01, 0x00}, program.Bytes())
}

func TestMoveMulConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Mul(SourceImm, ArchX64).SetDst(0x05).SetImm(0x04_03_02_01).Push()

	require.Equal(t, []byte{0x27, 0x05, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04}, program.Bytes())
}

func TestMoveDivConstantToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Div(SourceImm, ArchX64).SetDst(0x02).SetImm(0x00_ff_00_ff).Push()

	require.Equal(t, []byte{0x37, 0x02, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00}, program.Bytes())
}

func TestMoveBitOrConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.BitOr(SourceImm, ArchX64).SetDst(0x02).SetImm(0x00_11_00_22).Push()

	require.Equal(t, []byte{0x47, 0x02, 0x00, 0x00, 0x22, 0x00, 0x11, 0x00}, program.Bytes())
}

func TestMoveBitAndConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.BitAnd(SourceImm, ArchX64).SetDst(0x02).SetImm(0x11_22_33_44).Push()

	require.Equal(t, []byte{0x57, 0x02, 0x00, 0x00, 0x44, 0x33, 0x22, 0x11}, program.Bytes())
}

func TestMoveLeftShiftConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.LeftShift(SourceImm, ArchX64).SetDst(0x01).Push()

	require.Equal(t, []byte{0x67, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveLogicalRightShiftConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.RightShift(SourceImm, ArchX64).SetDst(0x01).Push()

	require.Equal(t, []byte{0x77, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveNegateRegister(t *testing.T) {
	program := NewBpfCode()
	program.Negate(ArchX64).SetDst(0x02).Push()

	require.Equal(t, []byte{0x87, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveModConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Modulo(SourceImm, ArchX64).SetDst(0x02).Push()

	require.Equal(t, []byte{0x97, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveBitXorConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.BitXor(SourceImm, ArchX64).SetDst(0x03).Push()

	require.Equal(t, []byte{0xa7, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Mov(SourceImm, ArchX64).SetDst(0x01).SetImm(0x00_00_00_FF).Push()

	require.Equal(t, []byte{0xb7, 0x01, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveSignedRightShiftConstToRegister(t *testing.T) {
	program := NewBpfCode()
	program.SignedRightShift(SourceImm, ArchX64).SetDst(0x05).Push()

	require.Equal(t, []byte{0xc7, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveAndAddFromRegister(t *testing.T) {
	program := NewBpfCode()
	program.Add(SourceReg, ArchX64).SetDst(0x03).SetSrc(0x02).Push()

	require.Equal(t, []byte{0x0f, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveSubFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Sub(SourceReg, ArchX64).SetDst(0x03).SetSrc(0x04).Push()

	require.Equal(t, []byte{0x1f, 0x43, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveMulFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Mul(SourceReg, ArchX64).SetDst(0x04).SetSrc(0x03).Push()

	require.Equal(t, []byte{0x2f, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveDivFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Div(SourceReg, ArchX64).SetDst(0x01).SetSrc(0x00).Push()

	require.Equal(t, []byte{0x3f, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveBitOrFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.BitOr(SourceReg, ArchX64).SetDst(0x03).SetSrc(0x01).Push()

	require.Equal(t, []byte{0x4f, 0x13, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveBitAndFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.BitAnd(SourceReg, ArchX64).SetDst(0x03).SetSrc(0x02).Push()

	require.Equal(t, []byte{0x5f, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveLeftShiftFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.LeftShift(SourceReg, ArchX64).SetDst(0x02).SetSrc(0x03).Push()

	require.Equal(t, []byte{0x6f, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveLogicalRightShiftFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.RightShift(SourceReg, ArchX64).SetDst(0x02).SetSrc(0x04).Push()

	require.Equal(t, []byte{0x7f, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveModFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.Modulo(SourceReg, ArchX64).SetDst(0x01).SetSrc(0x02).Push()

	require.Equal(t, []byte{0x9f, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveBitXorFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.BitXor(SourceReg, ArchX64).SetDst(0x02).SetSrc(0x04).Push()

	require.Equal(t, []byte{0xaf, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveFromRegisterToAnotherRegister(t *testing.T) {
	program := NewBpfCode()
	program.Mov(SourceReg, ArchX64).SetSrc(0x01).Push()

	require.Equal(t, []byte{0xbf, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveSignedRightShiftFromRegisterToRegister(t *testing.T) {
	program := NewBpfCode()
	program.SignedRightShift(SourceReg, ArchX64).SetDst(0x02).SetSrc(0x03).Push()

	require.Equal(t, []byte{0xcf, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveAndAddConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Add(SourceImm, ArchX32).SetDst(0x02).SetImm(0x01020304).Push()

	require.Equal(t, []byte{0x04, 0x02, 0x00, 0x00, 0x04, 0x03, 0x02, 0x01}, program.Bytes())
}

func TestMoveSubConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Sub(SourceImm, ArchX32).SetDst(0x04).SetImm(0x00010203).Push()

	require.Equal(t, []byte{0x14, 0x04, 0x00, 0x00, 0x03, 0x02, 0x01, 0x00}, program.Bytes())
}

func TestMoveMulConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Mul(SourceImm, ArchX32).SetDst(0x05).SetImm(0x04030201).Push()

	require.Equal(t, []byte{0x24, 0x05, 0x00, 0x00, 0x01, 0x02, 0x03, 0x04}, program.Bytes())
}

func TestMoveDivConstantToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Div(SourceImm, ArchX32).SetDst(0x02).SetImm(0x00ff00ff).Push()

	require.Equal(t, []byte{0x34, 0x02, 0x00, 0x00, 0xff, 0x00, 0xff, 0x00}, program.Bytes())
}

func TestMoveBitOrConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.BitOr(SourceImm, ArchX32).SetDst(0x02).SetImm(0x00110022).Push()

	require.Equal(t, []byte{0x44, 0x02, 0x00, 0x00, 0x22, 0x00, 0x11, 0x00}, program.Bytes())
}

func TestMoveBitAndConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.BitAnd(SourceImm, ArchX32).SetDst(0x02).SetImm(0x11223344).Push()

	require.Equal(t, []byte{0x54, 0x02, 0x00, 0x00, 0x44, 0x33, 0x22, 0x11}, program.Bytes())
}

func TestMoveLeftShiftConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.LeftShift(SourceImm, ArchX32).SetDst(0x01).Push()

	require.Equal(t, []byte{0x64, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveLogicalRightShiftConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.RightShift(SourceImm, ArchX32).SetDst(0x01).Push()

	require.Equal(t, []byte{0x74, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveNegateRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Negate(ArchX32).SetDst(0x02).Push()

	require.Equal(t, []byte{0x84, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveModConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Modulo(SourceImm, ArchX32).SetDst(0x02).Push()

	require.Equal(t, []byte{0x94, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveBitXorConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.BitXor(SourceImm, ArchX32).SetDst(0x03).Push()

	require.Equal(t, []byte{0xa4, 0x03, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Mov(SourceImm, ArchX32).SetDst(0x01).SetImm(0x00_00_00_FF).Push()

	require.Equal(t, []byte{0xb4, 0x01, 0x00, 0x00, 0xff, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveSignedRightShiftConstToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.SignedRightShift(SourceImm, ArchX32).SetDst(0x05).Push()

	require.Equal(t, []byte{0xc4, 0x05, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveAndAddFromRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Add(SourceReg, ArchX32).SetDst(0x03).SetSrc(0x02).Push()

	require.Equal(t, []byte{0x0c, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveSubFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Sub(SourceReg, ArchX32).SetDst(0x03).SetSrc(0x04).Push()

	require.Equal(t, []byte{0x1c, 0x43, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveMulFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Mul(SourceReg, ArchX32).SetDst(0x04).SetSrc(0x03).Push()

	require.Equal(t, []byte{0x2c, 0x34, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveDivFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Div(SourceReg, ArchX32).SetDst(0x01).SetSrc(0x00).Push()

	require.Equal(t, []byte{0x3c, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveBitOrFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.BitOr(SourceReg, ArchX32).SetDst(0x03).SetSrc(0x01).Push()

	require.Equal(t, []byte{0x4c, 0x13, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveBitAndFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.BitAnd(SourceReg, ArchX32).SetDst(0x03).SetSrc(0x02).Push()

	require.Equal(t, []byte{0x5c, 0x23, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveLeftShiftFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.LeftShift(SourceReg, ArchX32).SetDst(0x02).SetSrc(0x03).Push()

	require.Equal(t, []byte{0x6c, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveLogicalRightShiftFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.RightShift(SourceReg, ArchX32).SetDst(0x02).SetSrc(0x04).Push()

	require.Equal(t, []byte{0x7c, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveModFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Modulo(SourceReg, ArchX32).SetDst(0x01).SetSrc(0x02).Push()

	require.Equal(t, []byte{0x9c, 0x21, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveBitXorFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.BitXor(SourceReg, ArchX32).SetDst(0x02).SetSrc(0x04).Push()

	require.Equal(t, []byte{0xac, 0x42, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveFromRegisterToAnotherRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.Mov(SourceReg, ArchX32).SetDst(0x00).SetSrc(0x01).Push()

	require.Equal(t, []byte{0xbc, 0x10, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestMoveSignedRightShiftFromRegisterToRegister_x32(t *testing.T) {
	program := NewBpfCode()
	program.SignedRightShift(SourceReg, ArchX32).SetDst(0x02).SetSrc(0x03).Push()

	require.Equal(t, []byte{0xcc, 0x32, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00}, program.Bytes())
}

func TestExampleFromAssembler_x64(t *testing.T) {
	program := NewBpfCode().
		Add(SourceImm, ArchX64).
		SetDst(1).
		SetImm(0x605).
		Push().
		Mov(SourceImm, ArchX64).
		SetDst(2).
		SetImm(0x32).
		Push().
		Mov(SourceReg, ArchX64).
		SetSrc(0).
		SetDst(1).
		Push().
		SwapBytes(EndianBig).
		SetDst(0).
		SetImm(0x10).
		Push().
		Negate(ArchX64).
		SetDst(2).
		Push().
		Exit().
		Push()

	require.Equal(t, []byte{
		0x07, 0x01, 0x00, 0x00, 0x05, 0x06, 0x00, 0x00,
		0xb7, 0x02, 0x00, 0x00, 0x32, 0x00, 0x00, 0x00,
		0xbf, 0x01, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0xdc, 0x00, 0x00, 0x00, 0x10, 0x00, 0x00, 0x00,
		0x87, 0x02, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x95, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	}, program.Bytes())
}
