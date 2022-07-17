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

package main

import (
	"errors"
	"log"
	"os"

	"github.com/norbjd/pasuman/cmd"
	"github.com/spf13/cobra"
)

func main() {
	errExecute := cmd.RootCmd.Execute()
	if errors.Is(errExecute, cmd.ErrFileLocked) {
		os.Exit(1)
	}

	if cmd.RootCmd.Short == "remove-lock" {
		return
	}

	if cmd.RootCmd.Short == cobra.ShellCompRequestCmd || cmd.RootCmd.Short == cobra.ShellCompNoDescRequestCmd {
		return
	}

	errRemoveLock := cmd.RemoveLock()
	if errRemoveLock != nil {
		log.Fatal(errRemoveLock)
	}

	if errExecute != nil {
		os.Exit(1)
	}
}
