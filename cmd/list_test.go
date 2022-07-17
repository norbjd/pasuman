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
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		args   []string
		format cmdOutput
		setup  func() error
		output string
	}{
		{
			args:   []string{"list"},
			format: outputTable,
			output: "Unique ID\tDescription\tTags\tSite\t\n---------\t-----------\t----\t----\t\n",
		},
		{
			args:   []string{"list", "--output=json"},
			format: outputJSON,
			output: "[]",
		},
		{
			args:   []string{"list"},
			format: outputTable,
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
				"Unique ID\tDescription\tTags\t\tSite\t\t\t\t\n" +
				"---------\t-----------\t----\t\t----\t\t\t\t\n" +
				"id1\t\tA desc\t\ttag1,tag2\thttps://mysupersite.pasuman\t\n",
		},
		{
			args:   []string{"list", "--output=json"},
			format: outputJSON,
			output: `[
				{
					"unique_id": "id1",
					"description": "A desc",
					"tags": ["tag1", "tag2"],
					"site": "https://mysupersite.pasuman"
				}
			]`,
		},
		{
			args:   []string{"list"},
			format: outputTable,
			setup: func() error {
				_, err := add.Add(pasumantest.TestMasterPassword, data.Entry{
					UniqueID:    "id2",
					Description: "Another desc",
					Tags:        []string{"tag1", "tag3", "tag4"},
					Site:        "https://anothersite.pasuman",
					ID:          "otherId",
					Password:    "t0ps3cr3t!",
				})

				return err
			},
			output: "" +
				"Unique ID\tDescription\tTags\t\tSite\t\t\t\t\n" +
				"---------\t-----------\t----\t\t----\t\t\t\t\n" +
				"id1\t\tA desc\t\ttag1,tag2\thttps://mysupersite.pasuman\t\n" +
				"id2\t\tAnother desc\ttag1,tag3,tag4\thttps://anothersite.pasuman\t\n",
		},
		{
			args:   []string{"list", "--output=json"},
			format: outputJSON,
			output: `[
				{
					"unique_id": "id1",
					"description": "A desc",
					"tags": ["tag1", "tag2"],
					"site": "https://mysupersite.pasuman"
				},
				{
					"unique_id": "id2",
					"description": "Another desc",
					"tags": ["tag1", "tag3", "tag4"],
					"site": "https://anothersite.pasuman"
				}
			]`,
		},
	}

	for _, tt := range tests {
		if tt.setup != nil {
			require.NoError(t, tt.setup())
		}

		out, err := pasumantest.ExecuteCommand(RootCmd, tt.args...)
		require.NoError(t, err)

		if tt.format == outputJSON {
			require.JSONEq(t, tt.output, out)
		} else {
			require.Equal(t, tt.output, out)
		}

		pasumantest.Teardown(t, RootCmd)
	}
}
