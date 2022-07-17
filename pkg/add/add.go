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
	"github.com/norbjd/pasuman/pkg/config"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/encrypt"
	"github.com/norbjd/pasuman/pkg/util"
)

func Add(masterPassword string, e data.Entry) (string, error) {
	var d data.Data

	if err := d.FromFile(config.PasumanDataFile); err != nil {
		return "", err
	}

	var err error

	if e.UniqueID == "" {
		if e.UniqueID, err = util.NewUUIDV4(); err != nil {
			return "", err
		}
	}

	if e.ID, err = encrypt.Encrypt(masterPassword, e.ID); err != nil {
		return "", err
	}

	if e.Password, err = encrypt.Encrypt(masterPassword, e.Password); err != nil {
		return "", err
	}

	d.Entries = append(d.Entries, e)

	return e.UniqueID, d.ToFile(config.PasumanDataFile)
}
