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
	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// Downloader is the minimal interface that wf-engine needs from the gwat sync downloader.
// gwat's eth/downloader.Downloader satisfies this interface.
type Downloader interface {
	// Synchronising reports whether a sync is currently in progress.
	Synchronising() bool
	// OptimisticSpineSync syncs blocks reachable from the given spine hashes without
	// applying full state transitions.
	OptimisticSpineSync(spines common.HashArray) error
	// MainSync performs a full finalization sync starting from baseSpine.
	MainSync(baseSpine common.Hash, spines common.HashArray) error
	// DagSync retrieves all blocks starting from baseSpine and adds them to the
	// non-finalized chain.
	DagSync(baseSpine common.Hash, spines common.HashArray) error
	// Terminate interrupts the downloader, cancelling all pending operations.
	Terminate()
}
