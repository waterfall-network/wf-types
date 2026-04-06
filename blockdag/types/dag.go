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

package types

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math/big"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/common/hexutil"
)

// Checkpoint represents a coordinated checkpoint of coordinator and gwat nodes.
type Checkpoint struct {
	Epoch    uint64
	FinEpoch uint64
	Root     common.Hash
	Spine    common.Hash
}

type checkpointMarshaling struct {
	Epoch    *hexutil.Uint64 `json:"epoch"`
	FinEpoch *hexutil.Uint64 `json:"finEpoch"`
	Root     *common.Hash    `json:"root"`
	Spine    *common.Hash    `json:"spine"`
}

// Bytes encodes Checkpoint to its byte representation.
func (cp *Checkpoint) Bytes() []byte {
	cpLen := 8 + 8 + common.HashLength + common.HashLength
	res := make([]byte, 0, cpLen)
	epoch := make([]byte, 8)
	binary.BigEndian.PutUint64(epoch, cp.Epoch)
	finEpoch := make([]byte, 8)
	binary.BigEndian.PutUint64(finEpoch, cp.FinEpoch)
	res = append(res, epoch...)
	res = append(res, finEpoch...)
	res = append(res, cp.Root.Bytes()...)
	res = append(res, cp.Spine.Bytes()...)
	return res
}

// SetBytes restores Checkpoint from its byte representation.
func (cp *Checkpoint) SetBytes(data []byte) error {
	cpLen := 8 + 8 + common.HashLength + common.HashLength
	if len(data) != cpLen {
		return fmt.Errorf("bad bitlen: got=%d req=%d", len(data), cpLen)
	}
	var start, end int
	end += 8
	cp.Epoch = binary.BigEndian.Uint64(data[start:end])

	start = end
	end += 8
	cp.FinEpoch = binary.BigEndian.Uint64(data[start:end])

	start = end
	end += common.HashLength
	cp.Root = common.BytesToHash(data[start:end])

	start = end
	end += common.HashLength
	cp.Spine = common.BytesToHash(data[start:end])

	return nil
}

// BytesToCheckpoint creates a Checkpoint from its byte representation.
func BytesToCheckpoint(b []byte) (*Checkpoint, error) {
	var h Checkpoint
	if err := h.SetBytes(b); err != nil {
		return nil, err
	}
	return &h, nil
}

// Copy returns a deep copy of cp.
func (cp *Checkpoint) Copy() *Checkpoint {
	cpy := &Checkpoint{
		Epoch:    cp.Epoch,
		FinEpoch: cp.FinEpoch,
	}
	copy(cpy.Root[:], cp.Root[:])
	copy(cpy.Spine[:], cp.Spine[:])
	return cpy
}

// MarshalJSON implements json.Marshaler.
func (cp *Checkpoint) MarshalJSON() ([]byte, error) {
	out := checkpointMarshaling{
		Epoch:    (*hexutil.Uint64)(&cp.Epoch),
		FinEpoch: (*hexutil.Uint64)(&cp.FinEpoch),
		Root:     &cp.Root,
		Spine:    &cp.Spine,
	}
	return json.Marshal(out)
}

// UnmarshalJSON implements json.Unmarshaler.
func (cp *Checkpoint) UnmarshalJSON(input []byte) error {
	var dec checkpointMarshaling
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Epoch != nil {
		cp.Epoch = uint64(*dec.Epoch)
	}
	if dec.FinEpoch != nil {
		cp.FinEpoch = uint64(*dec.FinEpoch)
	}
	if dec.Root != nil {
		cp.Root = *dec.Root
	}
	if dec.Spine != nil {
		cp.Spine = *dec.Spine
	}
	return nil
}

// ValidatorSyncOp represents the type of validator synchronization operation.
type ValidatorSyncOp uint64

const (
	Activate      ValidatorSyncOp = iota
	Deactivate    ValidatorSyncOp = iota
	UpdateBalance ValidatorSyncOp = iota
)

// ValidatorSync holds data for a validator synchronization operation from coordinator.
type ValidatorSync struct {
	InitTxHash      common.Hash
	OpType          ValidatorSyncOp
	ProcEpoch       uint64
	Index           uint64
	Creator         common.Address
	Amount          *big.Int
	TxHash          *common.Hash
	Balance         *big.Int
	ActivationEpoch uint64
	ExitEpoch       uint64
}

