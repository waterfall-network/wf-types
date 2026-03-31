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
	// Args are the raw CLI arguments forwarded verbatim to gwat's own flag
	// parser (gopkg.in/urfave/cli.v1). The host does not interpret these;
	// gwat handles all flag semantics internally.
	//
	// Example (everything after "--" on the wf-engine command line):
	//   --datadir=.data/beacon-0 --networkid=333777333 --bootnodes=enode://...
	Args []string

	// DevMode is extracted from Args for use by the host-side dag.Backend.IsDevMode()
	// check. It mirrors the --dev flag in Args — the plugin also reads it from Args.
	DevMode bool
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

	// Downloader returns the sync downloader needed by the dag package.
	Downloader() Downloader

	// Creator returns the block creator needed by the dag package.
	Creator() BlockCreator

	// IsDevMode reports whether the node was started in development mode.
	IsDevMode() bool

	// SetDag replaces gwat's internal dag workloop with the provided wf-engine dag.
	// Must be called after Start(). gwat stops its internal dag and routes all
	// coordinator IPC calls to d instead.
	SetDag(d Dag)
}

// gwatPluginNoop is a no-op placeholder used only for the compile-time check.
type gwatPluginNoop struct{}

func (g *gwatPluginNoop) Init(_ *NodeConfig) error       { return nil }
func (g *gwatPluginNoop) Start() error                   { return nil }
func (g *gwatPluginNoop) Stop() error                    { return nil }
func (g *gwatPluginNoop) BlockChain() BlockChain         { return nil }
func (g *gwatPluginNoop) ValidatorChain() ValidatorChain { return nil }
func (g *gwatPluginNoop) TxPool() TxPool                 { return nil }
func (g *gwatPluginNoop) Downloader() Downloader         { return nil }
func (g *gwatPluginNoop) Creator() BlockCreator          { return nil }
func (g *gwatPluginNoop) IsDevMode() bool                { return false }
func (g *gwatPluginNoop) SetDag(_ Dag)                   {}
