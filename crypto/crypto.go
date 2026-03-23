// Copyright 2024 Blue Wave Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package crypto provides cryptographic utilities for wf-engine and wf-consensus.
// It implements the subset of gwat/crypto functions required by token, validator,
// and dag packages, using only Apache 2.0 / BSD-3-Clause dependencies.
package crypto

import (
	"encoding/binary"

	"golang.org/x/crypto/sha3"
)

// Keccak256 calculates and returns the Keccak256 hash of the concatenated input data.
func Keccak256(data ...[]byte) []byte {
	h := sha3.NewLegacyKeccak256()
	for _, b := range data {
		h.Write(b)
	}
	return h.Sum(nil)
}

// Keccak256Hash calculates the Keccak256 hash of the concatenated input data
// and returns it as a fixed-size 32-byte array.
func Keccak256Hash(data ...[]byte) [32]byte {
	var out [32]byte
	h := sha3.NewLegacyKeccak256()
	for _, b := range data {
		h.Write(b)
	}
	copy(out[:], h.Sum(nil))
	return out
}

// CreateAddress computes the Ethereum CREATE address from the sender address
// and nonce, using RLP encoding followed by Keccak256.
//
// Formula: keccak256(RLP([sender, nonce]))[12:]
//
// This is a self-contained implementation that does not depend on an external
// RLP library — the encoding of [address(20 bytes), uint64] is trivially fixed.
func CreateAddress(addr [20]byte, nonce uint64) [20]byte {
	data := rlpEncodeAddressNonce(addr, nonce)
	hash := Keccak256(data)
	var out [20]byte
	copy(out[:], hash[12:])
	return out
}

// rlpEncodeAddressNonce encodes the pair [addr, nonce] as an RLP list.
//
// RLP rules applied here:
//   - address (20-byte string): prefix 0x94 (= 0x80 + 20) followed by 20 bytes
//   - nonce (integer):
//   - 0          → 0x80 (empty byte string)
//   - 1..0x7f   → single byte (the value itself)
//   - otherwise  → 0x80+len prefix + minimal big-endian bytes
//   - list: 0xc0+total_payload_length prefix + concatenated items
//     (payload is always ≤ 30 bytes, so single-byte length is sufficient)
func rlpEncodeAddressNonce(addr [20]byte, nonce uint64) []byte {
	// Encode address: 0x94 + 20 bytes
	addrRLP := make([]byte, 21)
	addrRLP[0] = 0x94 // 0x80 + 20
	copy(addrRLP[1:], addr[:])

	// Encode nonce
	nonceRLP := rlpEncodeUint64(nonce)

	// Wrap in a list: 0xc0 + payload_length + payload
	payload := append(addrRLP, nonceRLP...)
	result := make([]byte, 1+len(payload))
	result[0] = byte(0xc0 + len(payload)) // payload <= 30, always fits in one byte
	copy(result[1:], payload)
	return result
}

// rlpEncodeUint64 returns the RLP encoding of a uint64 integer.
func rlpEncodeUint64(n uint64) []byte {
	if n == 0 {
		return []byte{0x80} // empty byte string
	}
	b := uint64ToMinimalBigEndian(n)
	if len(b) == 1 && b[0] < 0x80 {
		return b // single byte, value < 128: encoded as itself
	}
	// Multi-byte or value >= 0x80: 0x80+len prefix
	result := make([]byte, 1+len(b))
	result[0] = byte(0x80 + len(b))
	copy(result[1:], b)
	return result
}

// uint64ToMinimalBigEndian encodes n as big-endian bytes with leading zeros stripped.
func uint64ToMinimalBigEndian(n uint64) []byte {
	var buf [8]byte
	binary.BigEndian.PutUint64(buf[:], n)
	i := 0
	for i < 7 && buf[i] == 0 {
		i++
	}
	return buf[i:]
}
