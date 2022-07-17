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
	"github.com/norbjd/pasuman/pkg/get"
	"github.com/norbjd/pasuman/pkg/update"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		args    []string
		setup   func() error
		check   func()
		output  string
		wantErr error
	}{
		{
			args:    []string{"update", "does-not-exist"},
			wantErr: update.ErrNotFound,
		},
		{
			args: []string{"update", "id1", "--description=New desc"},
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
			check: func() {
				t.Helper()

				entry, err := get.Sensitive(pasumantest.TestMasterPassword, "id1")
				require.NoError(t, err)

				require.Equal(t, data.Entry{
					UniqueID:    "id1",
					Description: "New desc",
					Tags:        []string{"tag1", "tag2"},
					Site:        "https://mysupersite.pasuman",
					ID:          "myId",
					Password:    "p4$$w0rd!",
				}, entry)
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"Updated: id1\n",
		},
		{
			args: []string{"update", "id1", "--tags=tag2,tag3"},
			check: func() {
				t.Helper()

				entry, err := get.Sensitive(pasumantest.TestMasterPassword, "id1")
				require.NoError(t, err)

				require.Equal(t, data.Entry{
					UniqueID:    "id1",
					Description: "New desc",
					Tags:        []string{"tag2", "tag3"},
					Site:        "https://mysupersite.pasuman",
					ID:          "myId",
					Password:    "p4$$w0rd!",
				}, entry)
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"Updated: id1\n",
		},
		{
			args: []string{"update", "id1", "--site=https://mynewsite.pasuman"},
			check: func() {
				t.Helper()

				entry, err := get.Sensitive(pasumantest.TestMasterPassword, "id1")
				require.NoError(t, err)

				require.Equal(t, data.Entry{
					UniqueID:    "id1",
					Description: "New desc",
					Tags:        []string{"tag2", "tag3"},
					Site:        "https://mynewsite.pasuman",
					ID:          "myId",
					Password:    "p4$$w0rd!",
				}, entry)
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"Updated: id1\n",
		},
		{
			args: []string{"update", "id1", "--id=newId"},
			check: func() {
				t.Helper()

				entry, err := get.Sensitive(pasumantest.TestMasterPassword, "id1")
				require.NoError(t, err)

				require.Equal(t, data.Entry{
					UniqueID:    "id1",
					Description: "New desc",
					Tags:        []string{"tag2", "tag3"},
					Site:        "https://mynewsite.pasuman",
					ID:          "newId",
					Password:    "p4$$w0rd!",
				}, entry)
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"Updated: id1\n",
		},
		{
			args: []string{"update", "id1", "--password=n€wp4$$w0rd!"},
			check: func() {
				t.Helper()

				entry, err := get.Sensitive(pasumantest.TestMasterPassword, "id1")
				require.NoError(t, err)

				require.Equal(t, data.Entry{
					UniqueID:    "id1",
					Description: "New desc",
					Tags:        []string{"tag2", "tag3"},
					Site:        "https://mynewsite.pasuman",
					ID:          "newId",
					Password:    "n€wp4$$w0rd!",
				}, entry)
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"Updated: id1\n",
		},
		{
			args: []string{"update", "id1", "--unique-id=newId1"},
			check: func() {
				t.Helper()

				_, err := get.Sensitive(pasumantest.TestMasterPassword, "id1")
				require.ErrorIs(t, err, get.ErrNotFound)

				entry, err := get.Sensitive(pasumantest.TestMasterPassword, "newId1")
				require.NoError(t, err)

				require.Equal(t, data.Entry{
					UniqueID:    "newId1",
					Description: "New desc",
					Tags:        []string{"tag2", "tag3"},
					Site:        "https://mynewsite.pasuman",
					ID:          "newId",
					Password:    "n€wp4$$w0rd!",
				}, entry)
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"Updated: newId1\n",
		},
		{
			args: []string{
				"update", "newId1", "--unique-id=brandNewId1", "--description=Brand new desc",
				"--site=https://brandnewsite.pasuman", "--id=brandNewId", "--password=br4ndn€wp4$$w0rd!",
			},
			check: func() {
				t.Helper()

				_, err := get.Sensitive(pasumantest.TestMasterPassword, "newId1")
				require.ErrorIs(t, err, get.ErrNotFound)

				entry, err := get.Sensitive(pasumantest.TestMasterPassword, "brandNewId1")
				require.NoError(t, err)

				require.Equal(t, data.Entry{
					UniqueID:    "brandNewId1",
					Description: "Brand new desc",
					Tags:        []string{"tag2", "tag3"},
					Site:        "https://brandnewsite.pasuman",
					ID:          "brandNewId",
					Password:    "br4ndn€wp4$$w0rd!",
				}, entry)
			},
			output: "" +
				"Enter current master password: ✔\n" +
				"Updated: brandNewId1\n",
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

		if tt.check != nil {
			tt.check()
		}

		pasumantest.Teardown(t, RootCmd)
	}
}
