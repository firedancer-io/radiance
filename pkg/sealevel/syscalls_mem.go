package sealevel

import (
	"errors"

	"go.firedancer.io/radiance/pkg/sbpf"
	"go.firedancer.io/radiance/pkg/sbpf/cu"
)

var (
	ErrCopyOverlapping = errors.New("Overlapping copy")
)

func MemOpConsume(cuIn int, n uint64) int {
	perBytesCost := n / CuCpiBytesPerUnit
	return cu.ConsumeLowerBound(cuIn, CUMemOpBaseCost, int(perBytesCost))
}

func isNonOverlapping(src, dst, n uint64) bool {
	if src > dst {
		return src-dst >= n
	} else {
		return dst-src >= n
	}
}

func memmoveImplInternal(vm sbpf.VM, dst, src, n uint64) (err error) {
	srcBuf := make([]byte, n)
	err = vm.Read(src, srcBuf)
	if err != nil {
		return
	}
	err = vm.Write(dst, srcBuf)
	return
}

// SyscallMemcpyImpl is the implementation of the memcpy (sol_memcpy_) syscall.
// Overlapping src and dst for a given n bytes to be copied results in an error being returned.
func SyscallMemcpyImpl(vm sbpf.VM, dst, src, n uint64, cuIn int) (r0 uint64, cuOut int, err error) {
	cuOut = MemOpConsume(cuIn, n)
	if cuOut < 0 {
		return
	}

	// memcpy when src and dst are overlapping results in undefined behaviour,
	// hence check if there is an overlap and return early with an error if so.
	if !isNonOverlapping(src, dst, n) {
		return r0, cuOut, ErrCopyOverlapping
	}

	err = memmoveImplInternal(vm, dst, src, n)
	return
}

var SyscallMemcpy = sbpf.SyscallFunc3(SyscallMemcpyImpl)

// SyscallMemmoveImpl is the implementation for the memmove (sol_memmove_) syscall.
func SyscallMemmoveImpl(vm sbpf.VM, dst, src, n uint64, cuIn int) (r0 uint64, cuOut int, err error) {
	cuOut = MemOpConsume(cuIn, n)
	if cuOut < 0 {
		return
	}
	err = memmoveImplInternal(vm, dst, src, n)
	return
}

var SyscallMemmove = sbpf.SyscallFunc3(SyscallMemmoveImpl)
