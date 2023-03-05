package sealevel

import (
	"go.firedancer.io/radiance/pkg/sbpf"
)

var registry = Syscalls()

// Syscalls creates a registry of all Sealevel syscalls.
func Syscalls() sbpf.SyscallRegistry {
	reg := sbpf.NewSyscallRegistry()
	reg.Register("abort", SyscallAbort)
	reg.Register("sol_log_", SyscallLog)
	reg.Register("sol_log_64_", SyscallLog64)
	reg.Register("sol_log_compute_uits_", SyscallLogCUs)
	reg.Register("sol_log_pubkey", SyscallLogPubkey)
	return reg
}

func syscallCtx(vm sbpf.VM) *Execution {
	return vm.VMContext().(*Execution)
}
