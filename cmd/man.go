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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var manCmdPath string

func manCmdInit() {
	manCmd.Flags().StringVar(&manCmdPath, "path", "", "Output path")

	if err := manCmd.MarkFlagRequired("path"); err != nil {
		log.Fatal(err)
	}
}

var manCmd = &cobra.Command{
	Use:   "man",
	Short: "Generate man pages",
	RunE: func(cmd *cobra.Command, args []string) error {
		header := &doc.GenManHeader{
			Title:   "pasuman",
			Section: "1",
		}

		err := doc.GenManTree(RootCmd, header, manCmdPath)
		if err != nil {
			return err
		}

		return nil
	},
}
