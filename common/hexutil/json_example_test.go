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

package hexutil_test

import (
	"encoding/json"
	"fmt"
	"math/big"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/common/hexutil"
)

func ExampleBytes_marshal() {
	b := hexutil.Bytes{0xde, 0xad, 0xbe, 0xef}
	enc, _ := json.Marshal(b)
	fmt.Println(string(enc))
	// Output: "0xdeadbeef"
}

func ExampleBytes_unmarshal() {
	var b hexutil.Bytes
	_ = json.Unmarshal([]byte(`"0xdeadbeef"`), &b)
	fmt.Printf("%x\n", []byte(b))
	// Output: deadbeef
}

func ExampleUint64_marshal() {
	u := hexutil.Uint64(255)
	enc, _ := json.Marshal(u)
	fmt.Println(string(enc))
	// Output: "0xff"
}

func ExampleUint64_unmarshal() {
	var u hexutil.Uint64
	_ = json.Unmarshal([]byte(`"0xff"`), &u)
	fmt.Println(uint64(u))
	// Output: 255
}

func ExampleBig_marshal() {
	b := (*hexutil.Big)(big.NewInt(1024))
	enc, _ := json.Marshal(b)
	fmt.Println(string(enc))
	// Output: "0x400"
}

func ExampleBig_unmarshal() {
	var b hexutil.Big
	_ = json.Unmarshal([]byte(`"0x400"`), &b)
	fmt.Println(b.ToInt().Int64())
	// Output: 1024
}
