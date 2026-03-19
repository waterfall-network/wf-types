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

// TxPool abstracts access to the pending transaction pool.
// gwat implements this over *core.TxPool.
type TxPool interface {
	// PendingCount returns the number of pending transactions.
	PendingCount() int

	// AddLocal enqueues a single transaction from a local source.
	AddLocal(tx []byte) error

	// Pending returns pending transactions encoded as bytes,
	// keyed by sender address (hex string).
	Pending() (map[string][][]byte, error)
}
