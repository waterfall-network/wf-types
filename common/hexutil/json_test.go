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
	"encoding/json"
	"math/big"
	"reflect"
	"testing"
)

// ---- Bytes ----

func TestBytesJSONRoundTrip(t *testing.T) {
	cases := []Bytes{
		{},
		{0x01},
		{0xde, 0xad, 0xbe, 0xef},
	}
	for _, b := range cases {
		data, err := json.Marshal(b)
		if err != nil {
			t.Fatalf("Marshal Bytes %x: %v", []byte(b), err)
		}
		var decoded Bytes
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal Bytes %x: %v", []byte(b), err)
		}
		if string(b) != string(decoded) {
			t.Errorf("Bytes round-trip: got %x, want %x", []byte(decoded), []byte(b))
		}
	}
}

func TestBytesUnmarshalErrors(t *testing.T) {
	var b Bytes
	if err := json.Unmarshal([]byte(`123`), &b); err == nil {
		t.Error("Bytes: expected error for non-string JSON")
	}
	if err := json.Unmarshal([]byte(`"deadbeef"`), &b); err == nil {
		t.Error("Bytes: expected error for missing 0x prefix")
	}
	if err := json.Unmarshal([]byte(`"0x1"`), &b); err == nil {
		t.Error("Bytes: expected error for odd-length hex")
	}
}

func TestBytesString(t *testing.T) {
	b := Bytes{0xab, 0xcd}
	if b.String() != "0xabcd" {
		t.Errorf("Bytes.String() = %q, want %q", b.String(), "0xabcd")
	}
}

func TestBytesTextRoundTrip(t *testing.T) {
	b := Bytes{0x01, 0x02, 0x03}
	text, err := b.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	var decoded Bytes
	if err := decoded.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}
	if string(b) != string(decoded) {
		t.Errorf("Bytes text round-trip: got %x, want %x", []byte(decoded), []byte(b))
	}
}

// ---- Big ----

func TestBigJSONRoundTrip(t *testing.T) {
	values := []int64{0, 1, 255, 1<<32 - 1}
	for _, v := range values {
		orig := (*Big)(big.NewInt(v))
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("Marshal Big %d: %v", v, err)
		}
		var decoded Big
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal Big %d: %v", v, err)
		}
		if orig.ToInt().Cmp(decoded.ToInt()) != 0 {
			t.Errorf("Big round-trip %d: got %v", v, decoded.ToInt())
		}
	}
}

func TestBigUnmarshalErrors(t *testing.T) {
	var b Big
	if err := json.Unmarshal([]byte(`123`), &b); err == nil {
		t.Error("Big: expected error for non-string JSON")
	}
	if err := json.Unmarshal([]byte(`"0x01"`), &b); err == nil {
		t.Error("Big: expected error for leading zero")
	}
	// oversized: > 256 bits
	oversize := `"0x` + "1" + "0000000000000000000000000000000000000000000000000000000000000000" + `"`
	if err := json.Unmarshal([]byte(oversize), &b); err == nil {
		t.Error("Big: expected error for > 256 bits")
	}
}

func TestBigString(t *testing.T) {
	b := (*Big)(big.NewInt(255))
	if b.String() != "0xff" {
		t.Errorf("Big.String() = %q, want %q", b.String(), "0xff")
	}
}

// ---- Uint64 ----

func TestUint64JSONRoundTrip(t *testing.T) {
	values := []uint64{0, 1, 255, 1 << 32}
	for _, v := range values {
		orig := Uint64(v)
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("Marshal Uint64 %d: %v", v, err)
		}
		var decoded Uint64
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal Uint64 %d: %v", v, err)
		}
		if orig != decoded {
			t.Errorf("Uint64 round-trip %d: got %d", v, decoded)
		}
	}
}

func TestUint64UnmarshalErrors(t *testing.T) {
	var u Uint64
	if err := json.Unmarshal([]byte(`123`), &u); err == nil {
		t.Error("Uint64: expected error for non-string JSON")
	}
	// overflow: more than 16 hex digits
	if err := json.Unmarshal([]byte(`"0xffffffffffffffff0"`), &u); err == nil {
		t.Error("Uint64: expected error for overflow")
	}
}

