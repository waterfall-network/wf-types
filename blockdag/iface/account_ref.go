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

// Ref is the minimal interface for an entity identified by an Ethereum address.
// token.Processor and validator.Processor accept Ref as the caller argument.
// It is equivalent to gwat's vm.ContractRef.
type Ref interface {
	Address() [20]byte
}

// AccountRef is a concrete Ref backed by a raw 20-byte address.
// It is equivalent to gwat's vm.AccountRef and is used as the caller
// in token and validator processor operations.
type AccountRef [20]byte

// Address implements Ref.
func (ar AccountRef) Address() [20]byte { return [20]byte(ar) }
