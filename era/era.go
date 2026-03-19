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

// Package era implements era management for the Waterfall blockchain.
package era

import (
	"math"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/types"
)

// EraConfig holds chain configuration parameters relevant to era management.
type EraConfig struct {
	EpochsPerEra      uint64
	SlotsPerEpoch     uint64
	ValidatorsPerSlot uint64
	TransitionPeriod  uint64
	StartEpochsPerEra uint64
}

// Blockchain defines the interface used by era management to query chain state.
type Blockchain interface {
	GetSlotInfo() *types.SlotInfo
	GetLastCoordinatedCheckpoint() *types.Checkpoint
	GetEraInfo() *EraInfo
	Config() *EraConfig
}

// Era represents a period of time in the Waterfall blockchain.
type Era struct {
	Number    uint64      `json:"number"`
	From      uint64      `json:"fromEpoch"`
	To        uint64      `json:"toEpoch"`
	Root      common.Hash `json:"root"`
	BlockHash common.Hash `json:"blockHash"`
}

// NewEra creates a new Era.
func NewEra(number, from, to uint64, root, blockHash common.Hash) *Era {
	return &Era{
		Number:    number,
		From:      from,
		To:        to,
		Root:      root,
		BlockHash: blockHash,
	}
}

// Length returns the number of epochs in the era.
func (e *Era) Length() uint64 {
	return e.To - e.From + 1
}

// IsContainsEpoch returns true if epoch falls within this era.
func (e *Era) IsContainsEpoch(epoch uint64) bool {
	return epoch >= e.From && epoch <= e.To
}

// EraInfo holds metadata about the current era.
type EraInfo struct {
	currentEra *Era
	length     uint64
}

// NewEraInfo creates a new EraInfo.
func NewEraInfo(era *Era) *EraInfo {
	return &EraInfo{
		currentEra: era,
		length:     era.Length(),
	}
}

// Number returns the era number.
func (ei *EraInfo) Number() uint64 { return ei.currentEra.Number }

// GetEra returns the current Era object.
func (ei *EraInfo) GetEra() *Era { return ei.currentEra }

// ToEpoch returns the last epoch of the current era.
func (ei *EraInfo) ToEpoch() uint64 { return ei.GetEra().To }

// FromEpoch returns the first epoch of the current era.
func (ei *EraInfo) FromEpoch() uint64 { return ei.GetEra().From }

// EpochsPerEra returns the total number of epochs in the era.
func (ei *EraInfo) EpochsPerEra() uint64 {
	return ei.GetEra().To - ei.GetEra().From + 1
}

// FirstEpoch returns the first epoch of the era.
func (ei *EraInfo) FirstEpoch() uint64 { return ei.FromEpoch() }

// FirstSlot returns the first slot of the era.
func (ei *EraInfo) FirstSlot(bc Blockchain) uint64 {
	slot, err := bc.GetSlotInfo().SlotOfEpochStart(ei.FirstEpoch())
	if err != nil {
		return 0
	}
	return slot
}

// LastEpoch returns the last epoch of the era.
func (ei *EraInfo) LastEpoch() uint64 { return ei.ToEpoch() }

// LastSlot returns the last slot of the era.
func (ei *EraInfo) LastSlot(bc Blockchain) uint64 {
	slot, err := bc.GetSlotInfo().SlotOfEpochEnd(ei.LastEpoch())
	if err != nil {
		return 0
	}
	return slot
}

// IsTransitionPeriodEpoch returns true if epoch is within the transition period.
func (ei *EraInfo) IsTransitionPeriodEpoch(bc Blockchain, epoch uint64) bool {
	tp := bc.Config().TransitionPeriod
	return epoch >= ei.ToEpoch()-tp && epoch <= ei.ToEpoch()
}

// IsTransitionPeriodStartEpoch returns true if epoch is the start of the transition period.
func (ei *EraInfo) IsTransitionPeriodStartEpoch(bc Blockchain, epoch uint64) bool {
	return epoch == (ei.ToEpoch() - bc.Config().TransitionPeriod)
}

// IsTransitionPeriodStartSlot returns true if slot is the start of the transition period.
func (ei *EraInfo) IsTransitionPeriodStartSlot(bc Blockchain, slot uint64) bool {
	transitionEpoch := ei.ToEpoch() + 1 - bc.Config().TransitionPeriod
	currentEpoch := bc.GetSlotInfo().SlotToEpoch(slot)
	return currentEpoch == transitionEpoch && bc.GetSlotInfo().IsEpochStart(slot)
}

// NextEraFirstEpoch returns the first epoch of the next era.
func (ei *EraInfo) NextEraFirstEpoch() uint64 { return ei.ToEpoch() + 1 }

// NextEraFirstSlot returns the first slot of the next era.
func (ei *EraInfo) NextEraFirstSlot(bc Blockchain) uint64 {
	slot, err := bc.GetSlotInfo().SlotOfEpochStart(ei.NextEraFirstEpoch())
	if err != nil {
		return 0
	}
	return slot
}

// LenEpochs returns the total number of epochs in the current era.
func (ei *EraInfo) LenEpochs() uint64 { return ei.length }

// LenSlots returns the total number of slots in the current era.
func (ei *EraInfo) LenSlots() uint64 { return ei.length * 32 }

// IsContainsEpoch returns true if epoch is within the current era.
func (ei *EraInfo) IsContainsEpoch(epoch uint64) bool {
	return epoch >= ei.FromEpoch() && epoch <= ei.ToEpoch()
}

// NextEra calculates the next era from the current blockchain state.
func NextEra(bc Blockchain, root, blockHash common.Hash, numValidators uint64) *Era {
	nextEraNumber := bc.GetEraInfo().Number() + 1
	nextEraLength := EstimateEraLength(bc.Config(), numValidators, nextEraNumber)
	nextEraBegin := bc.GetEraInfo().ToEpoch() + 1
	nextEraEnd := bc.GetEraInfo().ToEpoch() + nextEraLength
	return NewEra(nextEraNumber, nextEraBegin, nextEraEnd, root, blockHash)
}

// EstimateEraLength calculates the length of an era based on chain configuration.
func EstimateEraLength(cfg *EraConfig, numberOfValidators, eraNumber uint64) uint64 {
	if eraNumber >= cfg.StartEpochsPerEra {
		return cfg.EpochsPerEra
	}
	eraLength := cfg.EpochsPerEra * roundUp(
		float64(numberOfValidators)/(float64(cfg.EpochsPerEra)*float64(cfg.SlotsPerEpoch)*float64(cfg.ValidatorsPerSlot)),
	)
	return eraLength
}

func roundUp(num float64) uint64 {
	return uint64(math.Ceil(num))
}
