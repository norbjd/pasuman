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

package remove

import (
	"os"
	"testing"

	"github.com/norbjd/pasuman/internal/pkg/pasumantest"
	"github.com/norbjd/pasuman/pkg/add"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/stretchr/testify/require"
)

func TestRemove(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	_, err := add.Add(pasumantest.TestMasterPassword, data.Entry{
		UniqueID:    "id1",
		Description: "A desc",
		Tags:        []string{"tag1", "tag2"},
		Site:        "https://mysupersite.pasuman",
		ID:          "myId",
		Password:    "p4$$w0rd!",
	})
	require.NoError(t, err)

	_, err = add.Add(pasumantest.TestMasterPassword, data.Entry{
		UniqueID:    "id2",
		Description: "Another desc",
		Tags:        []string{"tag1", "tag3", "tag4"},
		Site:        "https://anothersite.pasuman",
		ID:          "otherId",
		Password:    "t0ps3cr3t!",
	})
	require.NoError(t, err)

	tests := []struct {
		uniqueID string
		wantErr  error
	}{
		{
			uniqueID: "id3",
			wantErr:  ErrNotFound,
		},
		{
			uniqueID: "id1",
			wantErr:  nil,
		},
		{
			uniqueID: "id1",
			wantErr:  ErrNotFound,
		},
		{
			uniqueID: "id2",
			wantErr:  nil,
		},
		{
			uniqueID: "id2",
			wantErr:  ErrNotFound,
		},
	}

	for _, tt := range tests {
		err := Remove(tt.uniqueID)

		require.ErrorIs(t, tt.wantErr, err)
	}
}
