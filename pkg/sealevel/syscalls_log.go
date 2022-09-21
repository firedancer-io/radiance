package sealevel

import (
	"fmt"

	"go.firedancer.io/radiance/pkg/sbf"
	"go.firedancer.io/radiance/pkg/sbf/cu"
	"github.com/gagliardetto/solana-go"
)

func SyscallLogImpl(vm sbf.VM, ptr, strlen uint64, cuIn int) (r0 uint64, cuOut int, err error) {
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

var SyscallLog = sbf.SyscallFunc2(SyscallLogImpl)

func SyscallLog64Impl(vm sbf.VM, r1, r2, r3, r4, r5 uint64, cuIn int) (r0 uint64, cuOut int, err error) {
	cuOut = cuIn - CUSyscallBaseCost
	if cuOut < 0 {
		return
	}

	msg := fmt.Sprintf("Program log: %#x, %#x, %#x, %#x, %#x\n", r1, r2, r3, r4, r5)
	syscallCtx(vm).Log.Log(msg)
	return
}

var SyscallLog64 = sbf.SyscallFunc5(SyscallLog64Impl)

func SyscallLogCUsImpl(vm sbf.VM, cuIn int) (r0 uint64, cuOut int, err error) {
	cuOut = cuIn - CUSyscallBaseCost
	if cuOut < 0 {
		return
	}

	msg := fmt.Sprintf("Program consumption: %d units remaining", cuIn)
	syscallCtx(vm).Log.Log(msg)
	return
}

var SyscallLogCUs = sbf.SyscallFunc0(SyscallLogCUsImpl)

func SyscallLogPubkeyImpl(vm sbf.VM, pubkeyAddr uint64, cuIn int) (r0 uint64, cuOut int, err error) {
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

var SyscallLogPubkey = sbf.SyscallFunc1(SyscallLogPubkeyImpl)
