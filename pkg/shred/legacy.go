package shred

import (
	"encoding/base64"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type LegacyCode struct {
	Common  CommonHeader
	Payload []byte
}

const (
	LegacyHeaderSize  = 86
	LegacyPayloadSize = 1228
)

func LegacyCodeFromPayload(shred []byte) *LegacyCode {
	panic("legacy shred code unimplemented")
}

func (s *LegacyCode) CommonHeader() *CommonHeader {
	return &s.Common
}

func (s *LegacyCode) DataHeader() *DataHeader {
	return nil
}

func (s *LegacyCode) Data() ([]byte, bool) {
	return nil, false
}

func (s *LegacyCode) DataComplete() bool {
	return false
}

type LegacyData struct {
	Common  CommonHeader
	Header  DataHeader
	Payload []byte
}

func LegacyDataFromPayload(shred []byte) *LegacyData {
	data := new(LegacyData)
	dec := bin.NewBinDecoder(shred)
	if err := dec.Decode(&data.Common); err != nil {
		return nil
	}
	if err := dec.Decode(&data.Header); err != nil {
		return nil
	}
	if data.Common.Variant != LegacyDataID {
		return nil
	}
	if len(shred) < LegacyHeaderSize {
		return nil
	}
	// TODO Sanitize
	data.Payload = make([]byte, LegacyPayloadSize)
	copy(data.Payload, shred)
	return data
}

func (s *LegacyData) CommonHeader() *CommonHeader {
	return &s.Common
}

func (s *LegacyData) DataHeader() *DataHeader {
	return &s.Header
}

func (s *LegacyData) Data() ([]byte, bool) {
	return s.Payload[LegacyHeaderSize:1143], true
}

func (s *LegacyData) DataComplete() bool {
	return s.Header.Flags&FlagDataCompleteShred == 1
}

func (s *LegacyData) ReferenceTick() uint8 {
	return s.Header.Flags & FlagShredTickReferenceMask
}

func (s *LegacyData) MarshalYAML() (any, error) {
	item := struct {
		Signature   solana.Signature
		Variant     uint8
		Slot        uint64
		Index       uint32
		Version     uint16
		FECSetIndex uint32
		Payload     string
	}{
		Signature:   s.Common.Signature,
		Variant:     s.Common.Variant,
		Slot:        s.Common.Slot,
		Index:       s.Common.Index,
		Version:     s.Common.Version,
		FECSetIndex: s.Common.FECSetIndex,
		Payload:     base64.StdEncoding.EncodeToString(s.Payload),
	}
	return item, nil
}
