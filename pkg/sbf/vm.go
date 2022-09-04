package sbf

import (
	"errors"
	"fmt"
)

// VM is the virtual machine abstraction, implemented by each executor.
type VM interface {
	VMContext() any

	Read(addr uint64, p []byte) error
	Read8(addr uint64) (uint8, error)
	Read16(addr uint64) (uint16, error)
	Read32(addr uint64) (uint32, error)
	Read64(addr uint64) (uint64, error)

	Write(addr uint64, p []byte) error
	Write8(addr uint64, x uint8) error
	Write16(addr uint64, x uint16) error
	Write32(addr uint64, x uint32) error
	Write64(addr uint64, x uint64) error
}

// VMOpts specifies virtual machine parameters.
type VMOpts struct {
	// Machine parameters
	StackSize int
	HeapSize  int
	Syscalls  SyscallRegistry

	// Execution parameters
	Context any // passed to syscalls
	MaxCU   uint64
	Input   []byte // mapped at VaddrInput
}

// Syscall are callback handles from VM to Go. (work in progress)
type Syscall interface {
	Invoke(vm VM, r1, r2, r3, r4, r5 uint64, cuIn int64) (r0 uint64, cuOut int64, err error)
}

// Exception codes.
var (
	ExcDivideByZero   = errors.New("division by zero")
	ExcDivideOverflow = errors.New("divide overflow")
)

type ExcBadAccess struct {
	Addr   uint64
	Size   uint32
	Write  bool
	Reason string
}

func NewExcBadAccess(addr uint64, size uint32, write bool, reason string) ExcBadAccess {
	return ExcBadAccess{
		Addr:   addr,
		Size:   size,
		Write:  write,
		Reason: reason,
	}
}

func (e ExcBadAccess) Error() string {
	return fmt.Sprintf("bad memory access at %#x (size=%d write=%v), reason: %s", e.Addr, e.Size, e.Write, e.Reason)
}

// Convenience Methods

type SyscallFunc0 func(vm VM, cuIn int64) (r0 uint64, cuOut int64, err error)

func (f SyscallFunc0) Invoke(vm VM, _, _, _, _, _ uint64, cuIn int64) (r0 uint64, cuOut int64, err error) {
	return f(vm, cuIn)
}

type SyscallFunc1 func(vm VM, r1 uint64, cuIn int64) (r0 uint64, cuOut int64, err error)

func (f SyscallFunc1) Invoke(vm VM, r1, _, _, _, _ uint64, cuIn int64) (r0 uint64, cuOut int64, err error) {
	return f(vm, r1, cuIn)
}

type SyscallFunc2 func(vm VM, r1, r2 uint64, cuIn int64) (r0 uint64, cuOut int64, err error)

func (f SyscallFunc2) Invoke(vm VM, r1, r2, _, _, _ uint64, cuIn int64) (r0 uint64, cuOut int64, err error) {
	return f(vm, r1, r2, cuIn)
}

type SyscallFunc3 func(vm VM, r1, r2, r3 uint64, cuIn int64) (r0 uint64, cuOut int64, err error)

func (f SyscallFunc3) Invoke(vm VM, r1, r2, r3, _, _ uint64, cuIn int64) (r0 uint64, cuOut int64, err error) {
	return f(vm, r1, r2, r3, cuIn)
}

type SyscallFunc4 func(vm VM, r1, r2, r3, r4 uint64, cuIn int64) (r0 uint64, cuOut int64, err error)

func (f SyscallFunc4) Invoke(vm VM, r1, r2, r3, r4, _ uint64, cuIn int64) (r0 uint64, cuOut int64, err error) {
	return f(vm, r1, r2, r3, r4, cuIn)
}

type SyscallFunc5 func(vm VM, r1, r2, r3, r4, r5 uint64, cuIn int64) (r0 uint64, cuOut int64, err error)

func (f SyscallFunc5) Invoke(vm VM, r1, r2, r3, r4, r5 uint64, cuIn int64) (r0 uint64, cuOut int64, err error) {
	return f(vm, r1, r2, r3, r4, r5, cuIn)
}
