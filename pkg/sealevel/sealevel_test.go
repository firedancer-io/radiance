package sealevel

import (
	_ "embed"
	"log"
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

func TestExecute_Token(t *testing.T) {
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sealevel", "TokenkegQfeZyiNwAJbNbGKPFXCWuBvf9Ss623VQ5DA.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	tx := TxContext{}
	opts := tx.newVMOpts(&Params{
		Accounts: []AccountParam{
			// Mint
			{
				IsDuplicate:    false,
				DuplicateIndex: 0xFF,
				IsSigner:       true,
				IsWritable:     true,
				IsExecutable:   false,
				Key:            [32]byte{0x1},
				Owner:          [32]byte{6, 221, 246, 225, 215, 101, 161, 147, 217, 203, 225, 70, 206, 235, 121, 172, 28, 180, 133, 237, 95, 91, 55, 145, 58, 140, 245, 133, 126, 255, 0, 169},
				Lamports:       10000000,
				Data:           make([]byte, 82),
				RentEpoch:      0,
			},
			// Rent sysvar
			{
				IsDuplicate:    false,
				DuplicateIndex: 0xFF,
				IsSigner:       true,
				IsWritable:     true,
				IsExecutable:   false,
				Key:            [32]byte{0x06, 0xa7, 0xd5, 0x17, 0x19, 0x2c, 0x5c, 0x51, 0x21, 0x8c, 0xc9, 0x4c, 0x3d, 0x4a, 0xf1, 0x7f, 0x58, 0xda, 0xee, 0x08, 0x9b, 0xa1, 0xfd, 0x44, 0xe3, 0xdb, 0xd9, 0x8a, 0x00, 0x00, 0x00, 0x00},
				Owner:          [32]byte{0x06, 0xa7, 0xd5, 0x17, 0x18, 0x75, 0xf7, 0x29, 0xc7, 0x3d, 0x93, 0x40, 0x8f, 0x21, 0x61, 0x20, 0x06, 0x7e, 0xd8, 0x8c, 0x76, 0xe0, 0x8c, 0x28, 0x7f, 0xc1, 0x94, 0x60, 0x00, 0x00, 0x00, 0x00},
				Lamports:       10092,
				Data: []byte{
					0x98, 0x0d, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
					0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x40,
					0x64,
				},
				RentEpoch: 0,
			},
		},
		Data: []byte{
			0, 0, 0, 0,
			1,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
			0,
			0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
		},
		ProgramID: [32]byte{},
	})

	logger := log.Default()
	logger.SetFlags(0)
	opts.Tracer = logger

	interpreter := sbf.NewInterpreter(program, opts)
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	assert.NoError(t, err)

	logs := opts.Context.(*Execution).Log.(*LogRecorder).Logs
	assert.Equal(t, logs, []string{
		"Instruction: InitializeMint",
	})
}
