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

	"github.com/norbjd/pasuman/pkg/list"
	"github.com/norbjd/pasuman/pkg/masterpassword"
	"github.com/spf13/cobra"
)

var listCmd *cobra.Command

var listCmdOutput = outputTable

func listCmdInit() {
	listCmd = &cobra.Command{
		Use:   "list",
		Short: "List entries",
		RunE:  listCmdRunE,
	}

	listCmd.Flags().Var(&listCmdOutput, "output", fmt.Sprintf("Output format: %s", outputMessageHelp))

	if err := listCmd.RegisterFlagCompletionFunc("output", outputCompletion); err != nil {
		log.Fatal(err)
	}
}

func listCmdRunE(cmd *cobra.Command, args []string) error {
	masterPasswordSet, err := masterpassword.IsSet()
	if err != nil {
		return err
	}

	if !masterPasswordSet {
		return errNoMasterPasswordSet
	}

	entries, err := list.List(rootCmdProfile)
	if err != nil {
		return err
	}

	return printEntries(listCmd.OutOrStdout(), entries, listCmdOutput)
}