func TestUint64String(t *testing.T) {
	u := Uint64(255)
	if u.String() != "0xff" {
		t.Errorf("Uint64.String() = %q, want %q", u.String(), "0xff")
	}
}

// ---- Uint ----

func TestUintJSONRoundTrip(t *testing.T) {
	values := []uint{0, 1, 255, 65535}
	for _, v := range values {
		orig := Uint(v)
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("Marshal Uint %d: %v", v, err)
		}
		var decoded Uint
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal Uint %d: %v", v, err)
		}
		if orig != decoded {
			t.Errorf("Uint round-trip %d: got %d", v, decoded)
		}
	}
}

func TestUintString(t *testing.T) {
	u := Uint(16)
	if u.String() != "0x10" {
		t.Errorf("Uint.String() = %q, want %q", u.String(), "0x10")
	}
}

// ---- Uint8 ----

func TestUint8JSONRoundTrip(t *testing.T) {
	values := []uint8{0, 1, 127, 255}
	for _, v := range values {
		orig := Uint8(v)
		data, err := json.Marshal(orig)
		if err != nil {
			t.Fatalf("Marshal Uint8 %d: %v", v, err)
		}
		var decoded Uint8
		if err := json.Unmarshal(data, &decoded); err != nil {
			t.Fatalf("Unmarshal Uint8 %d: %v", v, err)
		}
		if orig != decoded {
			t.Errorf("Uint8 round-trip %d: got %d", v, decoded)
		}
	}
}

func TestUint8Overflow(t *testing.T) {
	var u Uint8
	if err := json.Unmarshal([]byte(`"0x100"`), &u); err == nil {
		t.Error("Uint8: expected error for value > 255")
	}
}

func TestUint8String(t *testing.T) {
	u := Uint8(0xff)
	if u.String() != "0xff" {
		t.Errorf("Uint8.String() = %q, want %q", u.String(), "0xff")
	}
}

// ---- UnmarshalFixedJSON / UnmarshalFixedText ----

func TestUnmarshalFixedJSON(t *testing.T) {
	typ := reflect.TypeOf(Bytes(nil))
	out := make([]byte, 4)
	if err := UnmarshalFixedJSON(typ, []byte(`"0xdeadbeef"`), out); err != nil {
		t.Fatalf("UnmarshalFixedJSON: %v", err)
	}
	if out[0] != 0xde || out[3] != 0xef {
		t.Errorf("UnmarshalFixedJSON: unexpected %x", out)
	}
	// wrong length
	if err := UnmarshalFixedJSON(typ, []byte(`"0xdeadbeef"`), make([]byte, 2)); err == nil {
		t.Error("UnmarshalFixedJSON: expected error for wrong length")
	}
	// non-string JSON
	if err := UnmarshalFixedJSON(typ, []byte(`123`), out); err == nil {
		t.Error("UnmarshalFixedJSON: expected error for non-string")
	}
}

func TestUnmarshalFixedText(t *testing.T) {
	out := make([]byte, 4)
	if err := UnmarshalFixedText("test", []byte("0xdeadbeef"), out); err != nil {
		t.Fatalf("UnmarshalFixedText: %v", err)
	}
	if out[0] != 0xde || out[3] != 0xef {
		t.Errorf("UnmarshalFixedText: unexpected %x", out)
	}
	// wrong length
	if err := UnmarshalFixedText("test", []byte("0xdeadbeef"), make([]byte, 2)); err == nil {
		t.Error("UnmarshalFixedText: expected error for wrong length")
	}
}

func TestUnmarshalFixedUnprefixedText(t *testing.T) {
	out := make([]byte, 4)
	// without prefix
	if err := UnmarshalFixedUnprefixedText("test", []byte("deadbeef"), out); err != nil {
		t.Fatalf("UnmarshalFixedUnprefixedText (no prefix): %v", err)
	}
	if out[0] != 0xde || out[3] != 0xef {
		t.Errorf("UnmarshalFixedUnprefixedText: unexpected %x", out)
	}
	// with prefix
	if err := UnmarshalFixedUnprefixedText("test", []byte("0xdeadbeef"), out); err != nil {
		t.Fatalf("UnmarshalFixedUnprefixedText (with prefix): %v", err)
	}
}
