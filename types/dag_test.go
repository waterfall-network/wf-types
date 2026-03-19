// Copyright 2024 Blue Wave Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0

package types

import (
	"encoding/json"
	"math/big"
	"testing"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

func TestCheckpointRoundTrip(t *testing.T) {
	cp := &Checkpoint{
		Epoch:    10,
		FinEpoch: 8,
		Root:     common.HexToHash("0xaaaa"),
		Spine:    common.HexToHash("0xbbbb"),
	}

	// bytes round-trip
	b := cp.Bytes()
	cp2, err := BytesToCheckpoint(b)
	if err != nil {
		t.Fatal(err)
	}
	if cp.Epoch != cp2.Epoch || cp.FinEpoch != cp2.FinEpoch || cp.Root != cp2.Root || cp.Spine != cp2.Spine {
		t.Fatalf("bytes round-trip mismatch")
	}

	// JSON round-trip
	data, err := json.Marshal(cp)
	if err != nil {
		t.Fatal(err)
	}
	var cp3 Checkpoint
	if err := json.Unmarshal(data, &cp3); err != nil {
		t.Fatal(err)
	}
	if cp.Epoch != cp3.Epoch || cp.Root != cp3.Root {
		t.Fatalf("JSON round-trip mismatch")
	}
}

func TestCheckpointCopy(t *testing.T) {
	cp := &Checkpoint{Epoch: 5, FinEpoch: 3, Root: common.HexToHash("0x1234"), Spine: common.HexToHash("0x5678")}
	cp2 := cp.Copy()
	cp2.Epoch = 99
	if cp.Epoch == cp2.Epoch {
		t.Fatal("Copy should be independent")
	}
}

func TestValidatorSyncJSONRoundTrip(t *testing.T) {
	vs := &ValidatorSync{
		InitTxHash:      common.HexToHash("0xabcd"),
		OpType:          Activate,
		ProcEpoch:       5,
		Index:           42,
		Creator:         common.HexToAddress("0x1234567890123456789012345678901234567890"),
		Amount:          big.NewInt(1000),
		Balance:         big.NewInt(2000),
		ActivationEpoch: 3,
		ExitEpoch:       100,
	}

	data, err := json.Marshal(vs)
	if err != nil {
		t.Fatal(err)
	}

	var vs2 ValidatorSync
	if err := json.Unmarshal(data, &vs2); err != nil {
		t.Fatal(err)
	}

	if vs.ProcEpoch != vs2.ProcEpoch || vs.Index != vs2.Index {
		t.Fatalf("JSON round-trip mismatch")
	}
	if vs.Amount.Cmp(vs2.Amount) != 0 {
		t.Fatalf("Amount mismatch: %s != %s", vs.Amount, vs2.Amount)
	}
}

func TestFinalizationParamsCopy(t *testing.T) {
	fp := &FinalizationParams{
		Spines:    common.HashArray{common.HexToHash("0x01"), common.HexToHash("0x02")},
		BaseSpine: func() *common.Hash { h := common.HexToHash("0x03"); return &h }(),
		Checkpoint: &Checkpoint{
			Epoch: 10, FinEpoch: 8,
			Root: common.HexToHash("0xaa"), Spine: common.HexToHash("0xbb"),
		},
		SyncMode: MainSync,
	}

	cp := fp.Copy()
	cp.Spines[0] = common.HexToHash("0xff")
	if fp.Spines[0] == cp.Spines[0] {
		t.Fatal("Copy Spines should be independent")
	}
}

func TestFinalizationParamsJSON(t *testing.T) {
	fp := &FinalizationParams{
		Spines:   common.HashArray{common.HexToHash("0x01")},
		SyncMode: HeadSync,
	}

	data, err := json.Marshal(fp)
	if err != nil {
		t.Fatal(err)
	}

	var fp2 FinalizationParams
	if err := json.Unmarshal(data, &fp2); err != nil {
		t.Fatal(err)
	}
	if fp.SyncMode != fp2.SyncMode {
		t.Fatalf("SyncMode mismatch: %d != %d", fp.SyncMode, fp2.SyncMode)
	}
}
