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

// Package iface defines interfaces used across wf-types modules.
package iface

// LogEntry represents an EVM event log.
type LogEntry struct {
	Address [20]byte
	Topics  [][32]byte
	Data    []byte
}

// StateDB abstracts execution state access for validator and token processors.
// gwat implements this over *state.StateDB.
type StateDB interface {
	GetBalance(addr [20]byte) []byte // *big.Int encoded as bytes
	SetBalance(addr [20]byte, amount []byte)
	AddBalance(addr [20]byte, amount []byte)
	SubBalance(addr [20]byte, amount []byte)

	GetState(addr [20]byte, key [32]byte) [32]byte
	SetState(addr [20]byte, key [32]byte, value [32]byte)

	GetCode(addr [20]byte) []byte
	SetCode(addr [20]byte, code []byte)

	GetNonce(addr [20]byte) uint64
	SetNonce(addr [20]byte, nonce uint64)

	Exist(addr [20]byte) bool
	CreateAccount(addr [20]byte)
	Suicide(addr [20]byte) bool

	AddLog(log *LogEntry)
	Snapshot() int
	RevertToSnapshot(id int)

	ForEachStorage(addr [20]byte, cb func(key, value [32]byte) bool) error
}
