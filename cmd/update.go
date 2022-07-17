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
	"strings"

	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/get"
	"github.com/norbjd/pasuman/pkg/masterpassword"
	"github.com/norbjd/pasuman/pkg/update"
	"github.com/norbjd/pasuman/pkg/util"
	"github.com/spf13/cobra"
)

var (
	updateCmdUniqueID    string
	updateCmdDescription string
	updateCmdTags        []string
	updateCmdSite        string
	updateCmdID          string
	updateCmdPassword    string
)

func updateCmdInit() {
	updateCmd.Flags().StringVar(&updateCmdUniqueID, "unique-id", "", "Unique ID")
	updateCmd.Flags().StringVar(&updateCmdDescription, "description", "", "Description")
	updateCmd.Flags().StringSliceVar(&updateCmdTags, "tags", []string{}, "Tags")
	updateCmd.Flags().StringVar(&updateCmdSite, "site", "", "Site")
	updateCmd.Flags().StringVar(&updateCmdID, "id", "", "ID")
	updateCmd.Flags().StringVar(&updateCmdPassword, "password", "", "Password")
}

var updateCmd = &cobra.Command{
	Use:               "update",
	Short:             "Update an entry",
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

		if _, err := get.NotSensitive(uniqueID); errors.Is(err, get.ErrNotFound) {
			return update.ErrNotFound
		}

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

		if updateCmdDescription == "" && len(updateCmdTags) == 0 && updateCmdSite == "" &&
			updateCmdID == "" && updateCmdPassword == "" {
			cmdPrintln(cmd, "INFO: Leave field empty if you don't want to update it")

			if updateCmdUniqueID == "" {
				cmdPrintf(cmd, "Enter unique id: ")
				updateCmdUniqueID, err = util.ReadLine()
				if err != nil {
					return err
				}
				updateCmdUniqueID = strings.TrimSpace(updateCmdUniqueID)
			}
			if updateCmdDescription == "" {
				cmdPrintf(cmd, "Enter description: ")
				updateCmdDescription, err = util.ReadLine()
				if err != nil {
					return err
				}
				updateCmdDescription = strings.TrimSpace(updateCmdDescription)
			}
			if len(updateCmdTags) == 0 {
				cmdPrintf(cmd, "Enter tags (comma-separated): ")
				tagsCommaSeparated, err := util.ReadLine()
				if err != nil {
					return err
				}
				tagsCommaSeparatedTrimmed := strings.TrimSpace(tagsCommaSeparated)
				if tagsCommaSeparatedTrimmed != "" {
					updateCmdTags = strings.Split(tagsCommaSeparatedTrimmed, ",")
				}
			}
			if updateCmdSite == "" {
				cmdPrintf(cmd, "Enter site: ")
				updateCmdSite, err = util.ReadLine()
				if err != nil {
					return err
				}
				updateCmdSite = strings.TrimSpace(updateCmdSite)
			}
			if updateCmdID == "" {
				cmdPrintf(cmd, "Enter id: ")
				updateCmdID, err = util.ReadLine()
				if err != nil {
					return err
				}
				updateCmdID = strings.TrimSpace(updateCmdID)
			}
			if updateCmdPassword == "" {
				cmdPrintf(cmd, "Enter password: ")
				updateCmdPassword, err = util.ReadPassword()
				if err != nil {
					return err
				}

				cmdPrintln(cmd, "✔")
			}
		}

		if err := update.Update(
			masterPassword,
			uniqueID,
			data.Entry{
				UniqueID:    updateCmdUniqueID,
				Description: updateCmdDescription,
				Tags:        updateCmdTags,
				Site:        updateCmdSite,
				ID:          updateCmdID,
				Password:    updateCmdPassword,
			},
		); err != nil {
			return err
		}

		if updateCmdUniqueID != "" {
			uniqueID = updateCmdUniqueID
		}
		cmdPrintf(cmd, "Updated: %s\n", uniqueID)

		return nil
	},
}
