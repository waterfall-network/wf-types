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

// ValidatorProcessor handles validator lifecycle operations and block reward
// distribution during state transitions.
// Implementations are created per-block via ValidatorProcessorFactory.
type ValidatorProcessor interface {
	// IsValidatorOp returns true if to is the validators state contract address.
	IsValidatorOp(to *[20]byte) bool

	// Call executes a validator operation against the current state.
	// data is the raw operation bytes taken directly from tx.Data();
	// txHash is the hash of the enclosing transaction.
	// The implementation is responsible for decoding the operation from data.
	Call(caller Ref, to [20]byte, value *big.Int, data []byte, txHash [32]byte) ([]byte, error)

	// GetValidatorsStateAddress returns the address of the validators state contract.
	GetValidatorsStateAddress() [20]byte

	// ProcessRewards distributes the block reward for coinbase, respecting
	// delegating stake profit-share rules where applicable.
	// The implementation uses its own internal state reference.
	ProcessRewards(coinbase [20]byte, reward *big.Int) error
}
