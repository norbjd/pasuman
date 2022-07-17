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

package util

import (
	"encoding/base64"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyToClipboard(t *testing.T) {
	t.Parallel()

	toCopy := "this should go to clipboard"

	realStdout := os.Stdout
	defer func() { os.Stdout = realStdout }()

	r, fakeStdout, err := os.Pipe()
	require.NoError(t, err)

	os.Stdout = fakeStdout

	CopyToClipboard(toCopy)

	err = fakeStdout.Close()
	require.NoError(t, err)

	fakeStdoutContent, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	toCopyBase64 := base64.StdEncoding.EncodeToString([]byte(toCopy))

	require.Equal(t, "\033]52;c;"+toCopyBase64+"\a", string(fakeStdoutContent))
}
