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

package common

import (
	"math"
	"math/big"
)

// Common big integers often used.
var (
	Big1      = big.NewInt(1)
	Big2      = big.NewInt(2)
	Big3      = big.NewInt(3)
	Big0      = big.NewInt(0)
	Big32     = big.NewInt(32)
	Big100    = big.NewInt(100)
	Big256    = big.NewInt(256)
	Big257    = big.NewInt(257)
	BigGwei   = new(big.Int).SetUint64(1_000_000_000)
	BigWat, _ = new(big.Int).SetString("1000000000000000000", 10)
)

// BnCanCastToUint64 returns true if bn is a non-negative integer that fits in uint64.
func BnCanCastToUint64(bn *big.Int) bool {
	return bn != nil && bn.Sign() >= 0 && new(big.Int).SetUint64(math.MaxUint64).Cmp(bn) >= 0
}
