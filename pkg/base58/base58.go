// Copyright 2022 Firedancer Contributors

// Package base58 converts between binary and Base58 for 32/64 length strings.
//
// Ported from Firedancer:
// https://github.com/firedancer-io/firedancer/blob/main/src/ballet/base58/fd_base58.h
//
// Original author: Philip Taffet <phtaffet@jumptrading.com>
package base58

import "encoding/binary"

// alphabet maps [0, 58) to the base58 character.
const alphabet = "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"

// encTable32 contains the unique values less than 58^5 such that:
//
//	2^(32*(7-j)) = sum_k table[j][k]*58^(5*(7-k))
//
// The second dimension of this table is actually ceil(log_(58^5)
// (2^(32*7)), but that's almost always 8
var encTable32 = [9][8]uint32{
	{513735, 77223048, 437087610, 300156666, 605448490, 214625350, 141436834, 379377856},
	{0, 78508, 646269101, 118408823, 91512303, 209184527, 413102373, 153715680},
	{0, 0, 11997, 486083817, 3737691, 294005210, 247894721, 289024608},
	{0, 0, 0, 1833, 324463681, 385795061, 551597588, 21339008},
	{0, 0, 0, 0, 280, 127692781, 389432875, 357132832},
	{0, 0, 0, 0, 0, 42, 537767569, 410450016},
	{0, 0, 0, 0, 0, 0, 6, 356826688},
	{0, 0, 0, 0, 0, 0, 0, 1},
}

func Encode32(out *[44]byte, in [32]byte) uint {
	const raw58sz = 45

	// Count leading zeros (needed for final output)
	var inLeading0s uint
	for i := range in {
		if in[i] != 0 {
			break
		}
		inLeading0s++
	}

	// X = sum_i bytes[i] * 2^(8*(32-1-i))

	// Convert N to 32-bit limbs:
	// X = sum_i binary[i] * 2^(32*(8-1-i))
	var limbs [8]uint32
	for i := range limbs {
		limbs[i] = binary.BigEndian.Uint32(in[4*i:])
	}

	r1Div := uint64(656356768) // = 58^5

	// Convert to the intermediate format:
	//   X = sum_i intermediate[i] * 58^(5*(INTERMEDIATE_SZ-1-i))
	// Initially, we don't require intermediate[i] < 58^5,
	// but we do want to make sure the sums don't overflow.

	var intermediate [9]uint64

	// The worst case is if binary[7] is (2^32)-1. In that case
	// intermediate[8] will be be just over 2^63, which is fine.

	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			intermediate[j+1] += uint64(limbs[i]) * uint64(encTable32[i][j])
		}
	}

	// Now we make sure each term is less than 58^5.
	// Again, we have to be a bit careful of overflow.
	//
	// For N==32, in the worst case, as before, intermediate[8] will be
	// just over 2^63 and intermediate[7] will be just over 2^62.6.  In
	// the first step, we'll add floor(intermediate[8]/58^5) to
	// intermediate[7].  58^5 is pretty big though, so intermediate[7]
	// barely budges, and this is still fine.
	//
	// For N==64, in the worst case, the biggest entry in intermediate
	// at this point is 2^63.87, and in the worst case, we add
	// (2^64-1)/58^5, which is still about 2^63.87.

	for i := 8; i > 0; i-- {
		intermediate[i-1] += intermediate[i] / r1Div
		intermediate[i] %= r1Div
	}

	// Convert intermediate form to base 58.
	//   X = sum_i raw_base58[i] * 58^(RAW58_SZ-1-i)

	var rawBase58 [45]byte
	for i := 0; i < 9; i++ {
		// We know intermediate[ i ] < 58^5 < 2^32 for all i, so casting
		// to a uint32 is safe.
		v := uint32(intermediate[i])
		rawBase58[5*i+4] = byte((v / 1) % 58)
		rawBase58[5*i+3] = byte((v / 58) % 58)
		rawBase58[5*i+2] = byte((v / 3364) % 58)
		rawBase58[5*i+1] = byte((v / 195112) % 58)
		// We know this one is less than 58
		rawBase58[5*i+0] = byte(v / 11316496)
	}

	// Finally, actually convert to the string.
	// We have to ignore all the leading zeros in rawBase58 and instead
	// insert inLeading0s leading '1' characters.  We can show that
	// rawBase58 actually has at least inLeading0s, so we'll do this
	// by skipping the first few leading zeros in rawBase58.

	var rawLeading0s uint
	for rawLeading0s = 0; rawLeading0s < raw58sz; rawLeading0s++ {
		if rawBase58[rawLeading0s] != 0 {
			break
		}
	}

	// It's not immediately obvious that rawLeading0s >= inLeading0s,
	// but it's true.  In base b, X has floor(log_b X)+1 digits.  That
	// means inLeading0s = N-1-floor(log_256 X) and rawLeading0s =
	// RAW58_SZ-1-floor(log_58 X).  Let X<256^N be given and consider:
	//
	//   rawLeading0s - inLeading0s =
	//     =  RAW58_SZ-N + floor( log_256 X ) - floor( log_58 X )
	//     >= RAW58_SZ-N - 1 + ( log_256 X - log_58 X ) .
	//
	// log_256 X - log_58 X is monotonically decreasing for X>0, so it
	// achieves it minimum at the maximum possible value for X, i.e.
	// 256^N-1.
	//   >= RAW58_SZ-N-1 + log_256(256^N-1) - log_58(256^N-1)
	//
	// When N==32, RAW58_SZ is 45, so this gives skip >= 0.29
	// When N==64, RAW58_SZ is 90, so this gives skip >= 1.59.
	//
	// Regardless, rawLeading0s - inLeading0s >= 0.

	skip := rawLeading0s - inLeading0s
	for i := uint(0); i < raw58sz-skip; i++ {
		out[i] = alphabet[rawBase58[i+skip]]
	}

	return raw58sz - skip
}
