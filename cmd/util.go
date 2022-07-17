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

package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/list"
	"github.com/norbjd/pasuman/pkg/util"
	"github.com/spf13/cobra"
)

func cmdPrintf(cmd *cobra.Command, format string, a ...interface{}) {
	fmt.Fprintf(cmd.OutOrStdout(), format, a...)
}

func cmdStderrPrintf(cmd *cobra.Command, format string, a ...interface{}) {
	fmt.Fprintf(cmd.ErrOrStderr(), format, a...)
}

func cmdPrintln(cmd *cobra.Command, message string) {
	cmdPrintf(cmd, "%s\n", message)
}

func cmdStderrPrintln(cmd *cobra.Command, message string) {
	cmdStderrPrintf(cmd, "%s\n", message)
}

func printEntries(w io.Writer, entries []data.Entry, output cmdOutput) error {
	switch output {
	case outputTable:
		headerColumns := []string{"Unique ID", "Description", "Tags", "Site"}

		var lines [][]string

		for _, entry := range entries {
			columns := []string{entry.UniqueID, entry.Description, strings.Join(entry.Tags, ","), entry.Site}

			lines = append(lines, columns)
		}

		util.RenderTable(w, headerColumns, lines)
	case outputJSON:
		type jsonEntry struct {
			UniqueID    string   `json:"unique_id"`
			Description string   `json:"description"`
			Tags        []string `json:"tags"`
			Site        string   `json:"site"`
		}

		jsonEntries := make([]jsonEntry, len(entries))

		for idx, e := range entries {
			jsonEntries[idx] = jsonEntry{UniqueID: e.UniqueID, Description: e.Description, Tags: e.Tags, Site: e.Site}
		}

		result, err := json.MarshalIndent(jsonEntries, "", "  ")
		if err != nil {
			return err
		}

		fmt.Fprintln(w, string(result))
	default:
		return errInvalidOutput
	}

	return nil
}

func autocomplete(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) != 0 {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	listResult, err := list.List(rootCmdProfile)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	uniqueIDs := make([]string, len(listResult))
	for idx := range listResult {
		uniqueIDs[idx] = listResult[idx].UniqueID
	}

	return uniqueIDs, cobra.ShellCompDirectiveNoFileComp
}
