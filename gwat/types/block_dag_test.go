// Copyright 2024 Blue Wave Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package types

import (
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

func TestBlockDAGBytesRoundTrip(t *testing.T) {
	b := &BlockDAG{
		Hash:     common.HexToHash("0x01"),
		Height:   100,
		Slot:     200,
		CpHash:   common.HexToHash("0x02"),
		CpHeight: 99,
		OrderedAncestorsHashes: common.HashArray{
			common.HexToHash("0xaa"),
			common.HexToHash("0xbb"),
		},
	}

	data := b.ToBytes()
	var b2 BlockDAG
	b2.SetBytes(data)

	if b.Hash != b2.Hash || b.Height != b2.Height || b.Slot != b2.Slot {
		t.Fatalf("BlockDAG bytes round-trip mismatch")
	}
	if !b.OrderedAncestorsHashes.IsEqualTo(b2.OrderedAncestorsHashes) {
		t.Fatalf("OrderedAncestorsHashes mismatch")
	}
}

func TestTipsOperations(t *testing.T) {
	tips := make(Tips)

	b1 := &BlockDAG{Hash: common.HexToHash("0x01"), Height: 10, Slot: 5}
	b2 := &BlockDAG{Hash: common.HexToHash("0x02"), Height: 20, Slot: 6}

	tips.Add(b1)
	tips.Add(b2)

	if len(tips) != 2 {
		t.Fatalf("expected 2 tips, got %d", len(tips))
	}

	// highest height first
	finDag := tips.GetFinalizingDag()
	if finDag == nil || finDag.Hash != b2.Hash {
		t.Fatalf("expected b2 as finalizing dag")
	}

	tips.Remove(b2.Hash)
	if len(tips) != 1 {
		t.Fatalf("expected 1 tip after remove, got %d", len(tips))
	}
}

func TestTipsNilIgnored(t *testing.T) {
	tips := make(Tips)
	tips.Add(nil) // should not panic or add
	if len(tips) != 0 {
		t.Fatal("nil blockDag should be ignored")
	}
}

func TestSlotInfo(t *testing.T) {
	si := &SlotInfo{
		GenesisTime:    1000,
		SecondsPerSlot: 12,
		SlotsPerEpoch:  32,
	}

	epoch := si.SlotToEpoch(64)
	if epoch != 2 {
		t.Fatalf("expected epoch 2, got %d", epoch)
	}

	if !si.IsEpochStart(0) {
		t.Fatal("slot 0 should be epoch start")
	}
	if !si.IsEpochStart(32) {
		t.Fatal("slot 32 should be epoch start")
	}
	if si.IsEpochStart(1) {
		t.Fatal("slot 1 should not be epoch start")
	}

	start, err := si.SlotOfEpochStart(3)
	if err != nil {
		t.Fatal(err)
	}
	if start != 96 {
		t.Fatalf("expected slot 96, got %d", start)
	}

	end, err := si.SlotOfEpochEnd(3)
	if err != nil {
		t.Fatal(err)
	}
	if end != 127 {
		t.Fatalf("expected slot 127, got %d", end)
	}

	si2 := si.Copy()
	si2.SlotsPerEpoch = 64
	if si.SlotsPerEpoch == si2.SlotsPerEpoch {
		t.Fatal("Copy should be independent")
	}
}
