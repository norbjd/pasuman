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
	"os"
	"testing"

	"github.com/norbjd/pasuman/internal/pkg/pasumantest"
	"github.com/norbjd/pasuman/pkg/config"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/stretchr/testify/require"
)

func TestListProfiles(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		args   []string
		setup  func() error
		output string
	}{
		{
			args:   []string{"list-profiles"},
			output: "default\n",
		},
		{
			args: []string{"list-profiles"},
			setup: func() error {
				profile := "otherprofile"

				config.Init(profile, true)

				return pasumantest.InitProfile(profile)
			},
			output: "default\notherprofile\n",
		},
		{
			args: []string{"list-profiles"},
			setup: func() error {
				profile := "anotherprofile"

				config.Init(profile, true)

				return pasumantest.InitProfile(profile)
			},
			output: "anotherprofile\ndefault\notherprofile\n",
		},
	}

	for _, tt := range tests {
		if tt.setup != nil {
			require.NoError(t, tt.setup())
		}

		out, err := pasumantest.ExecuteCommand(RootCmd, tt.args...)
		require.NoError(t, err)

		require.Equal(t, tt.output, out)

		pasumantest.Teardown(t, RootCmd)
	}
}
