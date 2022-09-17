package shred

// VersionByMainnetSlot is a lazy hack to get shred versioning on mainnet to work.
func VersionByMainnetSlot(slot uint64) int {
	// No idea if this slot number is correct.
	// Testing from Triton One needed here.
	if slot < 12960000 {
		return 1
	}
	return 2
}
