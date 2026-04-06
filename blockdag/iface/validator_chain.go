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

package iface

import (
	"math/big"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types/era"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// ProcessorCtx is the subset of *validator.Processor that gwat's
// FixValidatorSyncOpProcessing fix needs. Defining it here avoids a circular
// import between wf-engine/validator and wf-types.
type ProcessorCtx interface {
	GetBlockContext() BlockContext
}

// ValidatorSyncOpFix is the subset of operation.ValidatorSync that gwat's
// chain-fix functions read. operation.ValidatorSync from wf-engine satisfies
// this interface without any changes.
type ValidatorSyncOpFix interface {
	InitTxHash() common.Hash
	OpType() types.ValidatorSyncOp
	ProcEpoch() uint64
	Index() uint64
	Creator() common.Address
	Amount() *big.Int
}

// ValidatorChain is the blockchain interface needed by the validator and
// validatorsync packages in wf-engine. gwat implements this over *core.BlockChain.
//
// It is the union of validator.blockchain and validatorsync.blockchain from
// wf-engine, excluding two methods that cannot be expressed in wf-types:
//   - ValidatorStorage() — returns valStore.Storage from wf-engine (uses *Validator)
//   - FixValidatorSyncOpProcessing() — takes *validator.Processor (circular import)
//
// Those two methods are added locally in wf-engine on top of this interface.
type ValidatorChain interface {
	Config() *ChainConfig

	// --- Slot / era ---
	GetSlotInfo() *types.SlotInfo
	GetEraInfo() *era.EraInfo
	EpochToEra(epoch uint64) *era.Era
	ReadEra(eraNumber uint64) *era.Era

	// --- Headers ---
	GetHeaderByHash(hash common.Hash) *types.HeaderInfo

	// --- State ---
	StateAt(root common.Hash) (StateDB, error)

	// --- Validator sync data ---
	GetValidatorSyncData(initTxHash common.Hash) *types.ValidatorSync
	GetNotProcessedValidatorSyncData() map[common.Hash]*types.ValidatorSync
	SetValidatorSyncData(vs *types.ValidatorSync)

	// --- Transactions ---
	GetTransaction(txHash common.Hash) (*types.Transaction, common.Hash, uint64)
	GetTransactionReceipt(txHash common.Hash) (*types.Receipt, common.Hash, uint64)
	// GetTransactionSender recovers the sender of the transaction.
	// Replaces types.LatestSigner + types.Sender from gwat.
	GetTransactionSender(txHash common.Hash) (common.Address, error)

	// --- Hard-coded chain fixes ---
	// FixValidatorSyncOpProcessing applies mainnet-specific fixes during
	// validator sync op processing. gwat implements this in core/chain_fixes.go.
	// The processor performs the ValidatorSync type-assertion before calling this.
	FixValidatorSyncOpProcessing(ctx ProcessorCtx, op ValidatorSyncOpFix, txHash common.Hash, from, to common.Address) (isApplied bool, ret []byte, err error)
}
