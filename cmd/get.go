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
	"log"
	"strings"

	"github.com/norbjd/pasuman/pkg/get"
	"github.com/norbjd/pasuman/pkg/masterpassword"
	"github.com/norbjd/pasuman/pkg/util"
	"github.com/spf13/cobra"
)

var (
	getCmdNotSensitive              bool
	getCmdAll                       bool
	getCmdID                        bool
	getCmdPassword                  bool
	getCmdCopyIDPasswordToClipboard bool
	getCmdOutput                    = outputTable
)

func getCmdInit() {
	getCmd.Flags().BoolVar(&getCmdNotSensitive, "not-sensitive", true,
		"Get only not sensitive data")
	getCmd.Flags().BoolVar(&getCmdAll, "all", false, "Get all data")
	getCmd.Flags().BoolVar(&getCmdID, "id", false, "Get id")
	getCmd.Flags().BoolVar(&getCmdPassword, "password", false, "Get password")
	getCmd.Flags().BoolVar(&getCmdCopyIDPasswordToClipboard, "copy-id-password-to-clipboard", false,
		"Do not print id and password, just copy them to clipboard one after the other")
	getCmd.MarkFlagsMutuallyExclusive("not-sensitive", "all", "id", "password", "copy-id-password-to-clipboard")
	getCmd.Flags().Var(&getCmdOutput, "output", fmt.Sprintf("Output format: %s", outputMessageHelp))

	if err := getCmd.RegisterFlagCompletionFunc("output", outputCompletion); err != nil {
		log.Fatal(err)
	}
}

var getCmd = &cobra.Command{
	Use:               "get <unique id>",
	Short:             "Get an entry",
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

		uniqueID := args[0]

		if getCmdCopyIDPasswordToClipboard {
			getCmdAll = true
		}

		if getCmdNotSensitive && !getCmdAll && !getCmdID && !getCmdPassword {
			entry, err := get.NotSensitive(uniqueID)
			if err != nil {
				return err
			}

			switch getCmdOutput {
			case outputTable:
				headerColumns := []string{"Unique ID", "Description", "Tags", "Site"}
				columns := []string{entry.UniqueID, entry.Description, strings.Join(entry.Tags, ","), entry.Site}

				util.RenderTable(cmd.OutOrStdout(), headerColumns, [][]string{columns})
			case outputJSON:
				var jsonEntry struct {
					UniqueID    string   `json:"unique_id"`
					Description string   `json:"description"`
					Tags        []string `json:"tags"`
					Site        string   `json:"site"`
				}
				jsonEntry.UniqueID = entry.UniqueID
				jsonEntry.Description = entry.Description
				jsonEntry.Tags = entry.Tags
				jsonEntry.Site = entry.Site

				result, err := json.MarshalIndent(jsonEntry, "", "  ")
				if err != nil {
					return err
				}
				cmdPrintln(cmd, string(result))
			default:
				return errInvalidOutput
			}

			return nil
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

		entry, err := get.Sensitive(masterPassword, uniqueID)
		if err != nil {
			return err
		}

		if getCmdID {
			cmdPrintln(cmd, entry.ID)

			return nil
		}

		if getCmdPassword {
			cmdPrintln(cmd, entry.Password)

			return nil
		}

		var entryID string
		var entryPassword string

		if getCmdCopyIDPasswordToClipboard {
			entryID = entry.ID
			entry.ID = "XXX*"

			entryPassword = entry.Password
			entry.Password = "XXX*"
		}

		if getCmdAll {
			switch getCmdOutput {
			case outputTable:
				headerColumns := []string{"Unique ID", "Description", "Tags", "Site", "ID", "Password"}
				columns := []string{
					entry.UniqueID, entry.Description, strings.Join(entry.Tags, ","), entry.Site,
					entry.ID, entry.Password,
				}

				util.RenderTable(cmd.OutOrStdout(), headerColumns, [][]string{columns})
			case outputJSON:
				var jsonEntry struct {
					UniqueID    string   `json:"unique_id"`
					Description string   `json:"description"`
					Tags        []string `json:"tags"`
					Site        string   `json:"site"`
					ID          string   `json:"id"`
					Password    string   `json:"password"`
				}
				jsonEntry.UniqueID = entry.UniqueID
				jsonEntry.Description = entry.Description
				jsonEntry.Tags = entry.Tags
				jsonEntry.Site = entry.Site
				jsonEntry.ID = entry.ID
				jsonEntry.Password = entry.Password

				result, err := json.MarshalIndent(jsonEntry, "", "  ")
				if err != nil {
					return err
				}
				cmdPrintln(cmd, string(result))
			default:
				return errInvalidOutput
			}
		}

		if getCmdCopyIDPasswordToClipboard {
			cmdPrintln(cmd, "\n*pasuman uses ANSI OSC 52 to copy to clipboard, it may not work on some terminals\n"+
				"(see https://github.com/ojroques/vim-oscyank/blob/main/README.md)\n")
		}

		if getCmdCopyIDPasswordToClipboard {
			cmdPrintf(cmd, "ID copied to clipboard, press enter to copy password")

			util.CopyToClipboard(entryID)

			_, err := util.ReadLine()
			if err != nil {
				return err
			}

			cmdPrintf(cmd, "Password copied to clipboard, press enter to clear clipboard")

			util.CopyToClipboard(entryPassword)

			_, err = util.ReadLine()
			if err != nil {
				return err
			}

			util.CopyToClipboard("")
		}

		return nil
	},
}
