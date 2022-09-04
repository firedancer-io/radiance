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
	Syscalls  map[uint32]Syscall

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
