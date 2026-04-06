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

// Transaction holds a transaction's data needed by the plugin layer.
// gwat populates this from its *core/types.Transaction at the call site,
// keeping all LGPL code out of Apache 2.0 modules.
type Transaction struct {
	Raw  []byte // RLP-encoded full transaction (for signing / marshaling)
	Data []byte // transaction input data field (tx.Data())
}
