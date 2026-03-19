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

// Package hexutil implements hex encoding with 0x prefix,
// as used by the Waterfall/Ethereum JSON-RPC API.
package hexutil

import (
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"strconv"
)

// Errors
var (
	ErrEmptyString   = errors.New("empty hex string")
	ErrSyntax        = errors.New("invalid hex string")
	ErrMissingPrefix = errors.New("hex string without 0x prefix")
	ErrOddLength     = errors.New("hex string of odd length")
	ErrEmptyNumber   = errors.New("hex string \"0x\"")
	ErrLeadingZero   = errors.New("hex number with leading zero digits")
	ErrUint64Range   = errors.New("hex number > 64 bits")
	ErrUint8Range    = errors.New("hex number > 8 bits")
	ErrBig256Range   = errors.New("hex number > 256 bits")
)

// Encode encodes b as a hex string with 0x prefix.
func Encode(b []byte) string {
	return "0x" + hex.EncodeToString(b)
}

// Decode decodes a hex string with 0x prefix.
func Decode(s string) ([]byte, error) {
	if len(s) == 0 {
		return nil, ErrEmptyString
	}
	if !has0xPrefix(s) {
		return nil, ErrMissingPrefix
	}
	s = s[2:]
	if len(s)%2 != 0 {
		return nil, ErrOddLength
	}
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil, ErrSyntax
	}
	return b, nil
}

// MustDecode decodes a hex string and panics on error.
func MustDecode(s string) []byte {
	b, err := Decode(s)
	if err != nil {
		panic(err)
	}
	return b
}

// EncodeBig encodes bigint as a hex string with 0x prefix.
// The sign of the integer is ignored.
func EncodeBig(bigint *big.Int) string {
	if bigint.Sign() == 0 {
		return "0x0"
	}
	return "0x" + bigint.Text(16)
}

// DecodeBig decodes a hex string with 0x prefix as a quantity.
// Numbers larger than 256 bits are not accepted.
func DecodeBig(input string) (*big.Int, error) {
	raw, err := checkNumberText([]byte(input))
	if err != nil {
		return nil, err
	}
	if len(raw) > 64 {
		return nil, ErrBig256Range
	}
	dec, ok := new(big.Int).SetString(string(raw), 16)
	if !ok {
		return nil, ErrSyntax
	}
	return dec, nil
}

// DecodeUint64 decodes a hex string with 0x prefix as a uint64.
func DecodeUint64(input string) (uint64, error) {
	raw, err := checkNumberText([]byte(input))
	if err != nil {
		return 0, err
	}
	if len(raw) > 16 {
		return 0, ErrUint64Range
	}
	dec, err := strconv.ParseUint(string(raw), 16, 64)
	if err != nil {
		return 0, ErrSyntax
	}
	return dec, nil
}

// UnmarshalFixedText decodes input as a hex string with 0x prefix.
// The length of out determines the required input length.
func UnmarshalFixedText(typname string, input, out []byte) error {
	raw, err := checkText(input, true)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return fmt.Errorf("empty hex string for %s", typname)
	}
	if len(raw)%2 != 0 {
		return ErrOddLength
	}
	if len(raw)/2 != len(out) {
		return fmt.Errorf("hex string has length %d, want %d for %s", len(raw), 2*len(out), typname)
	}
	_, err = hex.Decode(out, raw)
	return err
}

// UnmarshalFixedUnprefixedText decodes input as a hex string without 0x prefix.
func UnmarshalFixedUnprefixedText(typname string, input, out []byte) error {
	raw, err := checkText(input, false)
	if err != nil {
		return err
	}
	if len(raw) == 0 {
		return fmt.Errorf("empty hex string for %s", typname)
	}
	if len(raw)%2 != 0 {
		return ErrOddLength
	}
	if len(raw)/2 != len(out) {
		return fmt.Errorf("hex string has length %d, want %d for %s", len(raw), 2*len(out), typname)
	}
	_, err = hex.Decode(out, raw)
	return err
}

// UnmarshalFixedJSON decodes the input as a JSON string with 0x prefix.
func UnmarshalFixedJSON(typname string, input, out []byte) error {
	if !isString(input) {
		return fmt.Errorf("non-string JSON for %s", typname)
	}
	return UnmarshalFixedText(typname, input[1:len(input)-1], out)
}

// has0xPrefix returns true if s has a 0x or 0X prefix.
func has0xPrefix(s string) bool {
	return len(s) >= 2 && s[0] == '0' && (s[1] == 'x' || s[1] == 'X')
}

// isString checks if input is a JSON string.
func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// checkText strips the 0x prefix and validates the hex string.
func checkText(input []byte, wantPrefix bool) ([]byte, error) {
	if len(input) == 0 {
		return nil, nil
	}
	if wantPrefix {
		if !has0xPrefix(string(input)) {
			return nil, ErrMissingPrefix
		}
		input = input[2:]
	}
	return input, nil
}

// checkNumberText validates a hex number (with 0x prefix, no leading zeros).
func checkNumberText(input []byte) (raw []byte, err error) {
	if len(input) == 0 {
		return nil, ErrEmptyString
	}
	if !has0xPrefix(string(input)) {
		return nil, ErrMissingPrefix
	}
	input = input[2:]
	if len(input) == 0 {
		return nil, ErrEmptyNumber
	}
	if len(input) > 1 && input[0] == '0' {
		return nil, ErrLeadingZero
	}
	return input, nil
}
