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
	LegacyDataHeaderSize   = 86
	LegacyDataV2HeaderSize = 88
	LegacyPayloadSize      = 1057 // TODO where does this number come from?
	LegacyV2PayloadSize    = 1051 // TODO idk???
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

// LegacyData is genesis-era shred data.
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
	if len(shred) < LegacyDataHeaderSize {
		return nil
	}
	// TODO Sanitize
	data.Payload = make([]byte, LegacyDataHeaderSize+LegacyPayloadSize)
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
	return s.Payload[LegacyDataHeaderSize : LegacyDataHeaderSize+LegacyPayloadSize], true
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

// LegacyDataV2 is Q2 2020-era shred data.
type LegacyDataV2 struct {
	Common  CommonHeader
	Header  DataV2Header
	Payload []byte
}

func LegacyDataV2FromPayload(shred []byte) *LegacyDataV2 {
	data := new(LegacyDataV2)
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
	if len(shred) < LegacyDataV2HeaderSize {
		return nil
	}
	// TODO Sanitize
	data.Payload = make([]byte, LegacyDataV2HeaderSize+LegacyV2PayloadSize)
	copy(data.Payload, shred)
	return data
}

func (s *LegacyDataV2) CommonHeader() *CommonHeader {
	return &s.Common
}

func (s *LegacyDataV2) DataHeader() *DataV2Header {
	return &s.Header
}

/*[156, 65, 28, 59, 253, 19, 249]*/
/*9c 41 1c 3b fd 13 */

func (s *LegacyDataV2) Data() ([]byte, bool) {
	if int(s.Header.Size) > len(s.Payload) {
		return nil, false
	}
	return s.Payload[LegacyDataV2HeaderSize:s.Header.Size], true
}

func (s *LegacyDataV2) DataComplete() bool {
	return s.Header.Flags&FlagDataCompleteShred == 1
}

func (s *LegacyDataV2) ReferenceTick() uint8 {
	return s.Header.Flags & FlagShredTickReferenceMask
}

func (s *LegacyDataV2) MarshalYAML() (any, error) {
	item := struct {
		Signature   solana.Signature
		Variant     uint8
		Slot        uint64
		Index       uint32
		Version     uint16
		FECSetIndex uint32
		Size        uint16
		Payload     string
	}{
		Signature:   s.Common.Signature,
		Variant:     s.Common.Variant,
		Slot:        s.Common.Slot,
		Index:       s.Common.Index,
		Version:     s.Common.Version,
		FECSetIndex: s.Common.FECSetIndex,
		Size:        s.Header.Size,
		Payload:     base64.StdEncoding.EncodeToString(s.Payload),
	}
	return item, nil
}
