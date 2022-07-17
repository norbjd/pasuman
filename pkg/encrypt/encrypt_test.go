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

package encrypt

import (
	"strings"
	"testing"

	"github.com/norbjd/pasuman/pkg/util"
	"github.com/stretchr/testify/require"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		masterPassword  string
		stringToEncrypt string
		wantErr         error
	}{
		{
			masterPassword:  "",
			stringToEncrypt: "p4$$w0rd!",
			wantErr:         util.ErrMasterPasswordMustNotBeEmpty,
		},
		{
			masterPassword:  "pass",
			stringToEncrypt: "",
		},
		{
			masterPassword:  "pass",
			stringToEncrypt: "p4$$w0rd!",
		},
		{
			masterPassword:  "pass",
			stringToEncrypt: "パスワード",
		},
		{
			masterPassword: "pass",
			// 513 bytes (1 character = 1 byte)
			stringToEncrypt: "verylongpassword/verylongpassword/verylongpassword/verylongpassword/" +
				"verylongpassword/verylongpassword/verylongpassword/verylongpassword/verylongpassword/" +
				"verylongpassword/verylongpassword/verylongpassword/verylongpassword/verylongpassword/" +
				"verylongpassword/verylongpassword/verylongpassword/verylongpassword/verylongpassword/" +
				"verylongpassword/verylongpassword/verylongpassword/verylongpassword/verylongpassword/" +
				"verylongpassword/verylongpassword/verylongpassword/verylongpassword/verylongpassword/" +
				"verylongpassword/ver",
		},
		{
			masterPassword: "pass",
			// 513 bytes (1 character = 3 bytes)
			stringToEncrypt: "パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・" +
				"パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・" +
				"パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワード・パスワ",
		},
	}

	for _, tt := range tests {
		got, err := Encrypt(tt.masterPassword, tt.stringToEncrypt)
		if tt.wantErr != nil {
			require.ErrorIs(t, err, tt.wantErr)

			continue
		}

		if err != nil {
			t.Fatal(err)
		}

		// if password length <= 512 bytes, encrypted string have length 810
		// due to padding
		require.GreaterOrEqual(t, len(got), 810)
		require.Len(t, strings.Split(got, encryptedMessageStringSeparator), encryptedMessageSplitStringLength)

		decrypted, err := Decrypt("pass", got)
		require.NoError(t, err)
		require.Equal(t, tt.stringToEncrypt, decrypted)
	}
}
