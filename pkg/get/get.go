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

package get

import (
	"errors"

	"github.com/norbjd/pasuman/pkg/config"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/encrypt"
)

var ErrNotFound = errors.New("entry not found")

func getEntry(uniqueID string) (data.Entry, error) {
	var d data.Data

	if err := d.FromFile(config.PasumanDataFile); err != nil {
		return data.Entry{}, err
	}

	var entry data.Entry

	for _, e := range d.Entries {
		if e.UniqueID == uniqueID {
			entry = e

			break
		}
	}

	if entry.UniqueID == "" {
		return data.Entry{}, ErrNotFound
	}

	return entry, nil
}

func NotSensitive(uniqueID string) (data.Entry, error) {
	entry, err := getEntry(uniqueID)
	if err != nil {
		return data.Entry{}, err
	}

	entry.ID = ""
	entry.Password = ""

	return entry, nil
}

func Sensitive(masterPassword string, uniqueID string) (data.Entry, error) {
	entry, err := getEntry(uniqueID)
	if err != nil {
		return data.Entry{}, err
	}

	entry.ID, err = encrypt.Decrypt(masterPassword, entry.ID)
	if err != nil {
		return data.Entry{}, err
	}

	entry.Password, err = encrypt.Decrypt(masterPassword, entry.Password)
	if err != nil {
		return data.Entry{}, err
	}

	return entry, nil
}
