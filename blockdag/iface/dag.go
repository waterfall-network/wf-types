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
	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// Dag is the coordinator-facing consensus interface implemented by wf-engine's dag.Dag.
// gwat's pluginimpl wraps this interface so that IPC calls from the coordinator are
// dispatched to wf-engine's dag instead of gwat's internal one.
type Dag interface {
	HandleFinalize(data *types.FinalizationParams) *types.FinalizationResult
	HandleCoordinatedState() *types.FinalizationResult
	HandleGetCandidates(slot uint64) *types.CandidatesResult
	HandleGetOptimisticSpines(fromSpine common.Hash) *types.OptimisticSpinesResult
	HandleValidateSpines(spines common.HashArray) (bool, error)
	HandleValidateFinalization(spines common.HashArray) (bool, error)
	HandleSyncSpines(spines common.HashArray) (bool, error)
	HandleSyncSlotInfo(slotInfo types.SlotInfo) (bool, error)
}