type validatorSyncMarshaling struct {
	InitTxHash      *common.Hash    `json:"initTxHash"`
	OpType          *hexutil.Uint64 `json:"opType"`
	ProcEpoch       *hexutil.Uint64 `json:"procEpoch"`
	Index           *hexutil.Uint64 `json:"index"`
	Creator         *common.Address `json:"creator"`
	Amount          *hexutil.Big    `json:"amount"`
	TxHash          *common.Hash    `json:"txHash"`
	Balance         *hexutil.Big    `json:"balance"`
	ActivationEpoch *hexutil.Uint64 `json:"activationEpoch"`
	ExitEpoch       *hexutil.Uint64 `json:"exitEpoch"`
}

// Copy returns a deep copy of vs.
func (vs *ValidatorSync) Copy() *ValidatorSync {
	cpy := &ValidatorSync{
		OpType:          vs.OpType,
		ProcEpoch:       vs.ProcEpoch,
		Index:           vs.Index,
		ActivationEpoch: vs.ActivationEpoch,
		ExitEpoch:       vs.ExitEpoch,
	}
	copy(cpy.Creator[:], vs.Creator[:])
	copy(cpy.InitTxHash[:], vs.InitTxHash[:])
	if vs.Amount != nil {
		cpy.Amount = new(big.Int).Set(vs.Amount)
	}
	if vs.TxHash != nil {
		cpy.TxHash = new(common.Hash)
		copy(cpy.TxHash[:], vs.TxHash[:])
	}
	if vs.Balance != nil {
		cpy.Balance = new(big.Int).Set(vs.Balance)
	}
	return cpy
}

// Print returns a string representation of vs.
func (vs *ValidatorSync) Print() string {
	if vs == nil {
		return "{nil}"
	}
	amount := "<nil>"
	if vs.Amount != nil {
		amount = vs.Amount.String()
	}
	balance := "<nil>"
	if vs.Balance != nil {
		balance = vs.Balance.String()
	}
	return fmt.Sprintf(
		"{InitTxHash: %#x, OpType: %d, ProcEpoch: %d, Index: %d, Creator: %#x, Amount: %s, Balance: %s, TxHash: %#x, ActivationEpoch: %d, ExitEpoch: %d}",
		vs.InitTxHash, vs.OpType, vs.ProcEpoch, vs.Index, vs.Creator,
		amount, balance, vs.TxHash, vs.ActivationEpoch, vs.ExitEpoch,
	)
}

// Key returns the identifying hash of vs.
func (vs *ValidatorSync) Key() common.Hash {
	if vs == nil {
		return common.Hash{}
	}
	return vs.InitTxHash
}

// MarshalJSON implements json.Marshaler.
func (vs *ValidatorSync) MarshalJSON() ([]byte, error) {
	out := validatorSyncMarshaling{
		OpType:          (*hexutil.Uint64)(&vs.OpType),
		ProcEpoch:       (*hexutil.Uint64)(&vs.ProcEpoch),
		Index:           (*hexutil.Uint64)(&vs.Index),
		Creator:         &vs.Creator,
		TxHash:          vs.TxHash,
		InitTxHash:      &vs.InitTxHash,
		ActivationEpoch: (*hexutil.Uint64)(&vs.ActivationEpoch),
		ExitEpoch:       (*hexutil.Uint64)(&vs.ExitEpoch),
	}
	if vs.Amount != nil {
		out.Amount = (*hexutil.Big)(vs.Amount)
	}
	if vs.Balance != nil {
		out.Balance = (*hexutil.Big)(vs.Balance)
	}
	return json.Marshal(out)
}

// UnmarshalJSON implements json.Unmarshaler.
func (vs *ValidatorSync) UnmarshalJSON(input []byte) error {
	var dec validatorSyncMarshaling
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.OpType != nil {
		vs.OpType = ValidatorSyncOp(*dec.OpType)
	}
	if dec.ProcEpoch != nil {
		vs.ProcEpoch = uint64(*dec.ProcEpoch)
	}
	if dec.Index != nil {
		vs.Index = uint64(*dec.Index)
	}
	if dec.Creator != nil {
		vs.Creator = *dec.Creator
	}
	if dec.Amount != nil {
		vs.Amount = (*big.Int)(dec.Amount)
	}
	if dec.TxHash != nil {
		vs.TxHash = dec.TxHash
	}
	if dec.InitTxHash != nil {
		vs.InitTxHash = *dec.InitTxHash
	}
	if dec.Balance != nil {
		vs.Balance = (*big.Int)(dec.Balance)
	}
	if dec.ActivationEpoch != nil {
		vs.ActivationEpoch = uint64(*dec.ActivationEpoch)
	}
	if dec.ExitEpoch != nil {
		vs.ExitEpoch = uint64(*dec.ExitEpoch)
	}
	return nil
}

