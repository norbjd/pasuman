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
	"errors"
	"os"

	"github.com/norbjd/pasuman/pkg/config"
	"github.com/spf13/cobra"
)

var errNoLock = errors.New("no lock")

var removeLockCmd = &cobra.Command{
	Use:   "remove-lock",
	Short: "Remove lock",
	RunE: func(cmd *cobra.Command, args []string) error {
		lockFile := config.PasumanDataFile + ".lock"

		err := os.Remove(lockFile)
		if os.IsNotExist(err) {
			return errNoLock
		}
		if err != nil {
			return err
		}

		cmdPrintln(cmd, "Lock has been removed!")

		return nil
	},
}
