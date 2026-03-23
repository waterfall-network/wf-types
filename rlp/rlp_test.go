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

package rlp

import (
	"bytes"
	"math/big"
	"testing"
)

// --- Primitives ---

func TestEncodeDecodeUint(t *testing.T) {
	cases := []struct {
		name string
		val  uint64
		enc  []byte
	}{
		{"zero", 0, []byte{0x80}},
		{"one", 1, []byte{0x01}},
		{"127", 127, []byte{0x7f}},
		{"128", 128, []byte{0x81, 0x80}},
		{"255", 255, []byte{0x81, 0xff}},
		{"256", 256, []byte{0x82, 0x01, 0x00}},
		{"65535", 65535, []byte{0x82, 0xff, 0xff}},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := encodeUint(c.val)
			if !bytes.Equal(got, c.enc) {
				t.Fatalf("encodeUint(%d) = %x, want %x", c.val, got, c.enc)
			}
			var out uint64
			if err := DecodeBytes(c.enc, &out); err != nil {
				t.Fatalf("DecodeBytes: %v", err)
			}
			if out != c.val {
				t.Fatalf("round-trip: got %d, want %d", out, c.val)
			}
		})
	}
}

func TestEncodeDecodeBool(t *testing.T) {
	trueEnc := encodeBool(true)
	falseEnc := encodeBool(false)
	if !bytes.Equal(trueEnc, []byte{0x01}) {
		t.Fatalf("encodeBool(true) = %x, want 01", trueEnc)
	}
	if !bytes.Equal(falseEnc, []byte{0x80}) {
		t.Fatalf("encodeBool(false) = %x, want 80", falseEnc)
	}

	var b bool
	if err := DecodeBytes(trueEnc, &b); err != nil || !b {
		t.Fatalf("decode true: err=%v b=%v", err, b)
	}
	if err := DecodeBytes(falseEnc, &b); err != nil || b {
		t.Fatalf("decode false: err=%v b=%v", err, b)
	}
}

func TestEncodeDecodeByteSlice(t *testing.T) {
	cases := [][]byte{
		nil,
		{},
		{0x01},
		{0x7f},
		{0x80},
		{0x00, 0x01, 0x02},
		make([]byte, 60),
	}
	for _, c := range cases {
		enc, err := EncodeToBytes(c)
		if err != nil {
			t.Fatalf("encode %x: %v", c, err)
		}
		var out []byte
		if err := DecodeBytes(enc, &out); err != nil {
			t.Fatalf("decode %x: %v", enc, err)
		}
		if !bytes.Equal(out, c) {
			t.Fatalf("round-trip %x: got %x", c, out)
		}
	}
}