// SyncMode specifies the synchronization mode for finalization.
type SyncMode uint8

const (
	NoSync   SyncMode = iota
	MainSync SyncMode = iota
	HeadSync SyncMode = iota
)

// FinalizationParams holds parameters for a finalization request.
type FinalizationParams struct {
	Spines      common.HashArray `json:"spines"`
	BaseSpine   *common.Hash     `json:"baseSpine"`
	Checkpoint  *Checkpoint      `json:"checkpoint"`
	ValSyncData []*ValidatorSync `json:"valSyncData"`
	SyncMode    SyncMode         `json:"syncMode"`
}

type finalizationParamsMarshaling struct {
	Spines      common.HashArray `json:"spines"`
	BaseSpine   *common.Hash     `json:"baseSpine"`
	Checkpoint  *Checkpoint      `json:"checkpoint"`
	ValSyncData []*ValidatorSync `json:"valSyncData"`
	SyncMode    *hexutil.Uint8   `json:"syncMode"`
}

// Copy returns a deep copy of fp.
func (fp *FinalizationParams) Copy() *FinalizationParams {
	cpy := &FinalizationParams{
		Spines: fp.Spines.Copy(),
	}
	if fp.BaseSpine != nil {
		cpy.BaseSpine = &common.Hash{}
		copy(cpy.BaseSpine[:], fp.BaseSpine[:])
	}
	if fp.Checkpoint != nil {
		cpy.Checkpoint = fp.Checkpoint.Copy()
	}
	if fp.ValSyncData != nil {
		cpy.ValSyncData = make([]*ValidatorSync, len(fp.ValSyncData))
		for i, vs := range fp.ValSyncData {
			cpy.ValSyncData[i] = vs.Copy()
		}
	}
	cpy.SyncMode = fp.SyncMode
	return cpy
}

// MarshalJSON implements json.Marshaler.
func (fp *FinalizationParams) MarshalJSON() ([]byte, error) {
	out := finalizationParamsMarshaling{
		Spines:      fp.Spines,
		BaseSpine:   fp.BaseSpine,
		Checkpoint:  fp.Checkpoint,
		ValSyncData: fp.ValSyncData,
		SyncMode:    (*hexutil.Uint8)(&fp.SyncMode),
	}
	return json.Marshal(out)
}

// UnmarshalJSON implements json.Unmarshaler.
func (fp *FinalizationParams) UnmarshalJSON(input []byte) error {
	type Decoding finalizationParamsMarshaling
	dec := Decoding{}
	if err := json.Unmarshal(input, &dec); err != nil {
		return err
	}
	if dec.Spines != nil {
		fp.Spines = dec.Spines
	}
	if dec.BaseSpine != nil {
		fp.BaseSpine = dec.BaseSpine
	}
	if dec.Checkpoint != nil {
		fp.Checkpoint = dec.Checkpoint
	}
	if dec.ValSyncData != nil {
		fp.ValSyncData = dec.ValSyncData
	}
	if dec.SyncMode != nil {
		fp.SyncMode = SyncMode(*dec.SyncMode)
	}
	return nil
}

// FinalizationResult holds the result of a finalization operation.
type FinalizationResult struct {
	Error   *string      `json:"error"`
	LFSpine *common.Hash `json:"lfSpine"`
	CpEpoch *uint64      `json:"cpEpoch"`
	CpRoot  *common.Hash `json:"cpRoot"`
}

// ConsensusResult holds the result of a consensus request.
type ConsensusResult struct {
	Error      *string            `json:"error"`
	Info       *map[string]string `json:"info"`
	Candidates common.HashArray   `json:"candidates"`
}

// CandidatesResult holds the result of a candidates request.
type CandidatesResult struct {
	Error      *string          `json:"error"`
	Candidates common.HashArray `json:"candidates"`
}

// OptimisticSpinesResult holds the result of an optimistic spines request.
type OptimisticSpinesResult struct {
	Data  []common.HashArray `json:"data"`
	Error *string            `json:"error"`
}
