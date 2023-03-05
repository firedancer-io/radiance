package sealevel

import (
	"errors"

	"go.firedancer.io/radiance/pkg/sbpf"
)

func SyscallAbortImpl(_ sbpf.VM, _ int) (r0 uint64, cuOut int, err error) {
	err = errors.New("aborted")
	return
}

var SyscallAbort = sbpf.SyscallFunc0(SyscallAbortImpl)
