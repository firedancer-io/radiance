package sealevel

import (
	"bytes"

	"github.com/certusone/radiance/pkg/sbf"
)

type TxContext struct{}

type Execution struct {
	Log Logger
}

func (t *TxContext) newVMOpts(params *Params) *sbf.VMOpts {
	execution := &Execution{
		Log: new(LogRecorder),
	}
	var buf bytes.Buffer
	params.Serialize(&buf)
	return &sbf.VMOpts{
		HeapSize: 32 * 1024,
		Syscalls: registry,
		Context:  execution,
		MaxCU:    1_400_000,
		Input:    buf.Bytes(),
	}
}
