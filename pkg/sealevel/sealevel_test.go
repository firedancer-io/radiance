package sealevel

import (
	_ "embed"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"go.firedancer.io/radiance/fixtures"
	"go.firedancer.io/radiance/pkg/sbf"
	"go.firedancer.io/radiance/pkg/sbf/loader"
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
		`Program log: Memo (len 3): "Bla"`,
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
		"Program log: entrypoint\x00",
		"Program log: 0x1, 0x2, 0x3, 0x4, 0x5\n",
	})
}

type executeCase struct {
	Name    string
	Program string
	Params  Params
	Logs    []string
}

func (e *executeCase) run(t *testing.T) {
	ld, err := loader.NewLoaderFromBytes(fixtures.Load(t, e.Program))
	require.NoError(t, err)
	require.NotNil(t, ld)

	program, err := ld.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	tx := TxContext{}
	opts := tx.newVMOpts(&e.Params)
	opts.Tracer = testLogger{t}

	interpreter := sbf.NewInterpreter(program, opts)
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	assert.NoError(t, err)

	logs := opts.Context.(*Execution).Log.(*LogRecorder).Logs
	assert.Equal(t, logs, e.Logs)
}

func TestExecute(t *testing.T) {
	// Collect test cases
	var cases []executeCase
	err := filepath.WalkDir(fixtures.Path(t, "sealevel"), func(path string, entry fs.DirEntry, err error) error {
		if !entry.Type().IsRegular() ||
			!strings.HasPrefix(filepath.Base(path), "test_") ||
			filepath.Ext(path) != ".json" {
			return nil
		}

		buf, err := os.ReadFile(path)
		require.NoError(t, err, path)

		var _case executeCase
		require.NoError(t, json.Unmarshal(buf, &_case), path)

		cases = append(cases, _case)
		return nil
	})
	require.NoError(t, err)

	for i := range cases {
		_case := cases[i]
		t.Run(_case.Name, func(t *testing.T) {
			t.Parallel()
			_case.run(t)
		})
	}
}

type testLogger struct {
	t *testing.T
}

func (t testLogger) Printf(format string, args ...any) {
	t.t.Logf(format, args...)
}
