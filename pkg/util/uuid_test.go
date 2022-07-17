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
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewUUIDV4(t *testing.T) {
	t.Parallel()

	uuids := make(map[string]struct{}, 1000000)

	for i := 1; i <= 1000000; i++ {
		got, err := NewUUIDV4()
		require.NoError(t, err)

		if _, alreadyExists := uuids[got]; alreadyExists {
			t.Logf("Collision found: %s\n", got)
			t.FailNow()
		}

		uuids[got] = struct{}{}

		uuidParsed, err := uuid.Parse(got)
		require.NoError(t, err)
		require.Equal(t, uuid.Version(4), uuidParsed.Version())
	}
}
