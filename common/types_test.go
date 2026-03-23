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

package common

import (
	"encoding/json"
	"math/big"
	"strings"
	"testing"
)

// ---- Hash ----

func TestBytesToHash(t *testing.T) {
	b := []byte{1, 2, 3}
	h := BytesToHash(b)
	if h[HashLength-1] != 3 || h[HashLength-2] != 2 || h[HashLength-3] != 1 {
		t.Errorf("BytesToHash: unexpected value %x", h)
	}
}

func TestHashSetBytes_Truncation(t *testing.T) {
	long := make([]byte, HashLength+5)
	long[len(long)-1] = 0xff
	var h Hash
	h.SetBytes(long)
	if h[HashLength-1] != 0xff {
		t.Error("SetBytes: last byte should be 0xff")
	}
}

func TestHashHex(t *testing.T) {
	h := BytesToHash([]byte{0xab, 0xcd})
	hex := h.Hex()
	if !strings.HasPrefix(hex, "0x") {
		t.Errorf("Hash.Hex() missing 0x prefix: %s", hex)
	}
	if !strings.HasSuffix(hex, "abcd") {
		t.Errorf("Hash.Hex() wrong value: %s", hex)
	}
}

func TestHashJSON(t *testing.T) {
	original := BytesToHash([]byte{0xde, 0xad, 0xbe, 0xef})
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var decoded Hash
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if original != decoded {
		t.Errorf("JSON round-trip: got %v, want %v", decoded, original)
	}
}

func TestHashBig(t *testing.T) {
	h := BigToHash(big.NewInt(255))
	if h.Big().Int64() != 255 {
		t.Errorf("HashBig round-trip failed: %v", h.Big())
	}
}

func TestHashPtr(t *testing.T) {
	h := BytesToHash([]byte{1})
	p := h.Ptr()
	if p == nil || *p != h {
		t.Error("Hash.Ptr() returned wrong pointer")
	}
}

func TestHashScan(t *testing.T) {
	src := make([]byte, HashLength)
	src[0] = 0xaa
	var h Hash
	if err := h.Scan(src); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if h[0] != 0xaa {
		t.Error("Scan: wrong value")
	}
	// wrong length
	if err := h.Scan([]byte{1, 2}); err == nil {
		t.Error("Scan: expected error for wrong length")
	}
	// wrong type
	if err := h.Scan("string"); err == nil {
		t.Error("Scan: expected error for wrong type")
	}
}

func TestHashValue(t *testing.T) {
	h := BytesToHash([]byte{0xff})
	v, err := h.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	b, ok := v.([]byte)
	if !ok || len(b) != HashLength {
		t.Error("Value: unexpected result")
	}
}

// ---- HashArray ----

func TestHashArrayBytesRoundTrip(t *testing.T) {
	ha := HashArray{BytesToHash([]byte{1}), BytesToHash([]byte{2})}
	b := ha.ToBytes()
	decoded := HashArrayFromBytes(b)
	if !ha.IsEqualTo(decoded) {
		t.Errorf("HashArray bytes round-trip failed")
	}
}

func TestHashArrayIntersection(t *testing.T) {
	a := HashArray{BytesToHash([]byte{1}), BytesToHash([]byte{2}), BytesToHash([]byte{3})}
	b := HashArray{BytesToHash([]byte{2}), BytesToHash([]byte{4})}
	got := a.Intersection(b)
	if len(got) != 1 || got[0] != BytesToHash([]byte{2}) {
		t.Errorf("Intersection: %v", got)
	}
}

func TestHashArrayDifference(t *testing.T) {
	a := HashArray{BytesToHash([]byte{1}), BytesToHash([]byte{2}), BytesToHash([]byte{3})}
	b := HashArray{BytesToHash([]byte{2})}
	got := a.Difference(b)
	if len(got) != 2 {
		t.Errorf("Difference: len=%d want 2", len(got))
	}
}

func TestHashArrayUniq(t *testing.T) {
	h := BytesToHash([]byte{1})
	ha := HashArray{h, h, BytesToHash([]byte{2})}
	got := ha.Uniq()
	if len(got) != 2 {
		t.Errorf("Uniq: len=%d want 2", len(got))
	}
}

