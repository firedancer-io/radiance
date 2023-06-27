package safemath

import (
	"math"
	"math/rand"
	"testing"
)

func TestCheckedAddU8_Overflow(t *testing.T) {
	var a uint8 = 0xff
	var b uint8 = 2

	_, err := CheckedAddU8(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d + %d", a, b)
	}
}

func TestCheckedAddU8_NoOverflow(t *testing.T) {
	var a uint8 = 0xf0
	var b uint8 = 2

	result, err := CheckedAddU8(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d + %d", a, b)
	}
	if result != (a + b) {
		t.Errorf("wrong result for %d + %d", a, b)
	}
}

func TestCheckedAddU16_Overflow(t *testing.T) {
	var a uint16 = math.MaxUint16
	var b uint16 = 2

	_, err := CheckedAddU16(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d + %d", a, b)
	}
}

func TestCheckedAddU16_NoOverflow(t *testing.T) {
	var a uint16 = 100
	var b uint16 = 2

	result, err := CheckedAddU16(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d + %d", a, b)
	}
	if result != (a + b) {
		t.Errorf("wrong result for %d + %d", a, b)
	}
}

func TestCheckedAddU32_Overflow(t *testing.T) {
	var a uint32 = math.MaxUint32
	var b uint32 = 2

	_, err := CheckedAddU32(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d + %d", a, b)
	}
}

func TestCheckedAddU32_NoOverflow(t *testing.T) {
	var a uint32 = 100
	var b uint32 = 2

	result, err := CheckedAddU32(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d + %d", a, b)
	}
	if result != (a + b) {
		t.Errorf("wrong result for %d + %d", a, b)
	}
}

func TestCheckedAddU64_Overflow(t *testing.T) {
	var a uint64 = math.MaxUint64
	var b uint64 = 2

	_, err := CheckedAddU64(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d + %d", a, b)
	}
}

func TestCheckedAddU64_NoOverflow(t *testing.T) {
	var a uint64 = 100
	var b uint64 = 2

	result, err := CheckedAddU64(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d + %d", a, b)
	}
	if result != (a + b) {
		t.Errorf("wrong result for %d + %d", a, b)
	}
}

func TestCheckedMulU8_Overflow(t *testing.T) {
	var a uint8 = 0xff / 2
	var b uint8 = 3

	_, err := CheckedMulU8(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d * %d", a, b)
	}
}

func TestCheckedMulU8_NoOverflow(t *testing.T) {
	var a uint8 = 22
	var b uint8 = 2

	result, err := CheckedMulU8(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d * %d", a, b)
	}
	if result != (a * b) {
		t.Errorf("wrong result for %d * %d", a, b)
	}
}

func TestCheckedMulU16_Overflow(t *testing.T) {
	var a uint16 = math.MaxUint16 / 2
	var b uint16 = 3

	_, err := CheckedMulU16(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d * %d", a, b)
	}
}

func TestCheckedMulU16_NoOverflow(t *testing.T) {
	var a uint16 = 22
	var b uint16 = 2

	result, err := CheckedMulU16(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d * %d", a, b)
	}
	if result != (a * b) {
		t.Errorf("wrong result for %d * %d", a, b)
	}
}

func TestCheckedMulU32_Overflow(t *testing.T) {
	var a uint32 = math.MaxUint32 / 2
	var b uint32 = 3

	_, err := CheckedMulU32(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d * %d", a, b)
	}
}

func TestCheckedMulU32_NoOverflow(t *testing.T) {
	var a uint32 = 22
	var b uint32 = 2

	result, err := CheckedMulU32(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d * %d", a, b)
	}
	if result != (a * b) {
		t.Errorf("wrong result for %d * %d", a, b)
	}
}

func TestCheckedMulU64_Overflow(t *testing.T) {
	var a uint64 = math.MaxUint64 / 2
	var b uint64 = 3

	_, err := CheckedMulU64(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d * %d", a, b)
	}
}

func TestCheckedMulU64_NoOverflow(t *testing.T) {
	var a uint64 = 22
	var b uint64 = 2

	result, err := CheckedMulU64(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d * %d", a, b)
	}
	if result != (a * b) {
		t.Errorf("wrong result for %d * %d", a, b)
	}
}

