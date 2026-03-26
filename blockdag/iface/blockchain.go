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

package iface

import (
	"math/big"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types/era"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// Database abstracts low-level key-value storage.
type Database interface {
	Get(key []byte) ([]byte, error)
	Put(key []byte, value []byte) error
	Delete(key []byte) error
	Has(key []byte) (bool, error)
	Close() error
}

// ChainConfig holds chain parameters needed by the plugin layer.
type ChainConfig struct {
	ChainID           uint64
	EpochsPerEra      uint64
	SlotsPerEpoch     uint64
	ValidatorsPerSlot uint64
	TransitionPeriod  uint64
	StartEpochsPerEra uint64
	// AcceptCpRootOnFinEpoch is a hard-coded safelist of (cpRoot → []finEpoch) pairs
	// used to skip invalid finalization requests on specific networks.
	AcceptCpRootOnFinEpoch map[[32]byte][]uint64

	// Validator state contract addresses.
	ValidatorsStateAddress    common.Address
	WaterfallDummyAddress     common.Address
	AllocationContractAddress common.Address

	// Validator consensus parameters.
	EffectiveBalance       *big.Int // effective validator balance in WAT (not Wei)
	ValidatorOpExpireSlots uint64   // number of slots before a validator op expires

	// Fork slots — enable features when the current slot >= the fork slot.
	ForkSlotDelegate      uint64
	ForkSlotValSyncProc   uint64
	ForkSlotValOpTracking uint64
}

// IsForkSlotDelegate reports whether the delegate fork is active at the given slot.
func (c *ChainConfig) IsForkSlotDelegate(slot uint64) bool {
	return slot >= c.ForkSlotDelegate
}

// IsForkSlotValSyncProc reports whether the validator-sync-proc fork is active at the given slot.
func (c *ChainConfig) IsForkSlotValSyncProc(slot uint64) bool {
	return slot >= c.ForkSlotValSyncProc
}

// IsForkSlotValOpTracking reports whether the validator-op-tracking fork is active at the given slot.
func (c *ChainConfig) IsForkSlotValOpTracking(slot uint64) bool {
	return slot >= c.ForkSlotValOpTracking
}

// BlockChain is the minimal interface that wf-engine and wf-consensus
// need from the gwat execution layer. gwat implements this over *core.BlockChain.
type BlockChain interface {
	// Block queries
	GetBlock(hash [32]byte) (*types.BlockInfo, error)
	GetBlockByNumber(number uint64) (*types.BlockInfo, error)
	GetLastFinalizedBlock() (*types.BlockInfo, error)
	GetLastFinalizedNumber() (uint64, error)
	GetHeaderByHash(hash [32]byte) (*types.HeaderInfo, error)

	// DAG queries
	GetTips() (types.Tips, error)
	GetSlotInfo() (*types.SlotInfo, error)
	SetSlotInfo(si *types.SlotInfo) error
	GetEraInfo() (*era.EraInfo, error)
	GetLastCoordinatedCheckpoint() (*types.Checkpoint, error)
	GetEpoch(epoch uint64) ([32]byte, error)

	// State
	StateAt(root [32]byte) (StateDB, error)

	// Config and DB
	Config() *ChainConfig
	Database() Database

	// DAG mutations
	InsertBlockDag(block *types.BlockInfo) error
}
