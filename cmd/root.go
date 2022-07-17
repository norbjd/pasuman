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
	"fmt"
	"log"
	"os"

	"github.com/norbjd/pasuman/pkg/config"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/spf13/cobra"
)

var ErrFileLocked = errors.New("file is locked")

var rootCmdProfile string

// nolint: gochecknoinits
func init() {
	helpFunc := RootCmd.HelpFunc()
	RootCmd.SetHelpFunc(func(cmd *cobra.Command, args []string) {
		if _, err := cmd.OutOrStdout().Write([]byte(Version() + "\n")); err != nil {
			log.Fatal(err)
		}

		helpFunc(cmd, args)
	})

	RootCmd.PersistentFlags().StringVar(&rootCmdProfile, "profile", constants.RootCmdDefaultProfile, "Profile")

	addCmdInit()
	RootCmd.AddCommand(addCmd)
	generateCmdInit()
	RootCmd.AddCommand(generateCmd)
	getCmdInit()
	RootCmd.AddCommand(getCmd)
	listCmdInit()
	RootCmd.AddCommand(listCmd)
	RootCmd.AddCommand(listProfilesCmd)
	manCmdInit()
	RootCmd.AddCommand(manCmd)
	RootCmd.AddCommand(masterPasswordCmd)
	RootCmd.AddCommand(removeCmd)
	RootCmd.AddCommand(removeLockCmd)
	searchCmdInit()
	RootCmd.AddCommand(searchCmd)
	updateCmdInit()
	RootCmd.AddCommand(updateCmd)

	RootCmd.SetVersionTemplate(Version())
}

var RootCmd = &cobra.Command{
	Use:     "pasuman",
	Short:   "A command-line password manager",
	Version: "1.0.0",
	PersistentPreRunE: func(cmd *cobra.Command, _ []string) error {
		cmd.SilenceUsage = true

		createDataFile := cmd == addCmd || cmd == masterPasswordCmd
		config.Init(rootCmdProfile, createDataFile)

		if cmd == removeLockCmd ||
			cmd.Name() == cobra.ShellCompRequestCmd || cmd.Name() == cobra.ShellCompNoDescRequestCmd {
			return nil
		}

		lockFile := config.PasumanDataFile + ".lock"

		_, err := os.Stat(lockFile)

		if errors.Is(err, os.ErrNotExist) {
			if _, err := os.Create(lockFile); err != nil {
				return err
			}
		} else {
			return fmt.Errorf("%w: use `pasuman remove-lock` command to remove it "+
				"(only if no other pasuman process is running!)", ErrFileLocked)
		}

		return nil
	},
	// caution: if the RunE fails, PersistentPostRun is not called
	PersistentPostRunE: func(cmd *cobra.Command, _ []string) error {
		if cmd.Name() == "__complete" {
			return nil
		}

		return RemoveLock()
	},
}

// RemoveLock - removes the lock created by the `PersistentPreRunE` of the root command.
// We want to execute this even if the command fails (in the `RunE` part).
// If there is an error during RunE, PersistentPostRunE is not run.
// So as far as I know, we should also call `RemoveLock` manually in the `main.go`.
func RemoveLock() error {
	lockFile := config.PasumanDataFile + ".lock"

	err := os.Remove(lockFile)
	if !errors.Is(err, os.ErrNotExist) && err != nil {
		return err
	}

	return nil
}
