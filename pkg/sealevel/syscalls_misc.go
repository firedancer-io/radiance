package sealevel

import (
	"errors"

	"github.com/certusone/radiance/pkg/sbf"
)

func SyscallAbortImpl(_ sbf.VM, _ int) (r0 uint64, cuOut int, err error) {
	err = errors.New("aborted")
	return
}

var SyscallAbort = sbf.SyscallFunc0(SyscallAbortImpl)
