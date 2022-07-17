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

package add

import (
	"os"
	"regexp"
	"testing"

	"github.com/norbjd/pasuman/internal/pkg/pasumantest"
	"github.com/norbjd/pasuman/pkg/config"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/encrypt"
	"github.com/norbjd/pasuman/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestAdd(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	type args struct {
		masterPassword string
		e              data.Entry
	}

	tests := []struct {
		args    args
		want    string
		check   func(s string)
		wantErr error
	}{
		{
			args: args{
				masterPassword: "",
				e:              data.Entry{},
			},
			wantErr: util.ErrMasterPasswordMustNotBeEmpty,
		},
		{
			args: args{
				masterPassword: "pass",
				e:              data.Entry{},
			},
			check: func(s string) {
				regexUUID := regexp.MustCompile("^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-" +
					"[89ab][0-9a-f]{3}-[0-9a-f]{12}$")
				require.Regexp(t, regexUUID, s)
			},
		},
		{
			args: args{
				masterPassword: "pass",
				e: data.Entry{
					UniqueID: "my-unique-id",
				},
			},
			want: "my-unique-id",
		},
		{
			args: args{
				masterPassword: "pass",
				e: data.Entry{
					UniqueID: "new-unique-id",
					ID:       "my-id",
					Password: "p4$$w0rd!",
				},
			},
			want: "new-unique-id",
			check: func(s string) {
				var d data.Data

				err := d.FromFile(config.PasumanDataFile)
				require.NoError(t, err)

				for _, entry := range d.Entries {
					if entry.UniqueID == "new-unique-id" {
						id, err := encrypt.Decrypt("pass", entry.ID)
						require.NoError(t, err)
						require.Equal(t, "my-id", id)

						password, err := encrypt.Decrypt("pass", entry.Password)
						require.NoError(t, err)
						require.Equal(t, "p4$$w0rd!", password)

						return
					}
				}

				t.Fail()
			},
		},
	}

	for _, tt := range tests {
		got, err := Add(tt.args.masterPassword, tt.args.e)

		if err != nil {
			require.ErrorIs(t, err, tt.wantErr)
		} else {
			if tt.want != "" {
				require.Equal(t, tt.want, got)
			}
			if tt.check != nil {
				tt.check(got)
			}
		}
	}
}
