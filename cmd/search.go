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
	"fmt"
	"log"

	"github.com/norbjd/pasuman/pkg/masterpassword"
	"github.com/norbjd/pasuman/pkg/search"
	"github.com/spf13/cobra"
)

var searchCmd *cobra.Command

var (
	searchCmdCaseSensitive bool
	searchCmdOutput        = outputTable
)

func searchCmdInit() {
	searchCmd = &cobra.Command{
		Use:   "search <term>",
		Short: "Search an entry by a term",
		Args:  cobra.ExactArgs(1),
		RunE:  searchCmdRunE,
	}

	searchCmd.Flags().BoolVar(&searchCmdCaseSensitive, "case-sensitive", false, "Case-sensitive search")
	searchCmd.Flags().Var(&searchCmdOutput, "output", fmt.Sprintf("Output format: %s", outputMessageHelp))

	if err := searchCmd.RegisterFlagCompletionFunc("output", outputCompletion); err != nil {
		log.Fatal(err)
	}
}

func searchCmdRunE(cmd *cobra.Command, args []string) error {
	masterPasswordSet, err := masterpassword.IsSet()
	if err != nil {
		return err
	}

	if !masterPasswordSet {
		return errNoMasterPasswordSet
	}

	term := args[0]

	entries, err := search.Search(rootCmdProfile, term, searchCmdCaseSensitive)
	if err != nil {
		return err
	}

	return printEntries(searchCmd.OutOrStdout(), entries, searchCmdOutput)
}
