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

package util

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenderTable(t *testing.T) {
	t.Parallel()

	headerColumns := []string{"col1", "col2", "col333333333333333333"}

	lines := [][]string{
		{"l1", "l2", "l3"},
		{"loooooooooooooongline", "l4", "l5"},
	}

	var b bytes.Buffer

	RenderTable(&b, headerColumns, lines)

	table := "" +
		"col1\t\t\tcol2\tcol333333333333333333\t\n" +
		"----\t\t\t----\t---------------------\t\n" +
		"l1\t\t\tl2\tl3\t\t\t\n" +
		"loooooooooooooongline\tl4\tl5\t\t\t\n"

	require.Equal(t, table, b.String())
}
