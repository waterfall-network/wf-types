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

package hexutil

import (
	"math/big"
	"testing"
)

// ---- Decode / Encode ----

func TestEncode(t *testing.T) {
	tests := []struct {
		input []byte
		want  string
	}{
		{[]byte{}, "0x"},
		{[]byte{0x01, 0x02}, "0x0102"},
		{[]byte{0xde, 0xad, 0xbe, 0xef}, "0xdeadbeef"},
	}
	for _, tt := range tests {
		got := Encode(tt.input)
		if got != tt.want {
			t.Errorf("Encode(%x) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDecode(t *testing.T) {
	tests := []struct {
		input   string
		want    []byte
		wantErr error
	}{
		{"0x", []byte{}, nil},
		{"0x0102", []byte{0x01, 0x02}, nil},
		{"0xdeadbeef", []byte{0xde, 0xad, 0xbe, 0xef}, nil},
		{"", nil, ErrEmptyString},
		{"deadbeef", nil, ErrMissingPrefix},
		{"0x1", nil, ErrOddLength},
		{"0xgg", nil, ErrSyntax},
	}
	for _, tt := range tests {
		got, err := Decode(tt.input)
		if tt.wantErr != nil {
			if err != tt.wantErr {
				t.Errorf("Decode(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("Decode(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if len(got) != len(tt.want) {
			t.Errorf("Decode(%q) len=%d want %d", tt.input, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("Decode(%q)[%d] = %x, want %x", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestMustDecode(t *testing.T) {
	got := MustDecode("0xabcd")
	if len(got) != 2 || got[0] != 0xab {
		t.Errorf("MustDecode unexpected: %x", got)
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDecode: expected panic for invalid input")
		}
	}()
	MustDecode("invalid")
}

// ---- DecodeUint64 / EncodeUint64 ----

func TestEncodeUint64(t *testing.T) {
	tests := []struct {
		input uint64
		want  string
	}{
		{0, "0x0"},
		{1, "0x1"},
		{255, "0xff"},
		{256, "0x100"},
	}
	for _, tt := range tests {
		got := EncodeUint64(tt.input)
		if got != tt.want {
			t.Errorf("EncodeUint64(%d) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDecodeUint64(t *testing.T) {
	tests := []struct {
		input   string
		want    uint64
		wantErr error
	}{
		{"0x0", 0, nil},
		{"0xff", 255, nil},
		{"0x100", 256, nil},
		{"", 0, ErrEmptyString},
		{"123", 0, ErrMissingPrefix},
		{"0x", 0, ErrEmptyNumber},
		{"0x01", 0, ErrLeadingZero},
		{"0xffffffffffffffff0", 0, ErrUint64Range},
	}
	for _, tt := range tests {
		got, err := DecodeUint64(tt.input)
		if tt.wantErr != nil {
			if err != tt.wantErr {
				t.Errorf("DecodeUint64(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("DecodeUint64(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("DecodeUint64(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestMustDecodeUint64(t *testing.T) {
	if MustDecodeUint64("0xff") != 255 {
		t.Error("MustDecodeUint64 failed")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDecodeUint64: expected panic for invalid input")
		}
	}()
	MustDecodeUint64("invalid")
}

// ---- DecodeBig / EncodeBig ----

func TestEncodeBig(t *testing.T) {
	tests := []struct {
		input *big.Int
		want  string
	}{
		{big.NewInt(0), "0x0"},
		{big.NewInt(255), "0xff"},
		{big.NewInt(256), "0x100"},
	}
	for _, tt := range tests {
		got := EncodeBig(tt.input)
		if got != tt.want {
			t.Errorf("EncodeBig(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestDecodeBig(t *testing.T) {
	tests := []struct {
		input   string
		want    int64
		wantErr error
	}{
		{"0x0", 0, nil},
		{"0xff", 255, nil},
		{"0x100", 256, nil},
		{"", 0, ErrEmptyString},
		{"256", 0, ErrMissingPrefix},
		{"0x", 0, ErrEmptyNumber},
	}
	for _, tt := range tests {
		got, err := DecodeBig(tt.input)
		if tt.wantErr != nil {
			if err != tt.wantErr {
				t.Errorf("DecodeBig(%q) error = %v, want %v", tt.input, err, tt.wantErr)
			}
			continue
		}
		if err != nil {
			t.Errorf("DecodeBig(%q) unexpected error: %v", tt.input, err)
			continue
		}
		if got.Int64() != tt.want {
			t.Errorf("DecodeBig(%q) = %d, want %d", tt.input, got.Int64(), tt.want)
		}
	}
}

func TestDecodeBig256Range(t *testing.T) {
	// 65 hex digits after 0x = > 256 bits
	_, err := DecodeBig("0x" + "1" + "0000000000000000000000000000000000000000000000000000000000000000")
	if err != ErrBig256Range {
		t.Errorf("DecodeBig oversized: error = %v, want ErrBig256Range", err)
	}
}

func TestMustDecodeBig(t *testing.T) {
	got := MustDecodeBig("0xff")
	if got.Int64() != 255 {
		t.Error("MustDecodeBig failed")
	}
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustDecodeBig: expected panic for invalid input")
		}
	}()
	MustDecodeBig("invalid")
}

func TestUint64RoundTrip(t *testing.T) {
	values := []uint64{0, 1, 127, 255, 65535, 1<<32 - 1, 1<<64 - 1}
	for _, v := range values {
		enc := EncodeUint64(v)
		dec, err := DecodeUint64(enc)
		if err != nil {
			t.Errorf("DecodeUint64(EncodeUint64(%d)): %v", v, err)
		}
		if dec != v {
			t.Errorf("Uint64 round-trip %d → %q → %d", v, enc, dec)
		}
	}
}

func TestBigRoundTrip(t *testing.T) {
	values := []*big.Int{
		big.NewInt(0),
		big.NewInt(1),
		big.NewInt(255),
		new(big.Int).Lsh(big.NewInt(1), 128),
	}
	for _, v := range values {
		enc := EncodeBig(v)
		dec, err := DecodeBig(enc)
		if err != nil {
			t.Errorf("DecodeBig(EncodeBig(%v)): %v", v, err)
		}
		if v.Cmp(dec) != 0 {
			t.Errorf("Big round-trip %v → %q → %v", v, enc, dec)
		}
	}
}
