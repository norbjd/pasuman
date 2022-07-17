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

package search

import (
	"os"
	"testing"

	"github.com/norbjd/pasuman/internal/pkg/pasumantest"
	"github.com/norbjd/pasuman/pkg/add"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/stretchr/testify/require"
)

func TestSearch(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	entry1 := data.Entry{
		UniqueID:    "id1",
		Description: "A desc",
		Tags:        []string{"tag1", "tag2"},
		Site:        "https://mysupersite.pasuman",
		ID:          "myId",
		Password:    "p4$$w0rd!",
	}
	_, err := add.Add(pasumantest.TestMasterPassword, entry1)
	require.NoError(t, err)

	entry2 := data.Entry{
		UniqueID:    "id2",
		Description: "Another desc",
		Tags:        []string{"tag1", "tag3", "tag4"},
		Site:        "https://anothersite.pasuman",
		ID:          "otherId",
		Password:    "t0ps3cr3t!",
	}
	_, err = add.Add(pasumantest.TestMasterPassword, entry2)
	require.NoError(t, err)

	type args struct {
		term          string
		caseSensitive bool
	}

	tests := []struct {
		args    args
		want    []data.Entry
		wantErr error
	}{
		{
			args: args{
				term:          "does-not-exist",
				caseSensitive: false,
			},
			want:    nil,
			wantErr: nil,
		},
		{
			args: args{
				term:          "id",
				caseSensitive: false,
			},
			want:    []data.Entry{entry1, entry2},
			wantErr: nil,
		},
		{
			args: args{
				term:          "desc",
				caseSensitive: false,
			},
			want:    []data.Entry{entry1, entry2},
			wantErr: nil,
		},
		{
			args: args{
				term:          "A desc",
				caseSensitive: false,
			},
			want:    []data.Entry{entry1},
			wantErr: nil,
		},
		{
			args: args{
				term:          "A DEsc",
				caseSensitive: true,
			},
			want:    nil,
			wantErr: nil,
		},
		{
			args: args{
				term:          "tag2",
				caseSensitive: false,
			},
			want:    []data.Entry{entry1},
			wantErr: nil,
		},
		{
			args: args{
				term:          "myId",
				caseSensitive: false,
			},
			want:    nil,
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		entries, err := Search(constants.RootCmdDefaultProfile, tt.args.term, tt.args.caseSensitive)

		require.ErrorIs(t, tt.wantErr, err)

		// do not compare ID and Password
		for i := range entries {
			entries[i].ID = ""
			entries[i].Password = ""
		}

		for i := range tt.want {
			tt.want[i].ID = ""
			tt.want[i].Password = ""
		}

		require.ElementsMatch(t, tt.want, entries)
	}
}
