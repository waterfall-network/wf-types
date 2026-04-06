// Copyright 2026 Digital Clever Solution Inc.
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

package rlp

import (
	"fmt"
	"math/big"
	"reflect"
)

// EncodeToBytes returns the RLP encoding of val.
func EncodeToBytes(val interface{}) ([]byte, error) {
	return encodeValue(reflect.ValueOf(val))
}

// encodeValue dispatches encoding by type.
func encodeValue(v reflect.Value) ([]byte, error) {
	// Dereference pointer
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			// nil pointer without context — caller should handle via struct field tags
			return []byte{0x80}, nil
		}
		return encodeValue(v.Elem())
	}

	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return encodeUint(v.Uint()), nil
	case reflect.Bool:
		return encodeBool(v.Bool()), nil
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return encodeByteSlice(v.Bytes()), nil
		}
		return encodeSlice(v)
	case reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			b := make([]byte, v.Len())
			reflect.Copy(reflect.ValueOf(b), v)
			return encodeByteSlice(b), nil
		}
		return encodeSlice(v)
	case reflect.Struct:
		if v.Type() == bigIntType {
			i := v.Interface().(big.Int)
			return encodeBigInt(&i), nil
		}
		return encodeStruct(v)
	default:
		return nil, fmt.Errorf("rlp: unsupported type %s", v.Type())
	}
}

var bigIntType = reflect.TypeOf(big.Int{})

// encodeUint encodes an unsigned integer as RLP.
// 0 → 0x80 (empty string), 1..0x7f → single byte, otherwise → string with big-endian bytes.
func encodeUint(i uint64) []byte {
	if i == 0 {
		return []byte{0x80}
	}
	if i < 128 {
		return []byte{byte(i)}
	}
	b := bigEndianBytes(i)
	return encodeString(b)
}

// encodeBool encodes a boolean as RLP (false=0x80, true=0x01).
func encodeBool(b bool) []byte {
	if b {
		return []byte{0x01}
	}
	return []byte{0x80}
}

// encodeBigInt encodes a *big.Int as RLP byte string (big-endian, no leading zeros).
// Zero and nil both encode as empty string 0x80.
func encodeBigInt(i *big.Int) []byte {
	if i == nil || i.Sign() == 0 {
		return []byte{0x80}
	}
	b := i.Bytes() // big-endian, no leading zeros
	return encodeString(b)
}

// encodeByteSlice encodes a []byte as RLP string.
func encodeByteSlice(b []byte) []byte {
	return encodeString(b)
}

// encodeString encodes raw bytes as RLP string.
func encodeString(b []byte) []byte {
	if len(b) == 1 && b[0] < 0x80 {
		return b
	}
	return append(encodeLength(len(b), 0x80), b...)
}

// encodeSlice encodes a slice (of non-byte elements) as RLP list.
func encodeSlice(v reflect.Value) ([]byte, error) {
	var payload []byte
	for i := 0; i < v.Len(); i++ {
		enc, err := encodeValue(v.Index(i))
		if err != nil {
			return nil, err
		}
		payload = append(payload, enc...)
	}
	return encodeList(payload), nil
}

// encodeStruct encodes a struct as an RLP list of its exported fields.
func encodeStruct(v reflect.Value) ([]byte, error) {
	fields, err := cachedStructFields(v.Type())
	if err != nil {
		return nil, err
	}

	// Find the last non-optional field index.
	lastRequired := -1
	for i, f := range fields {
		if !f.optional {
			lastRequired = i
		}
	}

	// Find the last optional field that has a non-zero value.
	// All optional fields up to and including it must be encoded,
	// even if some of them are zero (gwat/rlp compatibility: optional fields
	// are omitted only from the tail, never creating gaps in the list).
	lastIncluded := lastRequired
	for i, f := range fields {
		if f.optional && !isZeroValue(v.Field(f.index)) {
			lastIncluded = i
		}
	}

	var payload []byte
	for i, f := range fields {
		fv := v.Field(f.index)

		// Skip optional fields that trail beyond the last non-zero optional.
		if f.optional && i > lastIncluded {
			continue
		}

		enc, err := encodeField(fv, f)
		if err != nil {
			return nil, fmt.Errorf("rlp: field %d: %w", f.index, err)
		}
		payload = append(payload, enc...)
	}
	return encodeList(payload), nil
}

// encodeField encodes a single struct field, respecting tags.
func encodeField(v reflect.Value, f structField) ([]byte, error) {
	if f.tail {
		// For []byte tail fields: gwat/rlp ignores the tail tag during encoding
		// and encodes []byte as a regular RLP byte string (empty = 0x80).
		// For non-byte slices: each element is written individually (no list wrapper).
		if v.Kind() != reflect.Slice || v.Type().Elem().Kind() != reflect.Uint8 {
			return nil, fmt.Errorf("rlp: tail field must be []byte")
		}
		// Encode as a regular byte string — consistent with gwat/rlp behavior.
		return encodeByteSlice(v.Bytes()), nil
	}

	if f.nilOK && v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return nilEncoding(v.Type().Elem()), nil
		}
		return encodeValue(v.Elem())
	}

	return encodeValue(v)
}

// nilEncoding returns the RLP encoding for a nil pointer to the given element type.
// For byte-string types (integers, big.Int, byte arrays) → 0x80 (empty string).
// For struct types → 0xc0 (empty list).
func nilEncoding(elemType reflect.Type) []byte {
	switch elemType.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Bool, reflect.String:
		return []byte{0x80}
	default:
		// big.Int is treated as a byte string
		if elemType == bigIntType {
			return []byte{0x80}
		}
		if elemType.Kind() == reflect.Array && elemType.Elem().Kind() == reflect.Uint8 {
			return []byte{0x80}
		}
		if elemType.Kind() == reflect.Slice && elemType.Elem().Kind() == reflect.Uint8 {
			return []byte{0x80}
		}
		// struct or other composite → empty list
		return []byte{0xc0}
	}
}

// encodeList wraps payload bytes in an RLP list header.
func encodeList(payload []byte) []byte {
	return append(encodeLength(len(payload), 0xc0), payload...)
}

// encodeLength returns the length prefix for a string or list.
// smalltag is 0x80 for strings, 0xc0 for lists.
func encodeLength(length int, smalltag byte) []byte {
	if length < 56 {
		return []byte{smalltag + byte(length)}
	}
	b := bigEndianBytes(uint64(length))
	result := make([]byte, 1+len(b))
	result[0] = smalltag + 55 + byte(len(b))
	copy(result[1:], b)
	return result
}

// bigEndianBytes returns the minimal big-endian byte representation of n.
func bigEndianBytes(n uint64) []byte {
	var buf [8]byte
	buf[0] = byte(n >> 56)
	buf[1] = byte(n >> 48)
	buf[2] = byte(n >> 40)
	buf[3] = byte(n >> 32)
	buf[4] = byte(n >> 24)
	buf[5] = byte(n >> 16)
	buf[6] = byte(n >> 8)
	buf[7] = byte(n)
	i := 0
	for i < 7 && buf[i] == 0 {
		i++
	}
	return buf[i:]
}

// isZeroValue returns true if v is the zero value of its type.
func isZeroValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map:
		return v.IsNil()
	default:
		return v.IsZero()
	}
}
