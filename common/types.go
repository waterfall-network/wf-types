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

package common

import (
	"bytes"
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"sort"
	"strings"

	"gitlab.waterfall.network/waterfall/protocol/wf-types/common/hexutil"
	"golang.org/x/crypto/sha3"
)

// Lengths of hashes and addresses in bytes.
const (
	HashLength    = 32
	AddressLength = 20
)

var (
	hashT    = reflect.TypeOf(Hash{})
	addressT = reflect.TypeOf(Address{})
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BigToHash sets byte representation of b to hash.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash sets byte representation of s to hash.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// TerminalString implements log.TerminalStringer.
func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x..%x", h[:3], h[29:])
}

// String implements fmt.Stringer.
func (h Hash) String() string { return h.Hex() }

// Format implements fmt.Formatter.
func (h Hash) Format(s fmt.State, c rune) {
	hexb := make([]byte, 2+len(h)*2)
	copy(hexb, "0x")
	hex.Encode(hexb[2:], h[:])

	switch c {
	case 'x', 'X':
		if !s.Flag('#') {
			hexb = hexb[2:]
		}
		if c == 'X' {
			hexb = bytes.ToUpper(hexb)
		}
		fallthrough
	case 'v', 's':
		s.Write(hexb)
	case 'q':
		q := []byte{'"'}
		s.Write(q)
		s.Write(hexb)
		s.Write(q)
	case 'd':
		fmt.Fprint(s, ([len(h)]byte)(h))
	default:
		fmt.Fprintf(s, "%%!%c(hash=%x)", c, h)
	}
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON("Hash", input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// Ptr returns a pointer to a copy of h.
func (h Hash) Ptr() *Hash { return &h }

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}
	copy(h[HashLength-len(b):], b)
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != HashLength {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), HashLength)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) { return h[:], nil }

// UnprefixedHash allows marshaling a Hash without 0x prefix.
type UnprefixedHash Hash

// UnmarshalText decodes the hash from hex. The 0x prefix is optional.
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// MarshalText encodes the hash as hex.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

// HashArray represents a slice of Hashes.
type HashArray []Hash

// HashArrayFromBytes decodes the HashArray from byte representation.
func HashArrayFromBytes(data []byte) HashArray {
	count := len(data) / HashLength
	ha := make(HashArray, count)
	for i := 0; i < count; i++ {
		ha[i] = Hash{}
		ha[i].SetBytes(data[i*HashLength : (i+1)*HashLength])
	}
	return ha
}

// ToBytes encodes the HashArray to its byte representation.
func (ha HashArray) ToBytes() []byte {
	res := make([]byte, len(ha)*HashLength)
	for i := 0; i < len(ha); i++ {
		copy(res[i*HashLength:(i+1)*HashLength], ha[i][:])
	}
	return res
}

// Copy creates a copy of HashArray.
func (ha HashArray) Copy() HashArray {
	c := make(HashArray, len(ha))
	copy(c, ha)
	return c
}

// Intersection returns hashes present in both ha and hashArray.
func (ha HashArray) Intersection(hashArray HashArray) HashArray {
	c := make(HashArray, 0)
	m := make(map[Hash]bool)
	for _, item := range ha {
		m[item] = true
	}
	for _, item := range hashArray {
		if _, ok := m[item]; ok {
			c = append(c, item)
		}
	}
	return c
}

// SequenceIntersection returns the first contiguous sequence in ha that matches hashArray.
func (ha HashArray) SequenceIntersection(hashArray HashArray) HashArray {
	c := make(HashArray, 0)
	index := -1
	for _, item := range ha {
		i := hashArray.IndexOf(item)
		if index == -1 && i >= 0 {
			index = i + 1
			c = append(c, item)
			continue
		}
		if index >= 0 && index == i {
			index = i + 1
			c = append(c, item)
			continue
		}
		if index >= 0 && index != i {
			break
		}
	}
	return c
}

