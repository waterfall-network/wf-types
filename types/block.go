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

import "gitlab.waterfall.network/waterfall/protocol/wf-types/common"

// HeaderInfo is a lightweight block header metadata struct
// used across Apache 2.0 modules to avoid GPL/LGPL dependencies.
// The canonical hash is stored and must be set by the gwat adapter layer.
type HeaderInfo struct {
	Hash         common.Hash
	ParentHashes common.HashArray
	Slot         uint64
	Era          uint64
	Height       uint64
	GasLimit     uint64
	Time         uint64
	Coinbase     common.Address
}

// BlockInfo is a lightweight block metadata struct
// used across Apache 2.0 modules.
type BlockInfo struct {
	Header HeaderInfo
	// TxCount is the number of transactions in the block.
	TxCount int
}
