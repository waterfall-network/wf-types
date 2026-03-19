// Copyright 2024 Blue Wave Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package common

import (
	"encoding/json"
	"testing"
)

func TestHashRoundTrip(t *testing.T) {
	h := HexToHash("0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef")
	if h.Hex() != "0xdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeefdeadbeef" {
		t.Fatalf("unexpected hex: %s", h.Hex())
	}

	b, err := json.Marshal(h)
	if err != nil {
		t.Fatal(err)
	}
	var h2 Hash
	if err := json.Unmarshal(b, &h2); err != nil {
		t.Fatal(err)
	}
	if h != h2 {
		t.Fatalf("marshal/unmarshal mismatch: %s != %s", h, h2)
	}
}

func TestAddressChecksumHex(t *testing.T) {
	// EIP-55 checksum test
	a := HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
	expected := "0x5aAeb6053F3E94C9b9A09f33669435E7Ef1BeAed"
	if a.Hex() != expected {
		t.Fatalf("expected %s, got %s", expected, a.Hex())
	}
}

func TestHashArrayKey(t *testing.T) {
	ha1 := HashArray{HexToHash("0x01"), HexToHash("0x02")}
	ha2 := HashArray{HexToHash("0x02"), HexToHash("0x01")}
	// Key is based on unique+sorted bytes, so order should not matter
	if ha1.Key() != ha2.Key() {
		t.Fatalf("expected same key for same set: %s != %s", ha1.Key(), ha2.Key())
	}
}

func TestHashArrayOperations(t *testing.T) {
	ha := HashArray{HexToHash("0x01"), HexToHash("0x02"), HexToHash("0x01")}
	uniq := ha.Uniq()
	if len(uniq) != 2 {
		t.Fatalf("expected 2 unique, got %d", len(uniq))
	}

	diff := HashArray{HexToHash("0x01"), HexToHash("0x02"), HexToHash("0x03")}.Difference(HashArray{HexToHash("0x01")})
	if len(diff) != 2 {
		t.Fatalf("expected 2 in diff, got %d", len(diff))
	}

	inter := HashArray{HexToHash("0x01"), HexToHash("0x02")}.Intersection(HashArray{HexToHash("0x02"), HexToHash("0x03")})
	if len(inter) != 1 || inter[0] != HexToHash("0x02") {
		t.Fatalf("unexpected intersection: %v", inter)
	}
}

func TestHashArrayBytes(t *testing.T) {
	ha := HashArray{HexToHash("0x01"), HexToHash("0x02")}
	b := ha.ToBytes()
	ha2 := HashArrayFromBytes(b)
	if !ha.IsEqualTo(ha2) {
		t.Fatalf("round-trip failed: %v != %v", ha, ha2)
	}
}
