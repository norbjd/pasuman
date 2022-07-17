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
	"strings"

	"github.com/norbjd/pasuman/pkg/masterpassword"
	"github.com/norbjd/pasuman/pkg/remove"
	"github.com/norbjd/pasuman/pkg/util"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:               "remove <unique id>",
	Short:             "Remove an entry",
	Args:              cobra.ExactArgs(1),
	ValidArgsFunction: autocomplete,
	RunE: func(cmd *cobra.Command, args []string) error {
		masterPasswordSet, err := masterpassword.IsSet()
		if err != nil {
			return err
		}

		if !masterPasswordSet {
			return errNoMasterPasswordSet
		}

		uniqueID := strings.TrimSpace(args[0])

		cmdStderrPrintf(cmd, "Enter current master password: ")
		masterPassword, err := util.ReadPassword()
		if err != nil {
			return err
		}

		correct, err := masterpassword.IsCorrect(masterPassword)
		if err != nil {
			return err
		}

		if !correct {
			cmdStderrPrintln(cmd, "✘")

			return util.ErrMasterPasswordIncorrect
		}
		cmdStderrPrintln(cmd, "✔")

		if err := remove.Remove(uniqueID); err != nil {
			return err
		}

		cmdPrintf(cmd, "Removed entry: %s\n", uniqueID)

		return nil
	},
}
