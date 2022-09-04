package sealevel

import (
	"fmt"
	"strings"

	"github.com/certusone/radiance/pkg/sbf"
)

// TODO These are naive stubs

func SyscallLogImpl(vm sbf.VM, r1, r2 uint64, cuIn int64) (r0 uint64, cuOut int64, err error) {
	buf := make([]byte, r2)
	if err = vm.Read(r1, buf); err != nil {
		return
	}
	fmt.Println("Program Log:", strings.Trim(string(buf), " \t\x00"))
	//panic("log syscall unimplemented")
	return
}

var SyscallLog = sbf.SyscallFunc2(SyscallLogImpl)

func SyscallLog64Impl(vm sbf.VM, r1, r2, r3, r4, r5 uint64, cuIn int64) (r0 uint64, cuOut int64, err error) {
	fmt.Printf("Program Log: r1=%#x r2=%#x r3=%#x r4=%#x r5=%#x\n", r1, r2, r3, r4, r5)
	return
}

var SyscallLog64 = sbf.SyscallFunc5(SyscallLog64Impl)
