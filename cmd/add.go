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

	"github.com/norbjd/pasuman/pkg/add"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/get"
	"github.com/norbjd/pasuman/pkg/masterpassword"
	"github.com/norbjd/pasuman/pkg/util"
	"github.com/spf13/cobra"
)

var addCmd *cobra.Command

var (
	addCmdDescription string
	addCmdTags        []string
	addCmdSite        string
	addCmdID          string
	addCmdPassword    string
)

var (
	errEmptyPasswordNotAllowed = errors.New("empty password not allowed")
	errUniqueIDAlreadyExists   = errors.New("unique id already exists")
)

func addCmdInit() {
	addCmd = &cobra.Command{
		Use:   "add [unique id]",
		Short: "Add an entry",
		Args:  cobra.MaximumNArgs(1),
		RunE:  addCmdRunE,
	}

	addCmd.Flags().StringVar(&addCmdDescription, "description", "", "Description")
	addCmd.Flags().StringSliceVar(&addCmdTags, "tags", nil, "Tags")
	addCmd.Flags().StringVar(&addCmdSite, "site", "", "Site")
	addCmd.Flags().StringVar(&addCmdID, "id", "", "ID")
	addCmd.Flags().StringVar(&addCmdPassword, "password", "", "Password")
}

// nolint: funlen, gocognit
func addCmdRunE(cmd *cobra.Command, args []string) error {
	masterPasswordSet, err := masterpassword.IsSet()
	if err != nil {
		return err
	}

	if !masterPasswordSet {
		return errNoMasterPasswordSet
	}

	var uniqueID string

	if len(args) == 1 {
		uniqueID = strings.TrimSpace(args[0])

		if _, err := get.NotSensitive(uniqueID); !errors.Is(err, get.ErrNotFound) {
			return errUniqueIDAlreadyExists
		}
	}

	cmdPrintf(addCmd, "Enter current master password: ")

	masterPassword, err := util.ReadPassword()
	if err != nil {
		return err
	}

	correct, err := masterpassword.IsCorrect(masterPassword)
	if err != nil {
		return err
	}

	if !correct {
		cmdPrintln(addCmd, "✘")

		return util.ErrMasterPasswordIncorrect
	}

	cmdPrintln(addCmd, "✔")

	interactive := addCmdDescription == "" && len(addCmdTags) == 0 && addCmdSite == "" &&
		addCmdID == "" && addCmdPassword == ""

	// nolint: nestif
	if interactive {
		if uniqueID == "" {
			cmdPrintf(addCmd, "Enter unique id (leave empty to generate a random one): ")

			uniqueID, err = util.ReadLine()
			if err != nil {
				return err
			}

			uniqueID = strings.TrimSpace(uniqueID)

			if _, err := get.NotSensitive(uniqueID); err == nil {
				return errUniqueIDAlreadyExists
			}
		}

		cmdPrintf(addCmd, "Enter description: ")

		addCmdDescription, err = util.ReadLine()
		if err != nil {
			return err
		}

		addCmdDescription = strings.TrimSpace(addCmdDescription)

		cmdPrintf(addCmd, "Enter tags (comma-separated): ")

		tagsCommaSeparated, err := util.ReadLine()
		if err != nil {
			return err
		}

		tagsCommaSeparatedTrimmed := strings.TrimSpace(tagsCommaSeparated)
		if tagsCommaSeparatedTrimmed != "" {
			addCmdTags = strings.Split(tagsCommaSeparatedTrimmed, ",")
		}

		cmdPrintf(addCmd, "Enter site: ")

		addCmdSite, err = util.ReadLine()
		if err != nil {
			return err
		}

		addCmdSite = strings.TrimSpace(addCmdSite)

		cmdPrintf(addCmd, "Enter id: ")

		addCmdID, err = util.ReadLine()
		if err != nil {
			return err
		}

		cmdPrintf(addCmd, "Enter password: ")

		addCmdPassword, err = util.ReadPassword()
		if err != nil {
			return err
		}

		if addCmdPassword == "" {
			cmdPrintln(addCmd, "✘")

			return errEmptyPasswordNotAllowed
		}

		cmdPrintln(addCmd, "✔")
	}

	uniqueID, err = add.Add(
		masterPassword,
		data.Entry{
			UniqueID:    uniqueID,
			Description: addCmdDescription,
			Tags:        addCmdTags,
			Site:        addCmdSite,
			ID:          addCmdID,
			Password:    addCmdPassword,
		},
	)
	if err != nil {
		return err
	}

	cmdPrintf(addCmd, "New entry: %s\n", uniqueID)

	return nil
}
