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
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// BlockCreator is the minimal interface that wf-engine needs from the gwat block creator.
// gwat's dag/creator.Creator satisfies this interface.
type BlockCreator interface {
	// Start activates the block creator.
	Start()
	// Stop deactivates the block creator.
	Stop()
	// IsRunning reports whether the block creator is currently active.
	IsRunning() bool
	// IsCreatorActive reports whether the given address is registered as an active creator node.
	IsCreatorActive(coinbase common.Address) bool
	// RunBlockCreation triggers creation of a block for the given slot.
	RunBlockCreation(slot uint64, slotCreators []common.Address, tips types.Tips, checkpoint *types.Checkpoint) error
}
