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

// Package rlp implements RLP (Recursive Length Prefix) encoding and decoding.
//
// RLP is the serialization format used in Ethereum. This implementation is
// binary-compatible with go-ethereum's rlp package and supports the types
// used in wf-engine's token and validator packages.
//
// Supported Go types:
//   - Unsigned integers: uint8, uint16, uint32, uint64, uint
//   - Boolean: bool (encoded as 0x00 or 0x01)
//   - Byte slices: []byte (encoded as RLP string)
//   - Fixed-size byte arrays: [N]byte (encoded as RLP string)
//   - math/big.Int by value and by pointer (*big.Int)
//   - Slices of supported types
//   - Structs (encoded as RLP list of fields)
//
// Struct tags:
//   - `rlp:"-"` — field is ignored
//   - `rlp:"nil"` — pointer field may be nil; nil encodes as empty string (0x80)
//     for scalar types or empty list (0xc0) for struct types
//   - `rlp:"optional"` — field may be omitted at end of list; zero value on decode
//     if absent; only allowed at the end of a struct
//   - `rlp:"tail"` — field receives remaining undecoded bytes as []byte;
//     must be the last field and of type []byte
package rlp

import "errors"

// ErrCanonInt is returned when an integer has leading zero bytes.
var ErrCanonInt = errors.New("rlp: non-canonical integer format")

// ErrCanonSize is returned when a string or list has a non-minimal length prefix.
var ErrCanonSize = errors.New("rlp: non-canonical size information")

// ErrExpectedString is returned when a string value was expected but a list was found.
var ErrExpectedString = errors.New("rlp: expected String or Byte, got List")

// ErrExpectedList is returned when a list was expected but a string was found.
var ErrExpectedList = errors.New("rlp: expected List")

// ErrElemTooLarge is returned when an element is larger than the containing list.
var ErrElemTooLarge = errors.New("rlp: element is larger than containing list")

// ErrValueTooLarge is returned when the value does not fit the target type.
var ErrValueTooLarge = errors.New("rlp: value too large for type")

// ErrMoreThanOneValue is returned when the input contains more than one value.
var ErrMoreThanOneValue = errors.New("rlp: input contains more than one value")

// inputError wraps an error with position information.
type inputError struct {
	msg string
	err error
}

func (e *inputError) Error() string { return e.msg }
func (e *inputError) Unwrap() error { return e.err }