// Difference returns hashes in ha that are not in hashArray.
func (ha HashArray) Difference(hashArray HashArray) HashArray {
	mb := make(map[Hash]struct{})
	for _, x := range hashArray {
		mb[x] = struct{}{}
	}
	diff := []Hash{}
	for _, x := range ha {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
}

// Uniq returns a new HashArray with duplicate hashes removed.
func (ha HashArray) Uniq() HashArray {
	c := make(HashArray, 0)
	m := make(map[Hash]struct{})
	for _, item := range ha {
		if _, ok := m[item]; !ok {
			c = append(c, item)
			m[item] = struct{}{}
		}
	}
	return c
}

// Deduplicate removes duplicate hashes from ha in-place.
func (ha *HashArray) Deduplicate() {
	m := make(map[Hash]struct{})
	j := 0
	for _, item := range *ha {
		if _, ok := m[item]; !ok {
			(*ha)[j] = item
			j++
			m[item] = struct{}{}
		}
	}
	*ha = (*ha)[:j]
}

// IsUniq returns true if all hashes in ha are unique.
func (ha HashArray) IsUniq() bool {
	return len(ha) == len(ha.Copy().Uniq())
}

// Concat appends hashes to ha and returns the result.
func (ha HashArray) Concat(hashes HashArray) HashArray {
	ha = append(ha, hashes...)
	return ha
}

// IsEqualTo returns true if ha and hashArray have the same hashes in the same order.
func (ha HashArray) IsEqualTo(hashArray HashArray) bool {
	if len(ha) != len(hashArray) {
		return false
	}
	for i, v := range ha {
		if v != hashArray[i] {
			return false
		}
	}
	return true
}

// Sort sorts ha lexicographically and returns it.
func (ha HashArray) Sort() HashArray {
	strs := make([]string, len(ha))
	for i, hash := range ha {
		strs[i] = hash.String()
	}
	sort.Strings(strs)
	for i, str := range strs {
		ha[i] = HexToHash(str)
	}
	return ha
}

// Has returns true if ha contains hash.
func (ha HashArray) Has(hash Hash) bool {
	for _, h := range ha {
		if h == hash {
			return true
		}
	}
	return false
}

// IndexOf returns the index of hash in ha, or -1 if not found.
func (ha HashArray) IndexOf(hash Hash) int {
	for i, h := range ha {
		if h == hash {
			return i
		}
	}
	return -1
}

// Reverse returns a reversed copy of ha.
func (ha HashArray) Reverse() HashArray {
	cpy := ha.Copy()
	i := 0
	j := len(ha) - 1
	for i < j {
		cpy[i], cpy[j] = cpy[j], cpy[i]
		i++
		j--
	}
	return cpy
}

// Key computes a deterministic hash key from the unique, sorted hashes in ha.
func (ha HashArray) Key() Hash {
	c := ha.Copy()
	buf := c.Uniq().Sort().ToBytes()
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[:])
	key := sha.Sum(nil)
	return BytesToHash(key)
}

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

