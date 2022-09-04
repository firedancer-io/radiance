package sealevel

import "github.com/certusone/radiance/pkg/sbf"

var registry = Syscalls()

// Syscalls creates a registry of all Sealevel syscalls.
func Syscalls() sbf.SyscallRegistry {
	reg := sbf.NewSyscallRegistry()
	reg.Register("abort", SyscallAbort)
	reg.Register("sol_log_", SyscallLog)
	reg.Register("sol_log_64_", SyscallLog64)
	return reg
}

func syscallCtx(vm sbf.VM) *Execution {
	return vm.VMContext().(*Execution)
}
