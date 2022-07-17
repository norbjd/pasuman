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

package update

import (
	"errors"

	"github.com/norbjd/pasuman/pkg/config"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/encrypt"
)

var ErrNotFound = errors.New("entry not found")

func Update(masterPassword string, uniqueID string, e data.Entry) error {
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

	if e.UniqueID != "" {
		d.Entries[index].UniqueID = e.UniqueID
	}

	if e.Description != "" {
		d.Entries[index].Description = e.Description
	}

	if len(e.Tags) > 0 {
		d.Entries[index].Tags = e.Tags
	}

	if e.Site != "" {
		d.Entries[index].Site = e.Site
	}

	var err error

	if e.ID != "" {
		if d.Entries[index].ID, err = encrypt.Encrypt(masterPassword, e.ID); err != nil {
			return err
		}
	}

	if e.Password != "" {
		if d.Entries[index].Password, err = encrypt.Encrypt(masterPassword, e.Password); err != nil {
			return err
		}
	}

	return d.ToFile(config.PasumanDataFile)
}
