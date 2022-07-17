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
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		args    []string
		output  string
		wantErr error
	}{
		{
			args: []string{
				"add", "id1", "--description=A desc", "--tags=tag1,tag2",
				"--site=https://mysupersite.pasuman", "--id=myId", "--password=p4$$w0rd!",
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"New entry: id1\n",
		},
		{
			args: []string{
				"add", "id1", "--description=A desc", "--tags=tag1,tag2",
				"--site=https://mysupersite.pasuman", "--id=myId", "--password=p4$$w0rd!",
			},
			wantErr: errUniqueIDAlreadyExists,
		},
		{
			args: []string{
				"add", "id2", "--description=Another desc", "--tags=tag1,tag3,tag4",
				"--site=https://anothersite.pasuman", "--id=otherId", "--password=t0ps3cr3t!",
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"New entry: id2\n",
		},
	}

	for _, tt := range tests {
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
