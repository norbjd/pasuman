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

	"github.com/norbjd/pasuman/pkg/masterpassword"
	"github.com/norbjd/pasuman/pkg/util"
	"github.com/spf13/cobra"
)

var (
	errSamePasswordAsBefore = errors.New("same password as before")
	errNewPasswordsMismatch = errors.New("new master passwords mismatch")
)

var masterPasswordCmd = &cobra.Command{
	Use:   "master-password",
	Short: "Set or change master password",
	RunE: func(cmd *cobra.Command, args []string) error {
		masterPasswordSet, err := masterpassword.IsSet()
		if err != nil {
			return err
		}

		var currentMasterPassword string

		if masterPasswordSet {
			cmdStderrPrintf(cmd, "Enter current master password: ")
			currentMasterPassword, err = util.ReadPassword()
			if err != nil {
				return err
			}

			correct, err := masterpassword.IsCorrect(currentMasterPassword)
			if err != nil {
				return err
			}

			if !correct {
				cmdStderrPrintln(cmd, "✘")

				return util.ErrMasterPasswordIncorrect
			}

			cmdStderrPrintln(cmd, "✔")
		}

		cmdPrintf(cmd, "Enter new master password: ")
		newMasterPassword, err := util.ReadPassword()
		if err != nil {
			return err
		}
		if newMasterPassword == "" {
			cmdPrintln(cmd, "✘")

			return util.ErrMasterPasswordMustNotBeEmpty
		}

		if masterPasswordSet {
			sameAsCurrentPassword, err := masterpassword.IsCorrect(newMasterPassword)
			if err != nil {
				return err
			}

			if sameAsCurrentPassword {
				cmdPrintln(cmd, "✘")

				return errSamePasswordAsBefore
			}
		}

		cmdPrintln(cmd, "✔")

		cmdPrintf(cmd, "Enter new master password again: ")
		newMasterPasswordAgain, err := util.ReadPassword()
		if err != nil {
			return err
		}

		if newMasterPassword != newMasterPasswordAgain {
			cmdPrintln(cmd, "✘")

			return errNewPasswordsMismatch
		}
		cmdPrintln(cmd, "✔")

		if err := masterpassword.SetMasterPassword(currentMasterPassword, newMasterPassword); err != nil {
			return err
		}

		cmdPrintln(cmd, "Master password has been set!")

		return nil
	},
}
