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

package crypto

import (
	"encoding/hex"
	"testing"
)

// ---- Keccak256 ----

func TestKeccak256Empty(t *testing.T) {
	got := Keccak256()
	// keccak256("") = c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470
	want := "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470"
	if hex.EncodeToString(got) != want {
		t.Errorf("Keccak256() = %x, want %s", got, want)
	}
}

func TestKeccak256SingleChunk(t *testing.T) {
	got := Keccak256([]byte("hello"))
	// keccak256("hello") = 1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8
	want := "1c8aff950685c2ed4bc3174f3472287b56d9517b9c948127319a09a7a36deac8"
	if hex.EncodeToString(got) != want {
		t.Errorf("Keccak256(hello) = %x, want %s", got, want)
	}
}

func TestKeccak256MultiChunk(t *testing.T) {
	// keccak256("hello" + "world") == keccak256("helloworld")
	got := Keccak256([]byte("hello"), []byte("world"))
	want := Keccak256([]byte("helloworld"))
	if hex.EncodeToString(got) != hex.EncodeToString(want) {
		t.Errorf("Keccak256 multi-chunk: %x != %x", got, want)
	}
}

func TestKeccak256Length(t *testing.T) {
	if len(Keccak256([]byte("test"))) != 32 {
		t.Error("Keccak256 should return 32 bytes")
	}
}

// ---- Keccak256Hash ----

func TestKeccak256HashEmpty(t *testing.T) {
	got := Keccak256Hash()
	want := Keccak256()
	for i, b := range want {
		if got[i] != b {
			t.Errorf("Keccak256Hash()[%d] = %x, want %x", i, got[i], b)
		}
	}
}

func TestKeccak256HashMultiChunk(t *testing.T) {
	a := Keccak256Hash([]byte("foo"), []byte("bar"))
	b := Keccak256Hash([]byte("foobar"))
	if a != b {
		t.Errorf("Keccak256Hash multi-chunk: %x != %x", a, b)
	}
}

// ---- CreateAddress ----

// Test vectors taken from Ethereum Yellow Paper and confirmed against go-ethereum.
func TestCreateAddress(t *testing.T) {
	tests := []struct {
		addrHex string
		nonce   uint64
		wantHex string
	}{
		// nonce = 0
		{
			"6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0",
			0,
			"cd234a471b72ba2f1ccf0a70fcaba648a5eecd8d",
		},
		// nonce = 1
		{
			"6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0",
			1,
			"343c43a37d37dff08ae8c4a11544c718abb4fcf8",
		},
		// nonce = 2
		{
			"6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0",
			2,
			"f778b86fa74e846c4f0a1fbd1335fe81c00a0c91",
		},
		// nonce = 3
		{
			"6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0",
			3,
			"fffd933a0bc612844eaf0c6fe3e5b8e9b6c1d19c",
		},
	}
	for _, tt := range tests {
		addrBytes, err := hex.DecodeString(tt.addrHex)
		if err != nil {
			t.Fatalf("bad test addr: %v", err)
		}
		var addr [20]byte
		copy(addr[:], addrBytes)

		got := CreateAddress(addr, tt.nonce)
		want := tt.wantHex
		if hex.EncodeToString(got[:]) != want {
			t.Errorf("CreateAddress(%s, %d) = %x, want %s",
				tt.addrHex, tt.nonce, got, want)
		}
	}
}

func TestCreateAddressHighNonce(t *testing.T) {
	// nonce = 0x100 (256) — multi-byte nonce RLP encoding
	addrBytes, _ := hex.DecodeString("6ac7ea33f8831ea9dcc53393aaa88b25a785dbf0")
	var addr [20]byte
	copy(addr[:], addrBytes)

	got := CreateAddress(addr, 0x100)
	want := "3837c1ae70354f670550c746580199ac6a73cb0a"
	if hex.EncodeToString(got[:]) != want {
		t.Errorf("CreateAddress nonce=256: got %x, want %s", got, want)
	}
}

// ---- rlpEncodeUint64 (internal, tested via CreateAddress) ----

func TestRlpEncodeUint64(t *testing.T) {
	tests := []struct {
		n    uint64
		want []byte
	}{
		{0, []byte{0x80}},
		{1, []byte{0x01}},
		{0x7f, []byte{0x7f}},
		{0x80, []byte{0x81, 0x80}},
		{0xff, []byte{0x81, 0xff}},
		{0x100, []byte{0x82, 0x01, 0x00}},
		{0xffff, []byte{0x82, 0xff, 0xff}},
	}
	for _, tt := range tests {
		got := rlpEncodeUint64(tt.n)
		if len(got) != len(tt.want) {
			t.Errorf("rlpEncodeUint64(%d): got %x, want %x", tt.n, got, tt.want)
			continue
		}
		for i := range got {
			if got[i] != tt.want[i] {
				t.Errorf("rlpEncodeUint64(%d)[%d]: got %x, want %x", tt.n, i, got[i], tt.want[i])
			}
		}
	}
}
