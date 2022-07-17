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

package update

import (
	"os"
	"testing"

	"github.com/norbjd/pasuman/internal/pkg/pasumantest"
	"github.com/norbjd/pasuman/pkg/add"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/get"
	"github.com/stretchr/testify/require"
)

func TestUpdate(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	entry := data.Entry{
		UniqueID:    "id1",
		Description: "A desc",
		Tags:        []string{"tag1", "tag2"},
		Site:        "https://mysupersite.pasuman",
		ID:          "myId",
		Password:    "p4$$w0rd!",
	}
	_, err := add.Add(pasumantest.TestMasterPassword, entry)
	require.NoError(t, err)

	tests := []struct {
		uniqueID string
		entry    data.Entry
		wantErr  error
	}{
		{
			uniqueID: "id1",
			entry: data.Entry{
				UniqueID:    "newId1",
				Description: "New desc",
				Tags:        []string{"tag2", "tag3"},
				Site:        "https://newsite.pasuman",
				ID:          "newId",
				Password:    "n€wp4$$w0rd!",
			},
		},
		{
			uniqueID: "id3",
			wantErr:  ErrNotFound,
		},
	}

	for _, tt := range tests {
		err := Update(pasumantest.TestMasterPassword, tt.uniqueID, tt.entry)

		require.ErrorIs(t, err, tt.wantErr)

		if tt.wantErr == nil {
			updated, err := get.Sensitive(pasumantest.TestMasterPassword, tt.entry.UniqueID)
			require.NoError(t, err)

			require.Equal(t, tt.entry, updated)
		}
	}
}
