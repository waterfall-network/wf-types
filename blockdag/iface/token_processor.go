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

// TokenProcessor handles WRC-20 and WRC-721 token operations during state transitions.
// Implementations are created per-block via TokenProcessorFactory.
type TokenProcessor interface {
	// IsToken returns true if addr is a deployed token contract address.
	IsToken(addr [20]byte) bool

	// Call executes a token operation against the current state.
	// op is the raw ABI-encoded operation bytes taken directly from tx.Data();
	// the implementation is responsible for decoding the operation.
	Call(caller Ref, token [20]byte, value *big.Int, op []byte) ([]byte, error)
}
