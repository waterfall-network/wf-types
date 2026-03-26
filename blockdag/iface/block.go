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
	"errors"
	"sort"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// Block is the read-only interface that wf-engine needs from a DAG block.
// gwat implements this over *core/types.Block; the concrete type satisfies
// this interface without any changes since all fields are exposed as methods.
type Block interface {
	// Nr returns the finalized block number (0 if not yet finalized).
	Nr() uint64
	// Number returns the finalized block number as a pointer (nil if not finalized).
	Number() *uint64
	// Slot returns the DAG slot of the block.
	Slot() uint64
	// Height returns the DAG height of the block.
	Height() uint64
	// Hash returns the canonical block hash.
	Hash() common.Hash
	// ParentHashes returns the list of parent block hashes.
	ParentHashes() common.HashArray
}

/********** BlockMap **********/

// BlockMap represents a map of blocks keyed by hash.
type BlockMap map[common.Hash]Block

// ToArray returns all blocks in the map as a slice (order is undefined).
func (bm BlockMap) ToArray() []Block {
	arr := make([]Block, 0, len(bm))
	for _, b := range bm {
		arr = append(arr, b)
	}
	return arr
}

/********** Blocks **********/

// Blocks is an ordered slice of Block.
type Blocks []Block

// GetHashes returns the hashes of all blocks in the slice.
func (bs Blocks) GetHashes() *common.HashArray {
	hashes := make(common.HashArray, 0, len(bs))
	for _, block := range bs {
		hashes = append(hashes, block.Hash())
	}
	return &hashes
}

// GroupBySlot groups blocks by their slot number.
func (bs *Blocks) GroupBySlot() (map[uint64]Blocks, error) {
	if len(*bs) == 0 {
		return map[uint64]Blocks{}, nil
	}
	res := make(map[uint64]Blocks)
	for _, block := range *bs {
		if block == nil {
			return nil, errors.New("nil block found")
		}
		slot := block.Slot()
		res[slot] = append(res[slot], block)
	}
	return res, nil
}

/********** SlotSpineMap **********/

// SlotSpineMap maps slot numbers to their representative spine block.
type SlotSpineMap map[uint64]Block

// GetOrderedHashes returns block hashes sorted in ascending slot order.
func (ssm *SlotSpineMap) GetOrderedHashes() *common.HashArray {
	if len(*ssm) == 0 {
		return &common.HashArray{}
	}
	hashes := make(common.HashArray, 0, len(*ssm))
	slots := make(common.SorterAscU64, 0, len(*ssm))
	for sl := range *ssm {
		slots = append(slots, sl)
	}
	sort.Sort(slots)
	for _, slot := range slots {
		hashes = append(hashes, (*ssm)[slot].Hash())
	}
	return &hashes
}

/********** HeaderMap **********/

// HeaderMap maps block hashes to their lightweight header metadata.
type HeaderMap map[common.Hash]*types.HeaderInfo
