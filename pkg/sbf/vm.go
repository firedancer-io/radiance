package sbf

import "errors"

// VM is the virtual machine abstraction, implemented by each executor.
type VM interface {
	VMContext() any
	// TODO
}

// VMOpts specifies virtual machine parameters.
type VMOpts struct {
	StackSize int
	HeapSize  int
	Input     []byte // mapped at VaddrInput
	MaxCU     uint64
	Context   any // passed to syscalls
	Syscalls  map[uint32]Syscall
}

// Syscall are callback handles from VM to Go. (work in progress)
type Syscall interface {
	Invoke(vm VM, r1, r2, r3, r4, r5 uint64) (uint64, error)
}

// Exception codes.
var (
	ExcDivideByZero   = errors.New("division by zero")
	ExcDivideOverflow = errors.New("divide overflow")
)
