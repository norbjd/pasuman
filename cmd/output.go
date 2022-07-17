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

	"github.com/spf13/cobra"
)

type cmdOutput string

const (
	outputTable cmdOutput = "table"
	outputJSON  cmdOutput = "json"
)

var (
	outputMessageHelp = fmt.Sprintf(`must be "%s" or "%s"`, outputTable, outputJSON)
	errInvalidOutput  = fmt.Errorf("output %s", outputMessageHelp)
)

func (o *cmdOutput) String() string {
	return string(*o)
}

func (o *cmdOutput) Set(v string) error {
	switch v {
	case string(outputTable), string(outputJSON):
		*o = cmdOutput(v)

		return nil
	default:
		return errInvalidOutput
	}
}

func (o *cmdOutput) Type() string {
	return "output"
}

func outputCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return []string{
		string(outputTable),
		string(outputJSON),
	}, cobra.ShellCompDirectiveDefault
}
