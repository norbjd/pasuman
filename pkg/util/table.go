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
	"fmt"
	"io"
	"text/tabwriter"
)

const tabWidth = 8

func RenderTable(w io.Writer, headerColumns []string, lines [][]string) {
	tabW := tabwriter.NewWriter(w, 0, tabWidth, 0, '\t', 0)

	for _, col := range headerColumns {
		fmt.Fprintf(tabW, "%s\t", col)
	}

	fmt.Fprintln(tabW)

	for _, col := range headerColumns {
		for range col {
			fmt.Fprintf(tabW, "%s", "-")
		}

		fmt.Fprintf(tabW, "\t")
	}

	fmt.Fprintln(tabW)

	for _, line := range lines {
		for _, col := range line {
			fmt.Fprintf(tabW, "%s\t", col)
		}

		fmt.Fprintln(tabW)
	}

	tabW.Flush()
}
