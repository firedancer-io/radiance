package gossip

func MakeBitVecU8(bits []byte, len int) BitVecU8 {
	if len <= 0 {
		return BitVecU8{
			Bits: BitVecU8Inner{Value: nil},
			Len:  0,
		}
	}
	return BitVecU8{
		Bits: BitVecU8Inner{Value: &bits},
		Len:  uint64(len),
	}
}

func (bv *BitVecU8) Get(pos uint64) bool {
	if pos >= bv.Len {
		panic("get bit out of bounds")
	}
	return (*bv.Bits.Value)[pos/8]&(1<<(pos%8)) != 0
}

func (bv *BitVecU8) Set(pos uint64, b bool) {
	if pos >= bv.Len {
		panic("get bit out of bounds")
	}
	if b {
		(*bv.Bits.Value)[pos/8] |= 1 << (pos % 8)
	} else {
		(*bv.Bits.Value)[pos/8] &= ^uint8(1 << (pos % 8))
	}
}

func MakeBitVecU64(bits []uint64, len int) BitVecU64 {
	if len <= 0 {
		return BitVecU64{
			Bits: BitVecU64Inner{Value: nil},
			Len:  0,
		}
	}
	return BitVecU64{
		Bits: BitVecU64Inner{Value: &bits},
		Len:  uint64(len),
	}
}

func (bv *BitVecU64) Get(pos uint64) bool {
	if pos >= bv.Len {
		panic("get bit out of bounds")
	}
	return (*bv.Bits.Value)[pos/64]&(1<<(pos%64)) != 0
}

func (bv *BitVecU64) Set(pos uint64, b bool) {
	if pos >= bv.Len {
		panic("get bit out of bounds")
	}
	bits := *bv.Bits.Value
	if b {
		bits[pos/64] |= 1 << (pos % 64)
	} else {
		bits[pos/64] &= ^uint64(1 << (pos % 64))
	}
}
