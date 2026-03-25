// Copyright 2024   Blue Wave Inc.
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
	"gitlab.waterfall.network/waterfall/protocol/wf-types/crypto"
)

// EvtErrorLogSignature is the receipt error log topic.
var EvtErrorLogSignature = common.Hash(crypto.Keccak256Hash([]byte("error")))

// ParsedLog represents a transaction receipt log with parsed data.
type ParsedLog struct {
	Address     common.Address `json:"address"`
	Topics      []common.Hash  `json:"topics"`
	Data        []byte         `json:"data"`
	BlockNumber uint64         `json:"blockNumber"`
	TxHash      common.Hash    `json:"transactionHash"`
	TxIndex     uint           `json:"transactionIndex"`
	BlockHash   common.Hash    `json:"blockHash"`
	Index       uint           `json:"logIndex"`
	Removed     bool           `json:"removed"`
	// fields for parsed data
	ParsedData   interface{} `json:"parsedData"`
	ParsedTopics []string    `json:"parsedTopics"`
}

// ToParsedLog converts a LogEntry to a ParsedLog.
func (l *LogEntry) ToParsedLog() *ParsedLog {
	topics := make([]common.Hash, len(l.Topics))
	for i, t := range l.Topics {
		topics[i] = t
	}
	return &ParsedLog{
		Address:      l.Address,
		Topics:       topics,
		Data:         l.Data,
		ParsedData:   nil,
		ParsedTopics: nil,
	}
}
