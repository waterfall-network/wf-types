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

// DecodeBytes decodes an RLP-encoded byte slice into val (must be a pointer).
func DecodeBytes(b []byte, val interface{}) error {
	rv := reflect.ValueOf(val)
	if rv.Kind() != reflect.Ptr || rv.IsNil() {
		return fmt.Errorf("rlp: DecodeBytes requires a non-nil pointer")
	}
	rest, err := decodeValue(b, rv.Elem())
	if err != nil {
		return err
	}
	if len(rest) > 0 {
		return ErrMoreThanOneValue
	}
	return nil
}

// decodeValue decodes the first RLP item from b into v and returns remaining bytes.
func decodeValue(b []byte, v reflect.Value) ([]byte, error) {
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			v.Set(reflect.New(v.Type().Elem()))
		}
		return decodeValue(b, v.Elem())
	}

	switch v.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint:
		return decodeUint(b, v)
	case reflect.Bool:
		return decodeBool(b, v)
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return decodeByteSlice(b, v)
		}
		return decodeSlice(b, v)
	case reflect.Array:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			return decodeByteArray(b, v)
		}
		return decodeSlice(b, v)
	case reflect.Struct:
		if v.Type() == bigIntType {
			return decodeBigInt(b, v)
		}
		return decodeStruct(b, v)
	default:
		return nil, fmt.Errorf("rlp: unsupported type %s", v.Type())
	}
}

// readString reads an RLP-encoded string from b.
// Returns the string bytes and remaining input.
func readString(b []byte) (content, rest []byte, err error) {
	if len(b) == 0 {
		return nil, nil, io_EOF()
	}
	tag := b[0]
	switch {
	case tag < 0x80:
		// Single byte
		return b[:1], b[1:], nil
	case tag <= 0xb7:
		// Short string
		size := int(tag - 0x80)
		if len(b) < 1+size {
			return nil, nil, fmt.Errorf("rlp: value size exceeds available input")
		}
		return b[1 : 1+size], b[1+size:], nil
	case tag <= 0xbf:
		// Long string
		sizeLen := int(tag - 0xb7)
		if len(b) < 1+sizeLen {
			return nil, nil, fmt.Errorf("rlp: value size exceeds available input")
		}
		size, err := decodeLength(b[1:1+sizeLen], sizeLen)
		if err != nil {
			return nil, nil, err
		}
		if len(b) < 1+sizeLen+size {
			return nil, nil, fmt.Errorf("rlp: value size exceeds available input")
		}
		return b[1+sizeLen : 1+sizeLen+size], b[1+sizeLen+size:], nil
	default:
		return nil, nil, ErrExpectedString
	}
}

// readList reads an RLP-encoded list header from b.
// Returns the list content bytes and remaining input.
func readList(b []byte) (content, rest []byte, err error) {
	if len(b) == 0 {
		return nil, nil, io_EOF()
	}
	tag := b[0]
	switch {
	case tag >= 0xc0 && tag <= 0xf7:
		// Short list
		size := int(tag - 0xc0)
		if len(b) < 1+size {
			return nil, nil, fmt.Errorf("rlp: list size exceeds available input")
		}
		return b[1 : 1+size], b[1+size:], nil
	case tag >= 0xf8:
		// Long list
		sizeLen := int(tag - 0xf7)
		if len(b) < 1+sizeLen {
			return nil, nil, fmt.Errorf("rlp: list size exceeds available input")
		}
		size, err := decodeLength(b[1:1+sizeLen], sizeLen)
		if err != nil {
			return nil, nil, err
		}
		if len(b) < 1+sizeLen+size {
			return nil, nil, fmt.Errorf("rlp: list size exceeds available input")
		}
		return b[1+sizeLen : 1+sizeLen+size], b[1+sizeLen+size:], nil
	default:
		return nil, nil, ErrExpectedList
	}
}

// decodeLength decodes a multi-byte big-endian length prefix.
func decodeLength(b []byte, n int) (int, error) {
	if n > 8 {
		return 0, fmt.Errorf("rlp: length too large")
	}
	var result uint64
	for i := 0; i < n; i++ {
		result = result<<8 | uint64(b[i])
	}
	return int(result), nil
}

func io_EOF() error { return fmt.Errorf("rlp: value size exceeds available input") }

// decodeUint decodes an RLP integer into v.
func decodeUint(b []byte, v reflect.Value) ([]byte, error) {
	content, rest, err := readString(b)
	if err != nil {
		return nil, err
	}
	switch len(content) {
	case 0:
		v.SetUint(0)
	case 1:
		if content[0] == 0 {
			return nil, ErrCanonInt
		}
		v.SetUint(uint64(content[0]))
	default:
		if content[0] == 0 {
			return nil, ErrCanonInt
		}
		if len(content) > 8 {
			return nil, ErrValueTooLarge
		}
		var n uint64
		for _, byt := range content {
			n = n<<8 | uint64(byt)
		}
		v.SetUint(n)
	}
	return rest, nil
}

// decodeBool decodes an RLP boolean into v.
func decodeBool(b []byte, v reflect.Value) ([]byte, error) {
	content, rest, err := readString(b)
	if err != nil {
		return nil, err
	}
	switch len(content) {
	case 0:
		v.SetBool(false)
	case 1:
		switch content[0] {
		case 0x01:
			v.SetBool(true)
		case 0x00:
			// canonical false is 0x80 (empty string), but 0x00 is also accepted
			v.SetBool(false)
		default:
			return nil, fmt.Errorf("rlp: invalid bool value 0x%02x", content[0])
		}
	default:
		return nil, fmt.Errorf("rlp: invalid bool encoding")
	}
	return rest, nil
}

