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

package list

import (
	"os"
	"testing"

	"github.com/norbjd/pasuman/internal/pkg/pasumantest"
	"github.com/norbjd/pasuman/pkg/add"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/encrypt"
	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
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

	got, err := List(constants.RootCmdDefaultProfile)
	require.NoError(t, err)

	gotUniqueID1 := got[0]
	require.Equal(t, "id1", gotUniqueID1.UniqueID)
	require.Equal(t, "A desc", gotUniqueID1.Description)
	require.ElementsMatch(t, []string{"tag1", "tag2"}, gotUniqueID1.Tags)
	require.Equal(t, "https://mysupersite.pasuman", gotUniqueID1.Site)

	require.NotEqual(t, "myId", gotUniqueID1.ID)
	gotUniqueID1ID, err := encrypt.Decrypt(pasumantest.TestMasterPassword, gotUniqueID1.ID)
	require.NoError(t, err)
	require.Equal(t, "myId", gotUniqueID1ID)

	require.NotEqual(t, "p4$$w0rd!", gotUniqueID1.Password)
	gotUniqueID1Password, err := encrypt.Decrypt(pasumantest.TestMasterPassword, gotUniqueID1.Password)
	require.NoError(t, err)
	require.Equal(t, "p4$$w0rd!", gotUniqueID1Password)

	gotUniqueID2 := got[1]
	require.Equal(t, "id2", gotUniqueID2.UniqueID)
	require.Equal(t, "Another desc", gotUniqueID2.Description)
	require.ElementsMatch(t, []string{"tag1", "tag3", "tag4"}, gotUniqueID2.Tags)
	require.Equal(t, "https://anothersite.pasuman", gotUniqueID2.Site)

	require.NotEqual(t, "otherId", gotUniqueID2.ID)
	gotUniqueID2ID, err := encrypt.Decrypt(pasumantest.TestMasterPassword, gotUniqueID2.ID)
	require.NoError(t, err)
	require.Equal(t, "otherId", gotUniqueID2ID)

	require.NotEqual(t, "t0ps3cr3t!", gotUniqueID2.Password)
	gotUniqueID2Password, err := encrypt.Decrypt(pasumantest.TestMasterPassword, gotUniqueID2.Password)
	require.NoError(t, err)
	require.Equal(t, "t0ps3cr3t!", gotUniqueID2Password)
}

func TestListEmpty(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	got, err := List(constants.RootCmdDefaultProfile)
	require.NoError(t, err)
	require.Len(t, got, 0)
}
