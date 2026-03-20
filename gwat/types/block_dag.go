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

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/bits"
	"sort"
	"time"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

/********** Tips **********/

// Tips represents the tips of the DAG chain.
type Tips map[common.Hash]*BlockDAG

// Add adds or replaces a BlockDAG in tips. Nil blockDag is silently ignored.
func (tips Tips) Add(blockDag *BlockDAG) Tips {
	if blockDag == nil {
		return tips
	}
	blockDag.OrderedAncestorsHashes.Deduplicate()
	tips[blockDag.Hash] = blockDag
	return tips
}

// Remove removes the item with the given hash from tips.
func (tips Tips) Remove(hash common.Hash) Tips {
	delete(tips, hash)
	return tips
}

// Get returns the BlockDAG for the given hash.
func (tips Tips) Get(hash common.Hash) *BlockDAG {
	return tips[hash]
}

// GetHashes returns the hashes of all BlockDAGs in tips.
func (tips Tips) GetHashes() common.HashArray {
	hashes := make(common.HashArray, 0, len(tips))
	for _, b := range tips {
		hashes = append(hashes, b.Hash)
	}
	return hashes
}

// Copy returns a shallow copy of tips.
func (tips Tips) Copy() Tips {
	cpy := make(Tips, len(tips))
	for key, value := range tips {
		cpy[key] = value
	}
	return cpy
}

// GetOrderedAncestorsHashes collects all DAG hashes in finalizing order.
func (tips Tips) GetOrderedAncestorsHashes() common.HashArray {
	allHashes := common.HashArray{}
	dagHashes := tips.getOrderedHashes()
	for _, dagHash := range dagHashes {
		dag := tips.Get(dagHash)
		allHashes = append(allHashes, dag.OrderedAncestorsHashes...)
		allHashes = append(allHashes, dagHash).Uniq()
	}
	return allHashes.Uniq()
}

// GetStableStateHash retrieves the hash of the stable state.
func (tips Tips) GetStableStateHash() common.Hash {
	finDag := tips.GetFinalizingDag()
	if finDag == nil {
		return common.Hash{}
	}
	return finDag.Hash
}

// GetFinalizingDag retrieves the top DAG with stable state.
func (tips Tips) GetFinalizingDag() *BlockDAG {
	dagHashes := tips.getOrderedHashes()
	if len(dagHashes) < 1 {
		return nil
	}
	return tips.Get(dagHashes[0])
}

// GetAncestorsHashes retrieves all OrderedAncestorsHashes from all tips.
func (tips Tips) GetAncestorsHashes() common.HashArray {
	ancestors := common.HashArray{}
	for _, dag := range tips {
		ancestors = append(ancestors, dag.OrderedAncestorsHashes...)
	}
	return ancestors
}

// Print returns a string representation of tips.
func (tips Tips) Print() string {
	mapB, _ := json.Marshal(tips)
	return string(mapB)
}

// getOrderedHashes sorts tips in finalizing order (desc height, then by hash).
func (tips Tips) getOrderedHashes() common.HashArray {
	heightMap := map[uint64]common.HashArray{}
	for _, dag := range tips {
		key := dag.Height
		heightMap[key] = append(heightMap[key], dag.Hash)
	}
	keys := make(common.SorterDescU64, 0, len(heightMap))
	for k := range heightMap {
		keys = append(keys, k)
	}
	sort.Sort(keys)
	sortedHashes := common.HashArray{}
	for _, k := range keys {
		sortedHashes = sortedHashes.Concat(heightMap[k].Sort())
	}
	return sortedHashes
}

// GetMaxSlot returns the maximum slot among all tips.
func (tips Tips) GetMaxSlot() uint64 {
	maxSlot := uint64(0)
	for _, t := range tips {
		if t.Slot > maxSlot {
			maxSlot = t.Slot
		}
	}
	return maxSlot
}

/********** BlockDAG **********/

// BlockDAG represents a tip block (no descendants) of the directed acyclic graph.
type BlockDAG struct {
	Hash     common.Hash
	Height   uint64
	Slot     uint64
	CpHash   common.Hash
	CpHeight uint64
	// ordered non-finalized ancestors hashes
	OrderedAncestorsHashes common.HashArray
}

// ToBytes encodes BlockDAG to its byte representation.
func (b *BlockDAG) ToBytes() []byte {
	res := []byte{}
	res = append(res, b.Hash.Bytes()...)

	height := make([]byte, 8)
	binary.BigEndian.PutUint64(height, b.Height)
	res = append(res, height...)

	slot := make([]byte, 8)
	binary.BigEndian.PutUint64(slot, b.Slot)
	res = append(res, slot...)

	res = append(res, b.CpHash.Bytes()...)

	lastFinHeight := make([]byte, 8)
	binary.BigEndian.PutUint64(lastFinHeight, b.CpHeight)
	res = append(res, lastFinHeight...)

	lenDC := make([]byte, 4)
	binary.BigEndian.PutUint32(lenDC, uint32(len(b.OrderedAncestorsHashes)))
	res = append(res, lenDC...)
	res = append(res, b.OrderedAncestorsHashes.ToBytes()...)

	return res
}

