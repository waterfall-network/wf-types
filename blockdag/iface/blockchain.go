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
	"context"
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

// CreatorsBySlotStorage is the narrow validator-storage view used by the dag package.
// gwat implements this over its ValidatorStorage, returning the shuffled creator
// addresses for the given slot.
type CreatorsBySlotStorage interface {
	GetCreatorsBySlot(slot uint64) ([]common.Address, error)
}

// BlockChain is the interface that wf-engine's dag package needs from the gwat
// execution layer. Signatures match gwat's *core.BlockChain methods exactly so
// that gwat can implement this interface without any adaptation layer.
//
// It is the union of dag.blockChain and dag/finalizer.blockChain from wf-engine.
type BlockChain interface {
	Config() *ChainConfig

	// --- Block queries ---
	GetLastFinalizedBlock() Block
	GetBlockByHash(hash common.Hash) Block
	GetBlocksByHashes(hashes common.HashArray) BlockMap
	GetBlock(ctx context.Context, hash common.Hash) Block
	GetLastFinalizedNumber() uint64
	GetLastFinalizedHeader() *types.HeaderInfo
	GetHeaderByHash(hash common.Hash) *types.HeaderInfo
	GetHeaderByNumber(number uint64) *types.HeaderInfo
	GetHeadersByHashes(hashes common.HashArray) HeaderMap
	GetBlockDag(hash common.Hash) *types.BlockDAG
	GetTips() types.Tips
	GetOptimisticSpines(gtSlot uint64) ([]common.HashArray, error)
	GetFinalizedNumberByHash(hash common.Hash) *uint64
	GetFinalizedHashByNumber(number uint64) common.Hash
	GetBlockFinalizedNumber(hash common.Hash) *uint64

	// --- Finalization mutations ---
	FinalizeBlock(blockHash common.Hash, finNr uint64, stateBlockHash common.Hash, isHead bool) error
	FinalizeTips(finHashes common.HashArray, lastFinHash common.Hash, lastFinHeight uint64)
	RollbackFinalization(spineHash common.Hash, lfNr uint64) error
	SaveBlockDag(blockDag *types.BlockDAG)
	SetRollbackActive()
	ResetRollbackActive()

	// --- DAG traversal ---
	CollectAncestorsAftCpByTips(parents common.HashArray, cpHash common.Hash) (bool, HeaderMap, common.HashArray, types.Tips)
	CollectAncestorsAftCpByParents(parents common.HashArray, cpHash common.Hash) (bool, HeaderMap, common.HashArray, error)

	// --- Slot / era / sync state ---
	GetSlotInfo() *types.SlotInfo
	SetSlotInfo(si *types.SlotInfo) error
	GetEraInfo() *era.EraInfo
	HandleEra(cp *types.Checkpoint) error
	GetLastCoordinatedCheckpoint() *types.Checkpoint
	SetLastCoordinatedCheckpoint(cp *types.Checkpoint)
	IsSynced() bool
	SetIsSynced(synced bool)
	ResetSyncCheckpointCache()
	SetSyncCheckpointCache(cp *types.Checkpoint)

	// --- Validator sync ---
	AppendNotProcessedValidatorSyncData(valSyncData []*types.ValidatorSync)
	CleanInvalidNotProcessedValidatorSync()
	FixValidatorSyncOps(finEpoch uint64, valSyncData []*types.ValidatorSync) []*types.ValidatorSync

	// --- Mutex ---
	DagMuLock()
	DagMuUnlock()

	// --- Validator storage (narrow dag view: slot → creator addresses) ---
	ValidatorStorage() CreatorsBySlotStorage
}
