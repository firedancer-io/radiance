package sealevel

import (
	"bytes"

	"go.firedancer.io/radiance/pkg/sbpf"
)

type TxContext struct{}

type Execution struct {
	Log Logger
}

func (t *TxContext) newVMOpts(params *Params) *sbpf.VMOpts {
	execution := &Execution{
		Log: new(LogRecorder),
	}
	var buf bytes.Buffer
	params.Serialize(&buf)
	return &sbpf.VMOpts{
		HeapSize: 32 * 1024,
		Syscalls: registry,
		Context:  execution,
		MaxCU:    1_400_000,
		Input:    buf.Bytes(),
	}
}