// BytesToAddress returns Address with value b.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// BigToAddress returns Address with byte values of b.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// HexToAddress returns Address with byte values of s.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// IsHexAddress verifies whether a string can represent a valid hex-encoded Ethereum address.
func IsHexAddress(s string) bool {
	if has0xPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

// Contains reports whether desired is in addresses.
func Contains(addresses []Address, desired Address) (bool, int) {
	for i, a := range addresses {
		if a == desired {
			return true, i
		}
	}
	return false, -1
}

// SortAddresses sorts addresses lexicographically.
func SortAddresses(addresses []Address) []Address {
	res := make([]Address, len(addresses))
	strs := make([]string, len(addresses))
	for i, a := range addresses {
		strs[i] = a.String()
	}
	sort.Strings(strs)
	for i, str := range strs {
		res[i] = HexToAddress(str)
	}
	return res
}

// Bytes gets the byte representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a Address) Hex() string { return string(a.checksumHex()) }

// String implements fmt.Stringer.
func (a Address) String() string { return a.Hex() }

func (a *Address) checksumHex() []byte {
	buf := a.hex()
	sha := sha3.NewLegacyKeccak256()
	sha.Write(buf[2:])
	hash := sha.Sum(nil)
	for i := 2; i < len(buf); i++ {
		hashByte := hash[(i-2)/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if buf[i] > '9' && hashByte > 7 {
			buf[i] -= 32
		}
	}
	return buf[:]
}

func (a Address) hex() []byte {
	var buf [len(a)*2 + 2]byte
	copy(buf[:2], "0x")
	hex.Encode(buf[2:], a[:])
	return buf[:]
}

// Format implements fmt.Formatter.
func (a Address) Format(s fmt.State, c rune) {
	switch c {
	case 'v', 's':
		s.Write(a.checksumHex())
	case 'q':
		q := []byte{'"'}
		s.Write(q)
		s.Write(a.checksumHex())
		s.Write(q)
	case 'x', 'X':
		h := a.hex()
		if !s.Flag('#') {
			h = h[2:]
		}
		if c == 'X' {
			h = bytes.ToUpper(h)
		}
		s.Write(h)
	case 'd':
		fmt.Fprint(s, ([len(a)]byte)(a))
	default:
		fmt.Fprintf(s, "%%!%c(address=%x)", c, a)
	}
}

// SetBytes sets the address to the value of b.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses an address in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses an address in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON("Address", input, a[:])
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) { return a[:], nil }

// UnprefixedAddress allows marshaling an Address without 0x prefix.
type UnprefixedAddress Address

// UnmarshalText decodes the address from hex. The 0x prefix is optional.
func (a *UnprefixedAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAddress", input, a[:])
}

// MarshalText encodes the address as hex.
func (a UnprefixedAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

// MixedcaseAddress retains the original string, which may or may not be correctly checksummed.
type MixedcaseAddress struct {
	addr     Address
	original string
}

// NewMixedcaseAddress creates a MixedcaseAddress from an Address.
func NewMixedcaseAddress(addr Address) MixedcaseAddress {
	return MixedcaseAddress{addr: addr, original: addr.Hex()}
}

// NewMixedcaseAddressFromString creates a MixedcaseAddress from a hex string.
func NewMixedcaseAddressFromString(hexaddr string) (*MixedcaseAddress, error) {
	if !IsHexAddress(hexaddr) {
		return nil, errors.New("invalid address")
	}
	a := FromHex(hexaddr)
	return &MixedcaseAddress{addr: BytesToAddress(a), original: hexaddr}, nil
}

// UnmarshalJSON parses MixedcaseAddress.
func (ma *MixedcaseAddress) UnmarshalJSON(input []byte) error {
	if err := hexutil.UnmarshalFixedJSON("Address", input, ma.addr[:]); err != nil {
		return err
	}
	return json.Unmarshal(input, &ma.original)
}

// MarshalJSON marshals the original value.
func (ma *MixedcaseAddress) MarshalJSON() ([]byte, error) {
	if strings.HasPrefix(ma.original, "0x") || strings.HasPrefix(ma.original, "0X") {
		return json.Marshal(fmt.Sprintf("0x%s", ma.original[2:]))
	}
	return json.Marshal(fmt.Sprintf("0x%s", ma.original))
}

// Address returns the address.
func (ma *MixedcaseAddress) Address() Address { return ma.addr }

// String implements fmt.Stringer.
func (ma *MixedcaseAddress) String() string {
	if ma.ValidChecksum() {
		return fmt.Sprintf("%s [chksum ok]", ma.original)
	}
	return fmt.Sprintf("%s [chksum INVALID]", ma.original)
}

// ValidChecksum returns true if the address has valid EIP-55 checksum.
func (ma *MixedcaseAddress) ValidChecksum() bool {
	return ma.original == ma.addr.Hex()
}

// Original returns the mixed-case input string.
func (ma *MixedcaseAddress) Original() string { return ma.original }
