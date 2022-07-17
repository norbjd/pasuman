// This file is part of pasuman (https://github.com/norbjd/pasuman).
//
// pasuman is a command-line password manager.
// Copyright (C) 2022 norbjd
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3 of the License.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package generate

import (
	"strings"
	"testing"

	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		length uint
	}{
		{length: 0},
		{length: 10},
		{length: 128},
		{length: 256},
	}

	for _, tt := range tests {
		got, err := Generate(tt.length)
		require.NoError(t, err)

		require.Len(t, got, int(tt.length))

		for _, r := range got {
			if !strings.ContainsRune(constants.Alphabet, r) {
				t.Fatalf("Password contains character '%s' but should not", string(r))

				break
			}
		}
	}
}
