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
	"strings"
	"testing"
)

func TestStorageSizeString(t *testing.T) {
	tests := []struct {
		size    StorageSize
		wantSuf string
	}{
		{500, "B"},
		{1500, "KiB"},
		{2 * 1048576, "MiB"},
		{2 * 1073741824, "GiB"},
		{2 * 1099511627776, "TiB"},
	}
	for _, tt := range tests {
		s := tt.size.String()
		if !strings.HasSuffix(s, tt.wantSuf) {
			t.Errorf("StorageSize(%g).String() = %q, want suffix %q", float64(tt.size), s, tt.wantSuf)
		}
	}
}

func TestStorageSizeTerminalString(t *testing.T) {
	tests := []struct {
		size    StorageSize
		wantSuf string
	}{
		{500, "B"},
		{1500, "KiB"},
		{2 * 1048576, "MiB"},
		{2 * 1073741824, "GiB"},
		{2 * 1099511627776, "TiB"},
	}
	for _, tt := range tests {
		s := tt.size.TerminalString()
		if !strings.HasSuffix(s, tt.wantSuf) {
			t.Errorf("StorageSize(%g).TerminalString() = %q, want suffix %q", float64(tt.size), s, tt.wantSuf)
		}
	}
}