func TestHashArrayDeduplicate(t *testing.T) {
	h := BytesToHash([]byte{1})
	ha := HashArray{h, h, BytesToHash([]byte{2})}
	ha.Deduplicate()
	if len(ha) != 2 {
		t.Errorf("Deduplicate: len=%d want 2", len(ha))
	}
}

func TestHashArrayIsUniq(t *testing.T) {
	h1 := BytesToHash([]byte{1})
	h2 := BytesToHash([]byte{2})
	if !(HashArray{h1, h2}).IsUniq() {
		t.Error("IsUniq: should be true")
	}
	if (HashArray{h1, h1}).IsUniq() {
		t.Error("IsUniq: should be false for duplicates")
	}
}

func TestHashArrayHasIndexOf(t *testing.T) {
	h1 := BytesToHash([]byte{1})
	h2 := BytesToHash([]byte{2})
	ha := HashArray{h1, h2}
	if !ha.Has(h1) {
		t.Error("Has: should find h1")
	}
	if ha.Has(BytesToHash([]byte{99})) {
		t.Error("Has: should not find unknown hash")
	}
	if ha.IndexOf(h2) != 1 {
		t.Errorf("IndexOf: want 1, got %d", ha.IndexOf(h2))
	}
	if ha.IndexOf(BytesToHash([]byte{99})) != -1 {
		t.Error("IndexOf: should return -1 for missing hash")
	}
}

func TestHashArrayReverse(t *testing.T) {
	h1 := BytesToHash([]byte{1})
	h2 := BytesToHash([]byte{2})
	ha := HashArray{h1, h2}
	rev := ha.Reverse()
	if rev[0] != h2 || rev[1] != h1 {
		t.Errorf("Reverse: %v", rev)
	}
	// original unchanged
	if ha[0] != h1 {
		t.Error("Reverse modified original")
	}
}

func TestHashArrayKey(t *testing.T) {
	h1 := BytesToHash([]byte{1})
	h2 := BytesToHash([]byte{2})
	// Key is order-independent (sorts internally)
	k1 := HashArray{h1, h2}.Key()
	k2 := HashArray{h2, h1}.Key()
	if k1 != k2 {
		t.Error("Key: should be order-independent")
	}
	// Key of different sets differs
	k3 := HashArray{h1}.Key()
	if k1 == k3 {
		t.Error("Key: different sets should produce different keys")
	}
}

func TestHashArrayConcat(t *testing.T) {
	h1 := BytesToHash([]byte{1})
	h2 := BytesToHash([]byte{2})
	got := HashArray{h1}.Concat(HashArray{h2})
	if len(got) != 2 || got[1] != h2 {
		t.Errorf("Concat: %v", got)
	}
}

func TestHashArrayCopy(t *testing.T) {
	h := BytesToHash([]byte{1})
	ha := HashArray{h}
	c := ha.Copy()
	c[0] = BytesToHash([]byte{99})
	if ha[0] != h {
		t.Error("Copy: modifying copy affected original")
	}
}

func TestHashArraySequenceIntersection(t *testing.T) {
	h1 := BytesToHash([]byte{1})
	h2 := BytesToHash([]byte{2})
	h3 := BytesToHash([]byte{3})
	ha := HashArray{h1, h2, h3}
	// contiguous sequence starting from h2
	got := ha.SequenceIntersection(HashArray{h2, h3})
	if len(got) != 2 || got[0] != h2 || got[1] != h3 {
		t.Errorf("SequenceIntersection: %v", got)
	}
}

// ---- Address ----

func TestBytesToAddress(t *testing.T) {
	b := []byte{0xaa, 0xbb}
	a := BytesToAddress(b)
	if a[AddressLength-1] != 0xbb || a[AddressLength-2] != 0xaa {
		t.Errorf("BytesToAddress: unexpected value %x", a)
	}
}

func TestAddressHex(t *testing.T) {
	a := BytesToAddress([]byte{0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef,
		0x01, 0x23, 0x45, 0x67})
	hex := a.Hex()
	if !strings.HasPrefix(hex, "0x") {
		t.Errorf("Address.Hex() missing 0x prefix: %s", hex)
	}
}

