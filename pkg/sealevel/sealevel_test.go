package sealevel

import (
	_ "embed"
	"encoding/json"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.firedancer.io/radiance/fixtures"
	"go.firedancer.io/radiance/pkg/sbpf"
	"go.firedancer.io/radiance/pkg/sbpf/loader"
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

	interpreter := sbpf.NewInterpreter(program, opts)
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
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbpf", "noop.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbpf.NewSyscallRegistry()
	syscalls.Register("log", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)

	var log LogRecorder

	interpreter := sbpf.NewInterpreter(program, &sbpf.VMOpts{
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

// The TestInterpreter_Memcpy_Strings_Match tests that memcpy works as expected
// by running an SBPF program that uses the memcpy syscall to copy a string
// literal to a stack buffer, before testing for equality using memcmp.
// The expected result is that the two match.
func TestInterpreter_Memcpy_Strings_Match(t *testing.T) {
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbpf", "memcpy_and_memmove_test_matched.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbpf.NewSyscallRegistry()
	syscalls.Register("sol_log_", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)
	syscalls.Register("my_copy", SyscallMemcpy)

	var log LogRecorder

	interpreter := sbpf.NewInterpreter(program, &sbpf.VMOpts{
		HeapSize: 32 * 1024,
		Input:    nil,
		MaxCU:    10000,
		Syscalls: syscalls,
		Context:  &Execution{Log: &log},
	})
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	assert.Equal(t, log.Logs, []string{
		"Program log: Strings matched after copy.",
	})
	require.NoError(t, err)
}

// The TestInterpreter_Memcpy_Do_Not_Match tests that memcpy works as expected
// by running an SBPF program that uses the memcpy syscall to copy a string
// literal to a stack buffer, with the destination then modified before testing
// for equality using memcmp. The expected result  is that the two do NOT match,
// because of the modification before comparison.
func TestInterpreter_Memcpy_Do_Not_Match(t *testing.T) {
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbpf", "memcpy_and_memmove_test_not_matched.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbpf.NewSyscallRegistry()
	syscalls.Register("sol_log_", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)
	syscalls.Register("my_copy", SyscallMemcpy)

	var log LogRecorder

	interpreter := sbpf.NewInterpreter(program, &sbpf.VMOpts{
		HeapSize: 32 * 1024,
		Input:    nil,
		MaxCU:    10000,
		Syscalls: syscalls,
		Context:  &Execution{Log: &log},
	})
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	assert.Equal(t, log.Logs, []string{
		"Program log: Strings did not match after copy.",
	})
	require.NoError(t, err)
}

// The TestInterpreter_Memmove_Strings_Match tests that memove works as expected
// by running an SBPF program that uses the memcpy syscall to copy a string
// literal to a stack buffer, before testing for equality using memcmp.
// The expected result is that the two match.
func TestInterpreter_Memmove_Strings_Match(t *testing.T) {
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbpf", "memcpy_and_memmove_test_matched.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbpf.NewSyscallRegistry()
	syscalls.Register("sol_log_", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)
	syscalls.Register("my_copy", SyscallMemmove)

	var log LogRecorder

	interpreter := sbpf.NewInterpreter(program, &sbpf.VMOpts{
		HeapSize: 32 * 1024,
		Input:    nil,
		MaxCU:    10000,
		Syscalls: syscalls,
		Context:  &Execution{Log: &log},
	})
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	assert.Equal(t, log.Logs, []string{
		"Program log: Strings matched after copy.",
	})
	require.NoError(t, err)
}

// The TestInterpreter_Memmove_Do_Not_Match function tests that memmove works
// as expected by running an SBPF program that uses the memcpy syscall to
// copy a string literal to a stack buffer, with the destination then
// modified before testing for equality using memcmp. The expected result is
// that the two do NOT match, because of the modification before comparison.
func TestInterpreter_Memmove_Do_Not_Match(t *testing.T) {
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbpf", "memcpy_and_memmove_test_not_matched.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbpf.NewSyscallRegistry()
	syscalls.Register("sol_log_", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)
	syscalls.Register("my_copy", SyscallMemmove)

	var log LogRecorder

	interpreter := sbpf.NewInterpreter(program, &sbpf.VMOpts{
		HeapSize: 32 * 1024,
		Input:    nil,
		MaxCU:    10000,
		Syscalls: syscalls,
		Context:  &Execution{Log: &log},
	})
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	assert.Equal(t, log.Logs, []string{
		"Program log: Strings did not match after copy.",
	})
	require.NoError(t, err)
}

// The TestInterpreter_Memcpy_Overlapping function tests that memcpy works
// as expected by attempting to do a copy involving two overlapping buffers.
// The expected result is an "Overlapping copy" error being returned.
func TestInterpreter_Memcpy_Overlapping(t *testing.T) {
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbpf", "memcpy_overlapping.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbpf.NewSyscallRegistry()
	syscalls.Register("sol_log_", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)
	syscalls.Register("my_copy", SyscallMemcpy)

	var log LogRecorder

	interpreter := sbpf.NewInterpreter(program, &sbpf.VMOpts{
		HeapSize: 32 * 1024,
		Input:    nil,
		MaxCU:    10000,
		Syscalls: syscalls,
		Context:  &Execution{Log: &log},
	})
	require.NotNil(t, interpreter)

	err = interpreter.Run()

	// expecting an error here because the src and dst are overlapping in the
	// program being run.
	require.Error(t, err)
}

// The TestInterpreter_Memcmp_Matches function tests that the memcmp
// syscall works as expected by comparing two instances of "abcdabcd1234"
// The expected result is that the two strings match and the program
// writes "Memory chunks matched." to the program log.
func TestInterpreter_Memcmp_Matches(t *testing.T) {
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbpf", "memcmp_matched.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbpf.NewSyscallRegistry()
	syscalls.Register("sol_log_", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)
	syscalls.Register("my_memcmp", SyscallMemcmp)

	var log LogRecorder

	interpreter := sbpf.NewInterpreter(program, &sbpf.VMOpts{
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
		"Program log: Memory chunks matched.",
	})
}

// The TestInterpreter_Memcmp_Does_Not_Match function tests that the memcmp
// syscall works as expected by comparing the string literals "Bbcdabcd1234"
// and "aLAHabcd1234"
// The expected result is that the two strings do not match and the difference
// between the first non-matching characters (0x42 - 0x61 = -0x1f) is returned,
// and the program checks these and returns messages accordingly.
func TestInterpreter_Memcmp_Does_Not_Match(t *testing.T) {
	loader, err := loader.NewLoaderFromBytes(fixtures.Load(t, "sbpf", "memcmp_not_matched.so"))
	require.NoError(t, err)
	require.NotNil(t, loader)

	program, err := loader.Load()
	require.NoError(t, err)
	require.NotNil(t, program)

	require.NoError(t, program.Verify())

	syscalls := sbpf.NewSyscallRegistry()
	syscalls.Register("sol_log_", SyscallLog)
	syscalls.Register("log_64", SyscallLog64)
	syscalls.Register("my_memcmp", SyscallMemcmp)

	var log LogRecorder

	interpreter := sbpf.NewInterpreter(program, &sbpf.VMOpts{
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
		"Program log: Memory chunks did not match.",
		"Program log: Difference between non-matching character was correctly returned.",
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

	interpreter := sbpf.NewInterpreter(program, opts)
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
