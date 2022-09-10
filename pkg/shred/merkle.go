package shred

type MerkleCode struct {
	Common CommonHeader
}

func MerkleCodeFromPayload(shred []byte) *MerkleCode {
	panic("merkle shred code unimplemented")
}

func (s *MerkleCode) CommonHeader() *CommonHeader {
	return &s.Common
}

func (s *MerkleCode) DataHeader() *DataHeader {
	return nil
}

func (s *MerkleCode) Data() ([]byte, bool) {
	return nil, false
}

func (s *MerkleCode) DataComplete() bool {
	return false
}

type MerkleData struct {
	Common CommonHeader
	Header DataHeader
}

func MerkleDataFromPayload(shred []byte) *MerkleData {
	panic("merkle shred data unimplemented")
}

func (s *MerkleData) CommonHeader() *CommonHeader {
	return &s.Common
}

func (s *MerkleData) DataHeader() *DataHeader {
	return &s.Header
}

func (s *MerkleData) Data() ([]byte, bool) {
	panic("MerkleData.Data() unimplemented")
}

func (s *MerkleData) DataComplete() bool {
	return s.Header.Flags&FlagDataCompleteShred == 1
}
