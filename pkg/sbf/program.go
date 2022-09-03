package sbf

// Program is a loaded SBF program.
type Program struct {
	RO         []byte // read-only segment containing text and ELFs
	Text       []byte
	Entrypoint uint64 // PC
}

// Verify runs the static bytecode verifier.
func (p *Program) Verify() error {
	return NewVerifier(p).Verify()
}
