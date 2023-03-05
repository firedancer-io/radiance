package sealevel

import (
	"fmt"

	"github.com/gagliardetto/solana-go"
	"go.firedancer.io/radiance/pkg/sbpf"
	"go.firedancer.io/radiance/pkg/sbpf/cu"
)

func SyscallLogImpl(vm sbpf.VM, ptr, strlen uint64, cuIn int) (r0 uint64, cuOut int, err error) {
	if strlen > (1 << 30) {
		cuOut = -1
		return
	}
	cuOut = cu.ConsumeLowerBound(cuIn, CUSyscallBaseCost, int(strlen))
	if cuOut < 0 {
		return
	}

	buf := make([]byte, strlen)
	if err = vm.Read(ptr, buf); err != nil {
		return
	}
	syscallCtx(vm).Log.Log("Program log: " + string(buf))
	return
}

var SyscallLog = sbpf.SyscallFunc2(SyscallLogImpl)

func SyscallLog64Impl(vm sbpf.VM, r1, r2, r3, r4, r5 uint64, cuIn int) (r0 uint64, cuOut int, err error) {
	cuOut = cuIn - CUSyscallBaseCost
	if cuOut < 0 {
		return
	}

	msg := fmt.Sprintf("Program log: %#x, %#x, %#x, %#x, %#x\n", r1, r2, r3, r4, r5)
	syscallCtx(vm).Log.Log(msg)
	return
}

var SyscallLog64 = sbpf.SyscallFunc5(SyscallLog64Impl)

func SyscallLogCUsImpl(vm sbpf.VM, cuIn int) (r0 uint64, cuOut int, err error) {
	cuOut = cuIn - CUSyscallBaseCost
	if cuOut < 0 {
		return
	}

	msg := fmt.Sprintf("Program consumption: %d units remaining", cuIn)
	syscallCtx(vm).Log.Log(msg)
	return
}

var SyscallLogCUs = sbpf.SyscallFunc0(SyscallLogCUsImpl)

func SyscallLogPubkeyImpl(vm sbpf.VM, pubkeyAddr uint64, cuIn int) (r0 uint64, cuOut int, err error) {
	cuOut = cuIn - CUSyscallBaseCost
	if cuOut < 0 {
		return
	}

	// TODO alignment check
	var pubkey solana.PublicKey
	if err = vm.Read(pubkeyAddr, pubkey[:]); err != nil {
		return
	}

	syscallCtx(vm).Log.Log("Program log: " + pubkey.String())
	return
}

var SyscallLogPubkey = sbpf.SyscallFunc1(SyscallLogPubkeyImpl)
