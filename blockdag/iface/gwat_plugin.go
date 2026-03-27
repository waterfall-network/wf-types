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

// NodeConfig holds the configuration required to initialize a gwat plugin.
type NodeConfig struct {
	// DataDir is the root data directory.
	DataDir string
	// NetworkID is the Ethereum network identifier.
	NetworkID uint64
	// HTTPHost and HTTPPort configure the HTTP-RPC server.
	HTTPHost string
	HTTPPort int
	// IPCPath is the path to the IPC socket.
	IPCPath string
	// Extra holds any additional plugin-specific config as raw bytes.
	Extra []byte
}

// GwatPlugin is the lifecycle interface for the gwat execution plugin.
// On Linux, gwat exports this symbol from a .so loaded via plugin.Open.
// On other platforms, it is satisfied by a direct import (static linking).
var _ GwatPlugin = (*gwatPluginNoop)(nil) // compile-time interface check

// GwatPlugin defines the interface that gwat exposes to wf-engine.
type GwatPlugin interface {
	// Init initializes the plugin with the given configuration.
	Init(config *NodeConfig) error

	// Start starts the plugin's background services.
	Start() error

	// Stop stops the plugin's background services.
	Stop() error

	// BlockChain returns the blockchain backend needed by the dag package.
	BlockChain() BlockChain

	// ValidatorChain returns the blockchain backend needed by validator and
	// validatorsync packages. gwat may return the same underlying object as
	// BlockChain() if it satisfies both interfaces.
	ValidatorChain() ValidatorChain

	// TxPool returns an abstracted view of the transaction pool.
	TxPool() TxPool
}

// gwatPluginNoop is a no-op placeholder used only for the compile-time check.
type gwatPluginNoop struct{}

func (g *gwatPluginNoop) Init(_ *NodeConfig) error       { return nil }
func (g *gwatPluginNoop) Start() error                   { return nil }
func (g *gwatPluginNoop) Stop() error                    { return nil }
func (g *gwatPluginNoop) BlockChain() BlockChain         { return nil }
func (g *gwatPluginNoop) ValidatorChain() ValidatorChain { return nil }
func (g *gwatPluginNoop) TxPool() TxPool                 { return nil }