func TestEncodeDecodeByteArray(t *testing.T) {
	var arr [4]byte
	arr[0] = 0xde
	arr[1] = 0xad
	arr[2] = 0xbe
	arr[3] = 0xef
	enc, err := EncodeToBytes(arr)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var out [4]byte
	if err := DecodeBytes(enc, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out != arr {
		t.Fatalf("round-trip: got %x, want %x", out, arr)
	}
}

func TestEncodeDecodeBigInt(t *testing.T) {
	cases := []*big.Int{
		new(big.Int), // zero
		big.NewInt(1),
		big.NewInt(127),
		big.NewInt(128),
		big.NewInt(256),
		new(big.Int).Lsh(big.NewInt(1), 128), // very large
	}
	for _, c := range cases {
		enc, err := EncodeToBytes(*c)
		if err != nil {
			t.Fatalf("encode %s: %v", c, err)
		}
		var out big.Int
		if err := DecodeBytes(enc, &out); err != nil {
			t.Fatalf("decode: %v", err)
		}
		if c.Cmp(&out) != 0 {
			t.Fatalf("round-trip %s: got %s", c, &out)
		}
	}
}

// --- Structs ---

type simpleStruct struct {
	A uint64
	B []byte
	C [4]byte
}

func TestEncodeDecodeSimpleStruct(t *testing.T) {
	in := simpleStruct{
		A: 42,
		B: []byte{0x01, 0x02, 0x03},
		C: [4]byte{0xaa, 0xbb, 0xcc, 0xdd},
	}
	enc, err := EncodeToBytes(in)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var out simpleStruct
	if err := DecodeBytes(enc, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.A != in.A || !bytes.Equal(out.B, in.B) || out.C != in.C {
		t.Fatalf("round-trip mismatch: got %+v, want %+v", out, in)
	}
}

type nilStruct struct {
	A uint64
	B *big.Int  `rlp:"nil"`
	C *[32]byte `rlp:"nil"`
}

func TestEncodeDecodeNilFields(t *testing.T) {
	// nil pointer fields
	in := nilStruct{A: 7, B: nil, C: nil}
	enc, err := EncodeToBytes(in)
	if err != nil {
		t.Fatalf("encode nil fields: %v", err)
	}
	var out nilStruct
	if err := DecodeBytes(enc, &out); err != nil {
		t.Fatalf("decode nil fields: %v", err)
	}
	if out.A != in.A || out.B != nil || out.C != nil {
		t.Fatalf("nil round-trip: got %+v", out)
	}

	// non-nil pointer fields
	b := big.NewInt(999)
	var arr [32]byte
	arr[0] = 0xff
	in2 := nilStruct{A: 1, B: b, C: &arr}
	enc2, err := EncodeToBytes(in2)
	if err != nil {
		t.Fatalf("encode non-nil: %v", err)
	}
	var out2 nilStruct
	if err := DecodeBytes(enc2, &out2); err != nil {
		t.Fatalf("decode non-nil: %v", err)
	}
	if out2.B == nil || out2.B.Cmp(b) != 0 {
		t.Fatalf("non-nil B: got %v, want %v", out2.B, b)
	}
	if out2.C == nil || *out2.C != arr {
		t.Fatalf("non-nil C: got %v, want %v", out2.C, arr)
	}
}

type optionalStruct struct {
	A uint64
	B uint64 `rlp:"optional"`
	C []byte `rlp:"optional"`
}

func TestEncodeDecodeOptional(t *testing.T) {
	// with optional fields present
	in := optionalStruct{A: 1, B: 2, C: []byte{0x03}}
	enc, err := EncodeToBytes(in)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var out optionalStruct
	if err := DecodeBytes(enc, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.A != 1 || out.B != 2 || !bytes.Equal(out.C, []byte{0x03}) {
		t.Fatalf("optional full: got %+v", out)
	}

	// encode with zero optional fields — they should be omitted
	in2 := optionalStruct{A: 5}
	enc2, err := EncodeToBytes(in2)
	if err != nil {
		t.Fatalf("encode zero optionals: %v", err)
	}
	var out2 optionalStruct
	if err := DecodeBytes(enc2, &out2); err != nil {
		t.Fatalf("decode zero optionals: %v", err)
	}
	if out2.A != 5 || out2.B != 0 || out2.C != nil {
		t.Fatalf("optional zero: got %+v", out2)
	}
}

type tailStruct struct {
	A    uint64
	Rest []byte `rlp:"tail"`
}

func TestEncodeDecodeTail(t *testing.T) {
	in := tailStruct{A: 10, Rest: []byte{0x01, 0x02, 0x03, 0x04}}
	enc, err := EncodeToBytes(in)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var out tailStruct
	if err := DecodeBytes(enc, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.A != in.A || !bytes.Equal(out.Rest, in.Rest) {
		t.Fatalf("tail round-trip: got %+v, want %+v", out, in)
	}

	// empty tail
	in2 := tailStruct{A: 3}
	enc2, err := EncodeToBytes(in2)
	if err != nil {
		t.Fatalf("encode empty tail: %v", err)
	}
	var out2 tailStruct
	if err := DecodeBytes(enc2, &out2); err != nil {
		t.Fatalf("decode empty tail: %v", err)
	}
	if out2.A != 3 || len(out2.Rest) != 0 {
		t.Fatalf("empty tail: got %+v", out2)
	}
}

// --- Slices ---

func TestEncodeDecodeSliceOfStructs(t *testing.T) {
	in := []simpleStruct{
		{A: 1, B: []byte{0xaa}, C: [4]byte{1, 2, 3, 4}},
		{A: 2, B: []byte{0xbb}, C: [4]byte{5, 6, 7, 8}},
	}
	enc, err := EncodeToBytes(in)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var out []simpleStruct
	if err := DecodeBytes(enc, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(out) != len(in) {
		t.Fatalf("len: got %d, want %d", len(out), len(in))
	}
	for i := range in {
		if out[i].A != in[i].A || !bytes.Equal(out[i].B, in[i].B) || out[i].C != in[i].C {
			t.Fatalf("elem %d: got %+v, want %+v", i, out[i], in[i])
		}
	}
}

// --- Ignored fields ---

type ignoredStruct struct {
	A uint64
	B uint64 `rlp:"-"`
	C uint64
}

func TestIgnoredField(t *testing.T) {
	in := ignoredStruct{A: 1, B: 999, C: 3}
	enc, err := EncodeToBytes(in)
	if err != nil {
		t.Fatalf("encode: %v", err)
	}
	var out ignoredStruct
	if err := DecodeBytes(enc, &out); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if out.A != 1 || out.B != 0 || out.C != 3 {
		t.Fatalf("ignored: got %+v", out)
	}
}

// --- Known RLP vectors (Ethereum Yellow Paper / go-ethereum) ---

func TestKnownVectors(t *testing.T) {
	// dog → [0x83, 'd', 'o', 'g']
	enc, err := EncodeToBytes([]byte("dog"))
	if err != nil {
		t.Fatalf("encode dog: %v", err)
	}
	want := []byte{0x83, 'd', 'o', 'g'}
	if !bytes.Equal(enc, want) {
		t.Fatalf("dog: got %x, want %x", enc, want)
	}

	// ["cat","dog"] → [0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g']
	enc2, err := EncodeToBytes([][]byte{[]byte("cat"), []byte("dog")})
	if err != nil {
		t.Fatalf("encode [cat,dog]: %v", err)
	}
	want2 := []byte{0xc8, 0x83, 'c', 'a', 't', 0x83, 'd', 'o', 'g'}
	if !bytes.Equal(enc2, want2) {
		t.Fatalf("[cat,dog]: got %x, want %x", enc2, want2)
	}

	// empty string → [0x80]
	enc3, err := EncodeToBytes([]byte{})
	if err != nil {
		t.Fatalf("encode empty: %v", err)
	}
	if !bytes.Equal(enc3, []byte{0x80}) {
		t.Fatalf("empty: got %x", enc3)
	}

	// integer 0 → [0x80]
	enc4, err := EncodeToBytes(uint64(0))
	if err != nil {
		t.Fatalf("encode 0: %v", err)
	}
	if !bytes.Equal(enc4, []byte{0x80}) {
		t.Fatalf("uint64(0): got %x", enc4)
	}

	// integer 15 → [0x0f]
	enc5, err := EncodeToBytes(uint64(15))
	if err != nil {
		t.Fatalf("encode 15: %v", err)
	}
	if !bytes.Equal(enc5, []byte{0x0f}) {
		t.Fatalf("uint64(15): got %x", enc5)
	}

	// integer 1024 → [0x82, 0x04, 0x00]
	enc6, err := EncodeToBytes(uint64(1024))
	if err != nil {
		t.Fatalf("encode 1024: %v", err)
	}
	if !bytes.Equal(enc6, []byte{0x82, 0x04, 0x00}) {
		t.Fatalf("uint64(1024): got %x", enc6)
	}
}

// --- Error cases ---

func TestDecodeErrors(t *testing.T) {
	// extra bytes after value
	var n uint64
	if err := DecodeBytes([]byte{0x01, 0x02}, &n); err != ErrMoreThanOneValue {
		t.Fatalf("extra bytes: got %v, want ErrMoreThanOneValue", err)
	}

	// non-canonical integer (leading zero)
	if err := DecodeBytes([]byte{0x82, 0x00, 0x01}, &n); err != ErrCanonInt {
		t.Fatalf("leading zero int: got %v, want ErrCanonInt", err)
	}

	// expected list but got string
	var s simpleStruct
	if err := DecodeBytes([]byte{0x80}, &s); err != ErrExpectedList {
		t.Fatalf("expected list: got %v, want ErrExpectedList", err)
	}
}