// decodeByteSlice decodes an RLP string into a []byte.
func decodeByteSlice(b []byte, v reflect.Value) ([]byte, error) {
	content, rest, err := readString(b)
	if err != nil {
		return nil, err
	}
	dst := make([]byte, len(content))
	copy(dst, content)
	v.Set(reflect.ValueOf(dst))
	return rest, nil
}

// decodeByteArray decodes an RLP string into a [N]byte array.
func decodeByteArray(b []byte, v reflect.Value) ([]byte, error) {
	content, rest, err := readString(b)
	if err != nil {
		return nil, err
	}
	if len(content) != v.Len() {
		return nil, fmt.Errorf("rlp: input string has wrong length for [%d]byte, got %d", v.Len(), len(content))
	}
	reflect.Copy(v, reflect.ValueOf(content))
	return rest, nil
}

// decodeBigInt decodes an RLP string into a big.Int value.
func decodeBigInt(b []byte, v reflect.Value) ([]byte, error) {
	content, rest, err := readString(b)
	if err != nil {
		return nil, err
	}
	if len(content) > 0 && content[0] == 0 {
		return nil, ErrCanonInt
	}
	i := new(big.Int).SetBytes(content)
	v.Set(reflect.ValueOf(*i))
	return rest, nil
}

// decodeSlice decodes an RLP list into a slice or array.
func decodeSlice(b []byte, v reflect.Value) ([]byte, error) {
	content, rest, err := readList(b)
	if err != nil {
		return nil, err
	}
	elemType := v.Type().Elem()
	var elems []reflect.Value
	for len(content) > 0 {
		elem := reflect.New(elemType).Elem()
		content, err = decodeValue(content, elem)
		if err != nil {
			return nil, err
		}
		elems = append(elems, elem)
	}
	if v.Kind() == reflect.Array {
		if len(elems) != v.Len() {
			return nil, fmt.Errorf("rlp: array has wrong length %d, want %d", len(elems), v.Len())
		}
		for i, e := range elems {
			v.Index(i).Set(e)
		}
	} else {
		slice := reflect.MakeSlice(v.Type(), len(elems), len(elems))
		for i, e := range elems {
			slice.Index(i).Set(e)
		}
		v.Set(slice)
	}
	return rest, nil
}

// decodeStruct decodes an RLP list into a struct.
func decodeStruct(b []byte, v reflect.Value) ([]byte, error) {
	content, rest, err := readList(b)
	if err != nil {
		return nil, err
	}
	fields, err := cachedStructFields(v.Type())
	if err != nil {
		return nil, err
	}

	for _, f := range fields {
		fv := v.Field(f.index)

		if f.tail {
			// For []byte tail: gwat/rlp encodes it as a regular RLP string,
			// so decode it the same way (readString), not as raw remaining bytes.
			if len(content) == 0 {
				// no bytes left — leave as nil
				break
			}
			content, err = decodeByteSlice(content, fv)
			if err != nil {
				return nil, fmt.Errorf("rlp: tail field: %w", err)
			}
			break
		}

		if len(content) == 0 {
			if f.optional {
				// leave field at zero value
				continue
			}
			return nil, fmt.Errorf("rlp: too few elements for struct %s", v.Type())
		}

		if f.nilOK && fv.Kind() == reflect.Ptr {
			content, err = decodeNilablePtr(content, fv)
		} else {
			content, err = decodeValue(content, fv)
		}
		if err != nil {
			return nil, fmt.Errorf("rlp: field %d (%s): %w", f.index, v.Type().Field(f.index).Name, err)
		}
	}

	if len(content) > 0 {
		return nil, fmt.Errorf("rlp: %d unused bytes after struct", len(content))
	}
	return rest, nil
}

// decodeNilablePtr decodes a pointer field that may be nil.
// If the input is an empty string (0x80) or empty list (0xc0), leaves ptr nil.
// Otherwise allocates and decodes into the pointed-to value.
func decodeNilablePtr(b []byte, v reflect.Value) ([]byte, error) {
	if len(b) == 0 {
		return nil, io_EOF()
	}
	tag := b[0]
	elemType := v.Type().Elem()

	// Determine the "nil" encoding for this pointer type
	isNilString := isStringType(elemType)

	if isNilString && tag == 0x80 {
		// nil encoding for scalar/big.Int/byte pointer
		v.Set(reflect.Zero(v.Type()))
		return b[1:], nil
	}
	if !isNilString && tag == 0xc0 {
		// nil encoding for struct pointer
		v.Set(reflect.Zero(v.Type()))
		return b[1:], nil
	}

	// Non-nil: allocate and decode
	ptr := reflect.New(elemType)
	rest, err := decodeValue(b, ptr.Elem())
	if err != nil {
		return nil, err
	}
	v.Set(ptr)
	return rest, nil
}

// isStringType returns true if the type encodes as an RLP string (not a list).
func isStringType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uint,
		reflect.Bool, reflect.String:
		return true
	default:
		if t == bigIntType {
			return true
		}
		if t.Kind() == reflect.Array && t.Elem().Kind() == reflect.Uint8 {
			return true
		}
		if t.Kind() == reflect.Slice && t.Elem().Kind() == reflect.Uint8 {
			return true
		}
		return false
	}
}
