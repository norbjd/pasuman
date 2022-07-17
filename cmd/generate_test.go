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
	"strings"
	"testing"

	"github.com/norbjd/pasuman/internal/pkg/pasumantest"
	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/stretchr/testify/require"
)

func TestGenerate(t *testing.T) {
	tempDir := pasumantest.Init(t, constants.RootCmdDefaultProfile)
	defer os.RemoveAll(tempDir)

	tests := []struct {
		args           []string
		passwordLength int
		wantErrString  string
	}{
		{
			args:           []string{"generate"},
			passwordLength: passwordDefaultLength,
		},
		{
			args:           []string{"generate", "--length=0"},
			passwordLength: 0,
		},
		{
			args:           []string{"generate", "--length=1"},
			passwordLength: 1,
		},
		{
			args:           []string{"generate", "--length=256"},
			passwordLength: 256,
		},
		{
			args: []string{"generate", "--length=-1"},
			wantErrString: `invalid argument "-1" for "--length" flag: ` +
				`strconv.ParseUint: parsing "-1": invalid syntax`,
		},
	}

	for _, tt := range tests {
		out, err := pasumantest.ExecuteCommand(RootCmd, tt.args...)
		if tt.wantErrString != "" {
			require.EqualError(t, err, tt.wantErrString)
		} else {
			require.NoError(t, err)

			generatedPassword := strings.TrimSpace(out)
			require.Equal(t, tt.passwordLength, len(generatedPassword))

			for _, r := range generatedPassword {
				if !strings.ContainsRune(constants.Alphabet, r) {
					t.Fatalf("Password contains character '%s' but should not", string(r))

					break
				}
			}
		}

		pasumantest.Teardown(t, RootCmd)
	}
}