func TestCheckedSubU8_Underflow(t *testing.T) {
	var a uint8 = 0x10
	var b uint8 = 0xff

	_, err := CheckedSubU8(a, b)
	if err == nil {
		t.Errorf("should've detected underflow in calculating %d - %d", a, b)
	}
}

func TestCheckedSubU8_NoUnderflow(t *testing.T) {
	var a uint8 = 22
	var b uint8 = 2

	result, err := CheckedSubU8(a, b)
	if err != nil {
		t.Errorf("should NOT have detected underflow in calculating %d - %d", a, b)
	}
	if result != (a - b) {
		t.Errorf("wrong result for %d - %d", a, b)
	}
}

func TestCheckedSubU16_Underflow(t *testing.T) {
	var a uint16 = 0x10
	var b uint16 = 0xff

	_, err := CheckedSubU16(a, b)
	if err == nil {
		t.Errorf("should've detected underflow in calculating %d - %d", a, b)
	}
}

func TestCheckedSubU16_NoUnderflow(t *testing.T) {
	var a uint16 = 22
	var b uint16 = 2

	result, err := CheckedSubU16(a, b)
	if err != nil {
		t.Errorf("should NOT have detected underflow in calculating %d - %d", a, b)
	}
	if result != (a - b) {
		t.Errorf("wrong result for %d - %d", a, b)
	}
}

func TestCheckedSubU32_Underflow(t *testing.T) {
	var a uint32 = 0x10
	var b uint32 = 0xff

	_, err := CheckedSubU32(a, b)
	if err == nil {
		t.Errorf("should've detected underflow in calculating %d - %d", a, b)
	}
}

func TestCheckedSubU32_NoUnderflow(t *testing.T) {
	var a uint32 = 20
	var b uint32 = 2

	result, err := CheckedSubU32(a, b)
	if err != nil {
		t.Errorf("should NOT have detected underflow in calculating %d - %d", a, b)
	}
	if result != (a - b) {
		t.Errorf("wrong result for %d - %d", a, b)
	}
}

func TestCheckedSubU64_Underflow(t *testing.T) {
	var a uint64 = 0x10
	var b uint64 = 0xff

	_, err := CheckedSubU64(a, b)
	if err == nil {
		t.Errorf("should've detected overflow in calculating %d - %d", a, b)
	}
}

func TestCheckedSubU64_NoUnderflow(t *testing.T) {
	var a uint64 = 22
	var b uint64 = 2

	result, err := CheckedSubU64(a, b)
	if err != nil {
		t.Errorf("should NOT have detected overflow in calculating %d * %d", a, b)
	}
	if result != (a - b) {
		t.Errorf("wrong result for %d - %d", a, b)
	}
}

func TestCheckedDivU16_DivByZero(t *testing.T) {
	var a uint16 = 0x10
	var b uint16 = 0

	_, err := CheckedDivU16(a, b)
	if err == nil {
		t.Errorf("should've detected division-by-zero in calculating %d / %d", a, b)
	}
}

func TestCheckedDivU16_NoDivByZero(t *testing.T) {
	var a uint16 = 22
	var b uint16 = 2

	result, err := CheckedDivU16(a, b)
	if err != nil {
		t.Errorf("should NOT have detected division-by-zero in calculating %d / %d", a, b)
	}
	if result != (a / b) {
		t.Errorf("wrong result for %d / %d", a, b)
	}
}

func TestCheckedDivU32_DivByZero(t *testing.T) {
	var a uint32 = 0x10
	var b uint32 = 0

	_, err := CheckedDivU32(a, b)
	if err == nil {
		t.Errorf("should've detected division-by-zero in calculating %d / %d", a, b)
	}
}

func TestCheckedDivU32_NoDivByZero(t *testing.T) {
	var a uint32 = 20
	var b uint32 = 2

	result, err := CheckedDivU32(a, b)
	if err != nil {
		t.Errorf("should NOT have detected division-by-zero in calculating %d / %d", a, b)
	}
	if result != (a / b) {
		t.Errorf("wrong result for %d / %d", a, b)
	}
}

func TestCheckedDivU64_DivByZero(t *testing.T) {
	var a uint64 = 0x10
	var b uint64 = 0

	_, err := CheckedDivU64(a, b)
	if err == nil {
		t.Errorf("should've detected division-by-zero in calculating %d / %d", a, b)
	}
}

