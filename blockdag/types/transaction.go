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

package types

// Transaction holds a serialized (RLP-encoded) transaction.
// It is used as a boundary type between wf-engine and gwat:
// gwat converts its *core/types.Transaction to/from Transaction
// at the call site, keeping all LGPL code out of Apache 2.0 modules.
type Transaction struct {
	Raw []byte
}
