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
	"github.com/norbjd/pasuman/pkg/generate"
	"github.com/spf13/cobra"
)

const passwordDefaultLength = 128

var generateCmdLength uint

func generateCmdInit() {
	generateCmd.Flags().UintVar(&generateCmdLength, "length", passwordDefaultLength,
		"Length wanted for the generated password")
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a random password",
	RunE: func(cmd *cobra.Command, args []string) error {
		generated, err := generate.Generate(generateCmdLength)
		if err != nil {
			return err
		}

		cmdPrintln(cmd, generated)

		return nil
	},
}
