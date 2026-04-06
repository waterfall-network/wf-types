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

package common

import (
	"testing"
)

func TestFromHex(t *testing.T) {
	tests := []struct {
		input string
		want  []byte
	}{
		{"0x0102", []byte{1, 2}},
		{"0X0102", []byte{1, 2}},
		{"0102", []byte{1, 2}},
		{"0x1", []byte{1}}, // odd length: padded to "01"
		{"0x", []byte{}},
		{"", []byte{}},
		{"0xdeadbeef", []byte{0xde, 0xad, 0xbe, 0xef}},
	}
	for _, tt := range tests {
		got := FromHex(tt.input)
		if len(got) != len(tt.want) {
			t.Errorf("FromHex(%q): len=%d want %d", tt.input, len(got), len(tt.want))
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("FromHex(%q)[%d] = %x, want %x", tt.input, i, got[i], tt.want[i])
			}
		}
	}
}

func TestCopyBytes(t *testing.T) {
	src := []byte{1, 2, 3}
	dst := CopyBytes(src)
	if len(dst) != len(src) {
		t.Fatalf("CopyBytes: len=%d want %d", len(dst), len(src))
	}
	dst[0] = 99
	if src[0] == 99 {
		t.Error("CopyBytes: modifying copy affected original")
	}

	if CopyBytes(nil) != nil {
		t.Error("CopyBytes(nil) should return nil")
	}
}

func TestBytes2Hex(t *testing.T) {
	got := Bytes2Hex([]byte{0xde, 0xad, 0xbe, 0xef})
	if got != "deadbeef" {
		t.Errorf("Bytes2Hex = %q, want %q", got, "deadbeef")
	}
}

func TestHex2Bytes(t *testing.T) {
	got := Hex2Bytes("deadbeef")
	want := []byte{0xde, 0xad, 0xbe, 0xef}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Hex2Bytes[%d] = %x, want %x", i, got[i], want[i])
		}
	}
}

func TestHex2BytesFixed(t *testing.T) {
	// exact length
	got := Hex2BytesFixed("deadbeef", 4)
	if got[0] != 0xde || got[3] != 0xef {
		t.Errorf("Hex2BytesFixed exact: %x", got)
	}
	// shorter than flen: left-padded with zeros
	got = Hex2BytesFixed("ff", 4)
	if got[0] != 0 || got[1] != 0 || got[2] != 0 || got[3] != 0xff {
		t.Errorf("Hex2BytesFixed padded: %x", got)
	}
	// longer than flen: cropped from left
	got = Hex2BytesFixed("aabbccdd", 2)
	if got[0] != 0xcc || got[1] != 0xdd {
		t.Errorf("Hex2BytesFixed cropped: %x", got)
	}
}

func TestRightPadBytes(t *testing.T) {
	got := RightPadBytes([]byte{1, 2}, 5)
	if len(got) != 5 || got[0] != 1 || got[1] != 2 || got[2] != 0 {
		t.Errorf("RightPadBytes = %v", got)
	}
	// already long enough
	src := []byte{1, 2, 3}
	got = RightPadBytes(src, 2)
	if &got[0] != &src[0] {
		t.Error("RightPadBytes should return original slice when l <= len")
	}
}

func TestLeftPadBytes(t *testing.T) {
	got := LeftPadBytes([]byte{1, 2}, 5)
	if len(got) != 5 || got[0] != 0 || got[3] != 1 || got[4] != 2 {
		t.Errorf("LeftPadBytes = %v", got)
	}
	// already long enough
	src := []byte{1, 2, 3}
	got = LeftPadBytes(src, 2)
	if &got[0] != &src[0] {
		t.Error("LeftPadBytes should return original slice when l <= len")
	}
}

func TestTrimLeftZeroes(t *testing.T) {
	got := TrimLeftZeroes([]byte{0, 0, 1, 2})
	if len(got) != 2 || got[0] != 1 {
		t.Errorf("TrimLeftZeroes = %v", got)
	}
	got = TrimLeftZeroes([]byte{0, 0, 0})
	if len(got) != 0 {
		t.Errorf("TrimLeftZeroes all-zero = %v", got)
	}
}

func TestTrimRightZeroes(t *testing.T) {
	got := TrimRightZeroes([]byte{1, 2, 0, 0})
	if len(got) != 2 || got[1] != 2 {
		t.Errorf("TrimRightZeroes = %v", got)
	}
	got = TrimRightZeroes([]byte{0, 0, 0})
	if len(got) != 0 {
		t.Errorf("TrimRightZeroes all-zero = %v", got)
	}
}

func TestUint64Bytes(t *testing.T) {
	values := []uint64{0, 1, 255, 1<<32 - 1, 1<<64 - 1}
	for _, v := range values {
		b := Uint64ToBytes(v)
		if len(b) != 8 {
			t.Errorf("Uint64ToBytes(%d): len=%d", v, len(b))
		}
		got := BytesToUint64(b)
		if got != v {
			t.Errorf("BytesToUint64(Uint64ToBytes(%d)) = %d", v, got)
		}
	}
}
