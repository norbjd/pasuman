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
	"github.com/norbjd/pasuman/pkg/add"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/remove"
	"github.com/stretchr/testify/require"
)

func TestRemove(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		args    []string
		setup   func() error
		output  string
		wantErr error
	}{
		{
			args:    []string{"remove", "does-not-exist"},
			wantErr: remove.ErrNotFound,
		},
		{
			args: []string{"remove", "id1"},
			setup: func() error {
				_, err := add.Add(pasumantest.TestMasterPassword, data.Entry{
					UniqueID:    "id1",
					Description: "A desc",
					Tags:        []string{"tag1", "tag2"},
					Site:        "https://mysupersite.pasuman",
					ID:          "myId",
					Password:    "p4$$w0rd!",
				})

				return err
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"Removed entry: id1\n",
		},
		{
			args:    []string{"remove", "id1"},
			wantErr: remove.ErrNotFound,
		},
	}

	for _, tt := range tests {
		if tt.setup != nil {
			require.NoError(t, tt.setup())
		}

		out, err := pasumantest.ExecuteCommand(RootCmd, tt.args...)
		if tt.wantErr != nil {
			require.ErrorIs(t, err, tt.wantErr)
		} else {
			require.NoError(t, err)

			require.Equal(t, tt.output, out)
		}

		pasumantest.Teardown(t, RootCmd)
	}
}