// SetBytes restores BlockDAG from its byte representation.
func (b *BlockDAG) SetBytes(data []byte) *BlockDAG {
	start := 0
	end := common.HashLength
	b.Hash = common.BytesToHash(data[start:end])

	start = end
	end += 8
	b.Height = binary.BigEndian.Uint64(data[start:end])

	start = end
	end += 8
	b.Slot = binary.BigEndian.Uint64(data[start:end])

	start = end
	end += common.HashLength
	b.CpHash = common.BytesToHash(data[start:end])

	start = end
	end += 8
	b.CpHeight = binary.BigEndian.Uint64(data[start:end])

	start = end
	end += 4
	lenDC := binary.BigEndian.Uint32(data[start:end])
	start = end
	end += common.HashLength * int(lenDC)
	b.OrderedAncestorsHashes = common.HashArrayFromBytes(data[start:end])

	return b
}

// Print returns a string representation of BlockDAG.
func (b *BlockDAG) Print() string {
	mapB, _ := json.Marshal(b)
	return string(mapB)
}

/********** SlotInfo **********/

// SlotInfo holds slot timing configuration.
type SlotInfo struct {
	GenesisTime    uint64 `json:"genesisTime"`
	SecondsPerSlot uint64 `json:"secondsPerSlot"`
	SlotsPerEpoch  uint64 `json:"slotsPerEpoch"`
}

// StartSlotTime returns the start time of the given slot.
func (si *SlotInfo) StartSlotTime(slot uint64) (time.Time, error) {
	overflows, timeSinceGenesis := bits.Mul64(slot, si.SecondsPerSlot)
	if overflows > 0 {
		return time.Unix(0, 0), fmt.Errorf("slot (%d) is in the far distant future: %w", slot, errors.New("multiplication overflows"))
	}
	sTime, carry := bits.Add64(timeSinceGenesis, si.GenesisTime, 0)
	if carry > 0 {
		return time.Unix(0, 0), fmt.Errorf("slot (%d) is in the far distant future: %w", slot, errors.New("addition overflows"))
	}
	return time.Unix(int64(sTime), 0), nil
}

// CurrentSlot returns the current slot as determined by the local clock.
func (si *SlotInfo) CurrentSlot() uint64 {
	now := time.Now().Unix()
	genesis := int64(si.GenesisTime)
	if now < genesis {
		return 0
	}
	return uint64(now-genesis) / si.SecondsPerSlot
}

// IsEpochStart returns true if slot is the first slot of an epoch.
func (si *SlotInfo) IsEpochStart(slot uint64) bool {
	return slot%si.SlotsPerEpoch == 0
}

// IsEpochEnd returns true if slot is the last slot of an epoch.
func (si *SlotInfo) IsEpochEnd(slot uint64) bool {
	return si.IsEpochStart(slot + 1)
}

// SlotToEpoch returns the epoch number of the given slot.
func (si *SlotInfo) SlotToEpoch(slot uint64) uint64 {
	if si.SlotsPerEpoch == 0 {
		panic("integer divide by zero")
	}
	val, _ := bits.Div64(0, slot, si.SlotsPerEpoch)
	return val
}

// SlotInEpoch returns the position of slot within its epoch.
func (si *SlotInfo) SlotInEpoch(slot uint64) uint64 {
	return slot % si.SlotsPerEpoch
}

// SlotOfEpochStart returns the first slot of the given epoch.
func (si *SlotInfo) SlotOfEpochStart(epoch uint64) (uint64, error) {
	overflows, slot := bits.Mul64(si.SlotsPerEpoch, epoch)
	if overflows > 0 {
		return slot, errors.New("start slot calculation overflows: multiplication overflows")
	}
	return slot, nil
}

// SlotOfEpochEnd returns the last slot of the given epoch.
func (si *SlotInfo) SlotOfEpochEnd(epoch uint64) (uint64, error) {
	if epoch == math.MaxUint64 {
		return 0, errors.New("start slot calculation overflows")
	}
	slot, err := si.SlotOfEpochStart(epoch + 1)
	if err != nil {
		return 0, err
	}
	return slot - 1, nil
}

// Copy returns a copy of si.
func (si *SlotInfo) Copy() *SlotInfo {
	if si == nil {
		return nil
	}
	return &SlotInfo{
		GenesisTime:    si.GenesisTime,
		SecondsPerSlot: si.SecondsPerSlot,
		SlotsPerEpoch:  si.SlotsPerEpoch,
	}
}
