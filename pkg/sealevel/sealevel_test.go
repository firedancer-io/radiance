package sealevel

import (
	_ "embed"
	"testing"

	"github.com/certusone/radiance/fixtures"
	"github.com/certusone/radiance/pkg/sbf"
	"github.com/certusone/radiance/pkg/sbf/loader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExecute_Memo(t *testing.T) {
	tx := TxContext{}
	opts := tx.newVMOpts(&Params{
		Accounts:  nil,
		Data:      []byte("Bla"),
		ProgramID: [32]byte{},
	})

	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sealevel", "MemoSq4gqABAXKb96qnH8TysNcWxMyWCqXgDLGmfcHr.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	interpreter := sbf.NewInterpreter(program, opts)
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	assert.NoError(t, err)

	// TODO expose proper API for this
	logs := opts.Context.(*Execution).Log.(*LogRecorder).Logs
	assert.Equal(t, logs, []string{
		`Memo (len 3): "Bla"`,
	})
}

func TestInterpreter_Noop(t *testing.T) {
	// TODO simplify API?
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbf", "noop.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbf.NewSyscallRegistry()
	syscalls.Register("log", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)

	var log LogRecorder

	interpreter := sbf.NewInterpreter(program, &sbf.VMOpts{
		HeapSize: 32 * 1024,
		Input:    nil,
		MaxCU:    10000,
		Syscalls: syscalls,
		Context:  &Execution{Log: &log},
	})
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	require.NoError(t, err)

	assert.Equal(t, log.Logs, []string{
		"entrypoint\x00",
		"0x1, 0x2, 0x3, 0x4, 0x5\n",
	})
}
