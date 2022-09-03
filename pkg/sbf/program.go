package sbf

// Program is a loaded SBF program.
type Program struct {
	RO []byte // read-only segment containing text and ELFs
}
