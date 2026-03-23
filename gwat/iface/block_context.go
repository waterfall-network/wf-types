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

import "math/big"

// BlockContext holds block-level information passed to token and validator processors.
// It mirrors the data fields of gwat's vm.BlockContext, excluding EVM-specific
// function callbacks (CanTransfer, Transfer, GetHash) which are not needed
// outside the execution layer.
type BlockContext struct {
	// Coinbase is the block producer address (EVM opcode COINBASE).
	Coinbase [20]byte
	// GasLimit is the block gas limit (EVM opcode GASLIMIT).
	GasLimit uint64

	// BlockHeight is the Waterfall-specific block height.
	BlockHeight *big.Int
	// BlockNumber is the EVM block number (EVM opcode NUMBER).
	BlockNumber *big.Int
	// Time is the block timestamp in seconds (EVM opcode TIME).
	Time *big.Int
	// Difficulty is the block difficulty (EVM opcode DIFFICULTY).
	Difficulty *big.Int
	// BaseFee is the base fee per gas unit (EIP-1559, EVM opcode BASEFEE).
	// Nil before EIP-1559 activation.
	BaseFee *big.Int

	// Random holds the post-merge PREVRANDAO value (EVM opcode PREVRANDAO).
	// Nil before the merge.
	Random *[32]byte

	// Slot is the Waterfall consensus slot number.
	Slot uint64
	// Era is the Waterfall consensus era number.
	Era uint64
	// BlockHash is the hash of the current block.
	BlockHash [32]byte
}