func TestCheckedDivU64_NoDivByZero(t *testing.T) {
	var a uint64 = 22
	var b uint64 = 2

	result, err := CheckedDivU64(a, b)
	if err != nil {
		t.Errorf("should NOT have detected division-by-zero in calculating %d / %d", a, b)
	}
	if result != (a / b) {
		t.Errorf("wrong result for %d / %d", a, b)
	}
}

func TestSaturatingAddU8_ShouldSaturate(t *testing.T) {
	var a uint8 = 0xfe
	var b uint8 = 6

	result := SaturatingAddU8(a, b)
	if result != math.MaxUint8 {
		t.Errorf("result should've saturated in calculating %d + %d", a, b)
	}
}

func TestSaturatingAddU8_ShouldNotSaturate(t *testing.T) {
	var a uint8 = 10
	var b uint8 = 20

	result := SaturatingAddU8(a, b)
	if result != (a + b) {
		t.Errorf("wrong result in calculating %d + %d (got %d)", a, b, result)
	}
}

func TestSaturatingMulU8_ShouldSaturate(t *testing.T) {
	var a uint8 = 0xff / 2
	var b uint8 = 3

	result := SaturatingMulU8(a, b)
	if result != math.MaxUint8 {
		t.Errorf("result should've saturated in calculating %d * %d", a, b)
	}
}

func TestSaturatingMulU8_ShouldNotSaturate(t *testing.T) {
	var a uint8 = uint8(rand.Intn(5))
	var b uint8 = uint8(rand.Intn(5))

	result := SaturatingMulU8(a, b)
	if result != (a * b) {
		t.Errorf("wrong result in calculating %d * %d (got %d)", a, b, result)
	}
}

func TestSaturatingSubU8_ShouldSaturate(t *testing.T) {
	var a uint8 = 10
	var b uint8 = 0x20

	result := SaturatingSubU8(a, b)
	if result != 0 {
		t.Errorf("result should've saturated in calculating %d - %d", a, b)
	}
}

func TestSaturatingSubU8_ShouldNotSaturate(t *testing.T) {
	var a uint8 = uint8(rand.Intn(50) + 10)
	var b uint8 = uint8(rand.Intn(10))

	result := SaturatingSubU8(a, b)
	if result != (a - b) {
		t.Errorf("wrong result in calculating %d - %d (got %d)", a, b, result)
	}
}

func TestSaturatingAddU16_ShouldSaturate(t *testing.T) {
	var a uint16 = math.MaxUint16
	var b uint16 = 6

	result := SaturatingAddU16(a, b)
	if result != math.MaxUint16 {
		t.Errorf("result should've saturated in calculating %d + %d", a, b)
	}
}

func TestSaturatingAddU16_ShouldNotSaturate(t *testing.T) {
	var a uint16 = uint16(rand.Intn(200))
	var b uint16 = uint16(rand.Intn(200))

	result := SaturatingAddU16(a, b)
	if result != (a + b) {
		t.Errorf("wrong result in calculating %d + %d (got %d)", a, b, result)
	}
}

func TestSaturatingMulU16_ShouldSaturate(t *testing.T) {
	var a uint16 = math.MaxUint16 / 2
	var b uint16 = 3

	result := SaturatingMulU16(a, b)
	if result != math.MaxUint16 {
		t.Errorf("result should've saturated in calculating %d * %d", a, b)
	}
}

func TestSaturatingMulU16_ShouldNotSaturate(t *testing.T) {
	var a uint16 = uint16(rand.Intn(20))
	var b uint16 = uint16(rand.Intn(20))

	result := SaturatingMulU16(a, b)
	if result != (a * b) {
		t.Errorf("wrong result in calculating %d * %d (got %d)", a, b, result)
	}
}

func TestSaturatingSubU16_ShouldSaturate(t *testing.T) {
	var a uint16 = 10
	var b uint16 = 0x20

	result := SaturatingSubU16(a, b)
	if result != 0 {
		t.Errorf("result should've saturated in calculating %d - %d", a, b)
	}
}

func TestSaturatingSubU16_ShouldNotSaturate(t *testing.T) {
	var a uint16 = uint16(rand.Intn(100) + 50)
	var b uint16 = uint16(rand.Intn(50))

	result := SaturatingSubU16(a, b)
	if result != (a - b) {
		t.Errorf("wrong result in calculating %d - %d (got %d)", a, b, result)
	}
}