func TestAddressJSON(t *testing.T) {
	original := BytesToAddress([]byte{
		0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a,
		0x0b, 0x0c, 0x0d, 0x0e, 0x0f, 0x10, 0x11, 0x12, 0x13, 0x14,
	})
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	var decoded Address
	if err := json.Unmarshal(data, &decoded); err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if original != decoded {
		t.Errorf("JSON round-trip: got %v, want %v", decoded, original)
	}
}

func TestIsHexAddress(t *testing.T) {
	valid := "0x" + strings.Repeat("a", 40)
	if !IsHexAddress(valid) {
		t.Errorf("IsHexAddress(%s) should be true", valid)
	}
	if IsHexAddress("0xabc") {
		t.Error("IsHexAddress: short string should be false")
	}
	if IsHexAddress("notahex") {
		t.Error("IsHexAddress: invalid chars should be false")
	}
}

func TestContains(t *testing.T) {
	a1 := BytesToAddress([]byte{1})
	a2 := BytesToAddress([]byte{2})
	addrs := []Address{a1, a2}
	found, idx := Contains(addrs, a1)
	if !found || idx != 0 {
		t.Errorf("Contains: found=%v idx=%d", found, idx)
	}
	found, idx = Contains(addrs, BytesToAddress([]byte{99}))
	if found || idx != -1 {
		t.Errorf("Contains missing: found=%v idx=%d", found, idx)
	}
}

func TestSortAddresses(t *testing.T) {
	a1 := HexToAddress("0x" + strings.Repeat("01", 20))
	a2 := HexToAddress("0x" + strings.Repeat("02", 20))
	sorted := SortAddresses([]Address{a2, a1})
	if sorted[0] != a1 {
		t.Errorf("SortAddresses: first element should be smaller")
	}
}

func TestAddressScan(t *testing.T) {
	src := make([]byte, AddressLength)
	src[0] = 0xbb
	var a Address
	if err := a.Scan(src); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if a[0] != 0xbb {
		t.Error("Scan: wrong value")
	}
	if err := a.Scan([]byte{1}); err == nil {
		t.Error("Scan: expected error for wrong length")
	}
}

func TestMixedcaseAddress(t *testing.T) {
	addr := BytesToAddress([]byte{
		0xde, 0xad, 0xbe, 0xef,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x01,
	})
	ma := NewMixedcaseAddress(addr)
	if ma.Address() != addr {
		t.Error("MixedcaseAddress.Address() mismatch")
	}
	if ma.Original() != addr.Hex() {
		t.Error("MixedcaseAddress.Original() mismatch")
	}
	if !ma.ValidChecksum() {
		t.Error("MixedcaseAddress: checksum should be valid when created from Address")
	}

	// From string: valid checksum address
	hexStr := addr.Hex()
	maf, err := NewMixedcaseAddressFromString(hexStr)
	if err != nil {
		t.Fatalf("NewMixedcaseAddressFromString: %v", err)
	}
	if maf.Address() != addr {
		t.Error("NewMixedcaseAddressFromString: address mismatch")
	}
	// invalid
	if _, err := NewMixedcaseAddressFromString("not-an-address"); err == nil {
		t.Error("NewMixedcaseAddressFromString: expected error for invalid input")
	}
}

func TestUnprefixedHash(t *testing.T) {
	h := BytesToHash([]byte{0xab, 0xcd})
	uh := UnprefixedHash(h)
	text, err := uh.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	var decoded UnprefixedHash
	if err := decoded.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}
	if uh != decoded {
		t.Errorf("UnprefixedHash round-trip failed")
	}
}

func TestUnprefixedAddress(t *testing.T) {
	a := BytesToAddress([]byte{0x11, 0x22})
	ua := UnprefixedAddress(a)
	text, err := ua.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText: %v", err)
	}
	var decoded UnprefixedAddress
	if err := decoded.UnmarshalText(text); err != nil {
		t.Fatalf("UnmarshalText: %v", err)
	}
	if ua != decoded {
		t.Errorf("UnprefixedAddress round-trip failed")
	}
}
