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

package hexutil

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
)

var (
	bytesT  = reflect.TypeOf(Bytes(nil))
	uint64T = reflect.TypeOf(Uint64(0))
	uint8T  = reflect.TypeOf(Uint8(0))
	bigT    = reflect.TypeOf((*Big)(nil))
)

// Bytes marshals/unmarshals as a JSON string with 0x prefix.
type Bytes []byte

// MarshalText implements encoding.TextMarshaler.
func (b Bytes) MarshalText() ([]byte, error) {
	result := make([]byte, len(b)*2+2)
	copy(result, "0x")
	hex.Encode(result[2:], b)
	return result, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Bytes) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return &json.UnmarshalTypeError{Value: "non-string", Type: bytesT}
	}
	return b.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Bytes) UnmarshalText(input []byte) error {
	raw, err := checkText(input, true)
	if err != nil {
		return err
	}
	dec := make([]byte, len(raw)/2)
	if _, err = hex.Decode(dec, raw); err != nil {
		return err
	}
	*b = dec
	return nil
}

// Uint64 marshals/unmarshals as a JSON string with 0x prefix.
type Uint64 uint64

// MarshalText implements encoding.TextMarshaler.
func (b Uint64) MarshalText() ([]byte, error) {
	buf := make([]byte, 2, 18)
	copy(buf, "0x")
	buf = strconv.AppendUint(buf, uint64(b), 16)
	return buf, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Uint64) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return &json.UnmarshalTypeError{Value: "non-string", Type: uint64T}
	}
	return b.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Uint64) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 16 {
		return ErrUint64Range
	}
	dec, err := strconv.ParseUint(string(raw), 16, 64)
	if err != nil {
		return ErrSyntax
	}
	*b = Uint64(dec)
	return nil
}

// String returns the hex encoding of b.
func (b Uint64) String() string {
	return fmt.Sprintf("%#x", uint64(b))
}

// Uint8 marshals/unmarshals as a JSON string with 0x prefix.
type Uint8 uint8

// MarshalText implements encoding.TextMarshaler.
func (b Uint8) MarshalText() ([]byte, error) {
	buf := make([]byte, 2, 6)
	copy(buf, "0x")
	buf = strconv.AppendUint(buf, uint64(b), 16)
	return buf, nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Uint8) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return &json.UnmarshalTypeError{Value: "non-string", Type: uint8T}
	}
	return b.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Uint8) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 2 {
		return ErrUint8Range
	}
	dec, err := strconv.ParseUint(string(raw), 16, 8)
	if err != nil {
		return ErrSyntax
	}
	*b = Uint8(dec)
	return nil
}

// Big marshals/unmarshals as a JSON string with 0x prefix.
// The sign of the integer is ignored.
type Big big.Int

// MarshalText implements encoding.TextMarshaler.
func (b Big) MarshalText() ([]byte, error) {
	return []byte(EncodeBig((*big.Int)(&b))), nil
}

// UnmarshalJSON implements json.Unmarshaler.
func (b *Big) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return &json.UnmarshalTypeError{Value: "non-string", Type: bigT}
	}
	return b.UnmarshalText(input[1 : len(input)-1])
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (b *Big) UnmarshalText(input []byte) error {
	raw, err := checkNumberText(input)
	if err != nil {
		return err
	}
	if len(raw) > 64 {
		return ErrBig256Range
	}
	dec, ok := new(big.Int).SetString(string(raw), 16)
	if !ok {
		return ErrSyntax
	}
	*b = (Big)(*dec)
	return nil
}

// ToInt converts b to a *big.Int.
func (b *Big) ToInt() *big.Int {
	return (*big.Int)(b)
}

// String returns the hex encoding of b.
func (b *Big) String() string {
	return EncodeBig(b.ToInt())
}
