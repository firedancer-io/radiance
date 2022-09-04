package sealevel

import (
	_ "embed"
	"testing"

	"github.com/certusone/radiance/fixtures"
	"github.com/certusone/radiance/pkg/sbf"
	"github.com/certusone/radiance/pkg/sbf/loader"
	"github.com/stretchr/testify/require"
)

func TestInterpreter_Noop(t *testing.T) {
	// TODO simplify API?
	loader, err := loader.NewLoaderFromBytes(fixtures.SBF(t, "noop.so"))
	require.NotNil(t, loader)
	require.NoError(t, err)

	program, err := loader.Load()
	require.NotNil(t, program)
	require.NoError(t, err)

	require.NoError(t, program.Verify())

	syscalls := sbf.NewSyscallRegistry()
	syscalls.Register("log", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)

	interpreter := sbf.NewInterpreter(program, &sbf.VMOpts{
		StackSize: 1024,
		HeapSize:  1024, // TODO
		Input:     nil,
		MaxCU:     10000,
		Syscalls:  syscalls,
	})
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	require.NoError(t, err)
}
