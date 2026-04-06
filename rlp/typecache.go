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

package rlp

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
)

// structField describes one field of a struct for RLP encoding/decoding.
type structField struct {
	index    int  // field index in the struct
	nilOK    bool // rlp:"nil" — pointer may be nil
	optional bool // rlp:"optional" — may be absent at end of list
	tail     bool // rlp:"tail" — receives remaining bytes as []byte
}

var (
	structCacheMu sync.RWMutex
	structCache   = make(map[reflect.Type][]structField)
)

// cachedStructFields returns the RLP-relevant fields of a struct type,
// using a cache to avoid repeated reflection.
func cachedStructFields(t reflect.Type) ([]structField, error) {
	structCacheMu.RLock()
	fields, ok := structCache[t]
	structCacheMu.RUnlock()
	if ok {
		return fields, nil
	}

	fields, err := buildStructFields(t)
	if err != nil {
		return nil, err
	}

	structCacheMu.Lock()
	structCache[t] = fields
	structCacheMu.Unlock()
	return fields, nil
}

// buildStructFields inspects a struct type and returns its field descriptors.
func buildStructFields(t reflect.Type) ([]structField, error) {
	var fields []structField
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)

		// Skip unexported fields
		if f.PkgPath != "" {
			continue
		}

		tag := f.Tag.Get("rlp")
		if tag == "-" {
			continue
		}

		sf := structField{index: i}
		for _, part := range strings.Split(tag, ",") {
			switch strings.TrimSpace(part) {
			case "nil":
				sf.nilOK = true
			case "optional":
				sf.optional = true
			case "tail":
				sf.tail = true
			case "":
				// no tag or empty part — ignore
			default:
				return nil, fmt.Errorf("rlp: unknown struct tag %q on field %s.%s", part, t.Name(), f.Name)
			}
		}

		// Validate nilOK: only valid on pointer fields
		if sf.nilOK && f.Type.Kind() != reflect.Ptr {
			return nil, fmt.Errorf("rlp: field %s.%s has rlp:\"nil\" but is not a pointer", t.Name(), f.Name)
		}

		// Validate tail: only valid on []byte, must be last exported field
		if sf.tail {
			if f.Type.Kind() != reflect.Slice || f.Type.Elem().Kind() != reflect.Uint8 {
				return nil, fmt.Errorf("rlp: field %s.%s has rlp:\"tail\" but is not []byte", t.Name(), f.Name)
			}
			// Ensure no exported fields follow
			for j := i + 1; j < t.NumField(); j++ {
				next := t.Field(j)
				if next.PkgPath == "" && next.Tag.Get("rlp") != "-" {
					return nil, fmt.Errorf("rlp: field %s.%s has rlp:\"tail\" but is not the last field", t.Name(), f.Name)
				}
			}
		}

		fields = append(fields, sf)
	}

	// Validate optional: optional fields must only appear after all required fields
	// (i.e., once optional starts, everything after must also be optional or tail)
	seenOptional := false
	for _, f := range fields {
		if f.optional || f.tail {
			seenOptional = true
		} else if seenOptional {
			return nil, fmt.Errorf("rlp: struct %s has required field after optional field", t.Name())
		}
	}

	return fields, nil
}