func TestSaturatingAddU32_ShouldSaturate(t *testing.T) {
	var a uint32 = math.MaxUint32
	var b uint32 = 6

	result := SaturatingAddU32(a, b)
	if result != math.MaxUint32 {
		t.Errorf("result should've saturated in calculating %d + %d", a, b)
	}
}

func TestSaturatingAddU32_ShouldNotSaturate(t *testing.T) {
	var a uint32 = uint32(rand.Intn(1000))
	var b uint32 = uint32(rand.Intn(1000))

	result := SaturatingAddU32(a, b)
	if result != (a + b) {
		t.Errorf("wrong result in calculating %d + %d (got %d)", a, b, result)
	}
}

func TestSaturatingMulU32_ShouldSaturate(t *testing.T) {
	var a uint32 = math.MaxUint32 / 2
	var b uint32 = 3

	result := SaturatingMulU32(a, b)
	if result != math.MaxUint32 {
		t.Errorf("result should've saturated in calculating %d * %d", a, b)
	}
}

func TestSaturatingMulU32_ShouldNotSaturate(t *testing.T) {
	var a uint32 = uint32(rand.Intn(200))
	var b uint32 = uint32(rand.Intn(1000))

	result := SaturatingMulU32(a, b)
	if result != (a * b) {
		t.Errorf("wrong result in calculating %d * %d (got %d)", a, b, result)
	}
}

func TestSaturatingSubU32_ShouldSaturate(t *testing.T) {
	var a uint32 = 10
	var b uint32 = 0x20

	result := SaturatingSubU32(a, b)
	if result != 0 {
		t.Errorf("result should've saturated in calculating %d - %d", a, b)
	}
}

func TestSaturatingSubU32_ShouldNotSaturate(t *testing.T) {
	var a uint32 = uint32(rand.Intn(1000) + 500)
	var b uint32 = uint32(rand.Intn(500))

	result := SaturatingSubU32(a, b)
	if result != (a - b) {
		t.Errorf("wrong result in calculating %d - %d (got %d)", a, b, result)
	}
}

func TestSaturatingAddU64_ShouldSaturate(t *testing.T) {
	var a uint64 = math.MaxUint64
	var b uint64 = 6

	result := SaturatingAddU64(a, b)
	if result != math.MaxUint64 {
		t.Errorf("result should've saturated in calculating %d + %d", a, b)
	}
}

func TestSaturatingAddU64_ShouldNotSaturate(t *testing.T) {
	var a uint64 = uint64(rand.Intn(20000))
	var b uint64 = uint64(rand.Intn(20000))

	result := SaturatingAddU64(a, b)
	if result != (a + b) {
		t.Errorf("wrong result in calculating %d + %d (got %d)", a, b, result)
	}
}

func TestSaturatingMulU64_ShouldSaturate(t *testing.T) {
	var a uint64 = math.MaxUint64 / 2
	var b uint64 = 3

	result := SaturatingMulU64(a, b)
	if result != math.MaxUint64 {
		t.Errorf("result should've saturated in calculating %d * %d", a, b)
	}
}

func TestSaturatingMulU64_ShouldNotSaturate(t *testing.T) {
	var a uint64 = uint64(rand.Intn(50000))
	var b uint64 = uint64(rand.Intn(1000))

	result := SaturatingMulU64(a, b)
	if result != (a * b) {
		t.Errorf("wrong result in calculating %d * %d (got %d)", a, b, result)
	}
}

func TestSaturatingSubU64_ShouldSaturate(t *testing.T) {
	var a uint64 = 10
	var b uint64 = 0x20

	result := SaturatingSubU64(a, b)
	if result != 0 {
		t.Errorf("result should've saturated in calculating %d - %d", a, b)
	}
}

func TestSaturatingSubU64_ShouldNotSaturate(t *testing.T) {
	var a uint64 = uint64(rand.Intn(50000) + 1000)
	var b uint64 = uint64(rand.Intn(1000))

	result := SaturatingSubU64(a, b)
	if result != (a - b) {
		t.Errorf("wrong result in calculating %d - %d (got %d)", a, b, result)
	}
}
