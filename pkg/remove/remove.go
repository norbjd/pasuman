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

package remove

import (
	"errors"

	"github.com/norbjd/pasuman/pkg/config"
	"github.com/norbjd/pasuman/pkg/data"
)

var ErrNotFound = errors.New("entry not found")

func Remove(uniqueID string) error {
	var d data.Data

	if err := d.FromFile(config.PasumanDataFile); err != nil {
		return err
	}

	index := -1

	for idx := range d.Entries {
		if d.Entries[idx].UniqueID == uniqueID {
			index = idx

			break
		}
	}

	if index == -1 {
		return ErrNotFound
	}

	d.Entries = append(d.Entries[:index], d.Entries[index+1:]...)

	return d.ToFile(config.PasumanDataFile)
}
