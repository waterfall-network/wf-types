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

package iface

import (
	"crypto/ecdsa"
	"math/big"

	blockdagtypes "gitlab.waterfall.network/waterfall/protocol/wf-types/blockdag/types"
	"gitlab.waterfall.network/waterfall/protocol/wf-types/common"
)

// Account represents a keystore account identified by its address.
type Account struct {
	Address common.Address
}

// Wallet abstracts a keystore wallet that can hold one or more accounts.
type Wallet interface {
	Contains(account Account) bool
}

// InternalKeyStore abstracts the underlying keystore implementation.
// In gwat this is satisfied by a wrapper around *accounts/keystore.KeyStore.
// AccountsCache().ScanAccounts() is flattened into ScanAccounts() for simplicity.
type InternalKeyStore interface {
	Accounts() []Account
	Wallets() []Wallet
	ScanAccounts() error
	IsUnlocked(addr common.Address) bool
	GetKey(addr common.Address) (*ecdsa.PrivateKey, error)
	SignTx(a Account, tx blockdagtypes.Transaction, chainID *big.Int) (blockdagtypes.Transaction, error)
	Unlock(account Account, passphrase string) error
}

// KeystoreConfig provides path configuration for the verifier keystore.
type KeystoreConfig interface {
	KeyDir() (string, error)
	PasswordsFilePath() (string, error)
}

// PasswordPrompter reads a password from the user interactively.
// In gwat this is satisfied by gwat/console/prompt.Stdin.
type PasswordPrompter interface {
	PromptPassword(prompt string) (string, error)
}
