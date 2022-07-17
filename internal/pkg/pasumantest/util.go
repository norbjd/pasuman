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

package pasumantest

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/alexedwards/argon2id"
	"github.com/norbjd/pasuman/pkg/config"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/stretchr/testify/require"
)

const (
	defaultDirMode  = os.FileMode(0o700)
	defaultFileMode = os.FileMode(0o600)

	TestMasterPassword = "pass"
)

func Init(t *testing.T, profile string) string {
	t.Helper()

	var err error
	tempDir, err := os.MkdirTemp(os.TempDir(), t.Name())
	require.NoError(t, err)

	config.ConfigDir = tempDir + string(os.PathSeparator) + ".config"
	err = os.Mkdir(config.ConfigDir, defaultDirMode)
	require.NoError(t, err)

	dataDir := tempDir + string(os.PathSeparator) + "data"
	err = os.Mkdir(dataDir, defaultDirMode)
	require.NoError(t, err)

	err = os.Mkdir(config.ConfigDir+string(os.PathSeparator)+"pasuman", defaultDirMode)
	require.NoError(t, err)

	err = os.WriteFile(config.ConfigDir+string(os.PathSeparator)+"pasuman"+string(os.PathSeparator)+"config.json",
		[]byte(`{"data_directory": "`+dataDir+`"}`), defaultFileMode)
	require.NoError(t, err)

	config.Init(profile, true)
	require.NoError(t, InitProfile(profile))

	return tempDir
}

func InitProfile(profile string) error {
	masterPasswordHash, err := argon2id.CreateHash(TestMasterPassword, argon2id.DefaultParams)
	if err != nil {
		return err
	}

	err = os.Setenv("PASUMAN_MASTER_PASSWORD", TestMasterPassword)
	if err != nil {
		return err
	}

	return os.WriteFile(config.PasumanDataFile, []byte(`{"master_password":"`+masterPasswordHash+`"}`), 0)
}

func ExecuteCommand(c *cobra.Command, args ...string) (string, error) {
	buf := new(bytes.Buffer)

	c.SetOut(buf)
	c.SetErr(buf)
	c.SetArgs(args)

	err := c.Execute()

	return buf.String(), err
}

// Teardown - reset flags to default value
// and run all "PostRun" functions.
// For some reason, `f.Value.Set(f.DefValue)` does not work
// if the value is of type slice and default value is nil
// because if so, f.DefValue is the string "[]".
func Teardown(t *testing.T, c *cobra.Command) {
	t.Helper()

	c.Flags().VisitAll(func(f *pflag.Flag) {
		var err error

		if strings.HasSuffix(f.Value.Type(), "Slice") {
			err = f.Value.Set("")
		} else {
			err = f.Value.Set(f.DefValue)
		}

		require.NoError(t, err)
	})

	for _, subCommand := range c.Commands() {
		subCommand.Flags().VisitAll(func(f *pflag.Flag) {
			var err error

			if strings.HasSuffix(f.Value.Type(), "Slice") {
				err = f.Value.Set("")
			} else {
				err = f.Value.Set(f.DefValue)
			}

			require.NoError(t, err)
		})
	}

	if c.PostRun != nil {
		c.PostRun(c, nil)
	}

	if c.PostRunE != nil {
		err := c.PostRunE(c, nil)
		require.NoError(t, err)
	}

	if c.PersistentPostRun != nil {
		c.PersistentPostRun(c, nil)
	}

	if c.PersistentPostRunE != nil {
		err := c.PersistentPostRunE(c, nil)
		require.NoError(t, err)
	}
}
