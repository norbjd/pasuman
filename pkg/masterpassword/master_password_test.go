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

package masterpassword

import (
	"os"
	"testing"

	"github.com/norbjd/pasuman/internal/pkg/pasumantest"
	"github.com/norbjd/pasuman/pkg/add"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/get"
	"github.com/norbjd/pasuman/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestIsSet(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	isSet, err := IsSet()
	require.NoError(t, err)
	require.True(t, isSet)
}

func TestIsCorrect(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	isCorrect, err := IsCorrect(pasumantest.TestMasterPassword)
	require.NoError(t, err)
	require.True(t, isCorrect)

	isCorrect, err = IsCorrect("")
	require.NoError(t, err)
	require.False(t, isCorrect)

	isCorrect, err = IsCorrect(pasumantest.TestMasterPassword + "a")
	require.NoError(t, err)
	require.False(t, isCorrect)
}

func TestSetMasterPassword(t *testing.T) {
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

	type args struct {
		oldMasterPassword string
		newMasterPassword string
	}

	tests := []struct {
		args                  args
		currentMasterPassword string
		wantErr               error
	}{
		{
			args: args{
				oldMasterPassword: "",
				newMasterPassword: pasumantest.TestMasterPassword,
			},
			currentMasterPassword: "",
			wantErr:               nil,
		},
		{
			args: args{
				oldMasterPassword: pasumantest.TestMasterPassword,
				newMasterPassword: "",
			},
			currentMasterPassword: pasumantest.TestMasterPassword,
			wantErr:               util.ErrMasterPasswordMustNotBeEmpty,
		},
		{
			args: args{
				oldMasterPassword: pasumantest.TestMasterPassword,
				newMasterPassword: "newMasterPass",
			},
			currentMasterPassword: pasumantest.TestMasterPassword,
			wantErr:               nil,
		},
		{
			args: args{
				oldMasterPassword: "newMasterPass",
				newMasterPassword: pasumantest.TestMasterPassword,
			},
			currentMasterPassword: "newMasterPass",
			wantErr:               nil,
		},
	}

	for _, tt := range tests {
		err := SetMasterPassword(tt.args.oldMasterPassword, tt.args.newMasterPassword)

		require.ErrorIs(t, tt.wantErr, err)

		if tt.wantErr == nil {
			_, err := get.Sensitive(tt.args.oldMasterPassword, "id1")
			require.Error(t, err)

			gotUniqueID1, err := get.Sensitive(tt.args.newMasterPassword, "id1")
			require.NoError(t, err)

			require.Equal(t, "myId", gotUniqueID1.ID)
			require.Equal(t, "p4$$w0rd!", gotUniqueID1.Password)

			_, err = get.Sensitive(tt.args.oldMasterPassword, "id2")
			require.Error(t, err)

			gotUniqueID2, err := get.Sensitive(tt.args.newMasterPassword, "id2")
			require.NoError(t, err)

			require.Equal(t, "otherId", gotUniqueID2.ID)
			require.Equal(t, "t0ps3cr3t!", gotUniqueID2.Password)
		} else {
			gotUniqueID1, err := get.Sensitive(tt.currentMasterPassword, "id1")
			require.NoError(t, err)

			require.Equal(t, "myId", gotUniqueID1.ID)
			require.Equal(t, "p4$$w0rd!", gotUniqueID1.Password)

			gotUniqueID2, err := get.Sensitive(tt.currentMasterPassword, "id2")
			require.NoError(t, err)

			require.Equal(t, "otherId", gotUniqueID2.ID)
			require.Equal(t, "t0ps3cr3t!", gotUniqueID2.Password)
		}
	}
}
