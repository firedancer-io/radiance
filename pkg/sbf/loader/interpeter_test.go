package loader

import (
	_ "embed"
	"testing"

	"github.com/certusone/radiance/pkg/sbf"
	"github.com/stretchr/testify/require"
)

func TestInterpreter_Noop(t *testing.T) {
	loader, err := NewLoaderFromBytes(soNoop)
	require.NotNil(t, loader)
	require.NoError(t, err)

	program, err := loader.Load()
	require.NotNil(t, program)
	require.NoError(t, err)

	require.NoError(t, program.Verify())

	interpreter := sbf.NewInterpreter(program, &sbf.VMOpts{
		StackSize: 1024,
		HeapSize:  1024, // TODO
		Input:     nil,
		MaxCU:     10000,
	})
	require.NotNil(t, interpreter)

	err = interpreter.Run()
	require.NoError(t, err)
}
