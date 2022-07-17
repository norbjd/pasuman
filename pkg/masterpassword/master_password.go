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

package masterpassword

import (
	"fmt"

	"github.com/alexedwards/argon2id"
	"github.com/norbjd/pasuman/pkg/config"
	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/encrypt"
	"github.com/norbjd/pasuman/pkg/util"
)

// nolint: gomnd
var argon2idParams = argon2id.Params{
	SaltLength:  64,
	KeyLength:   32,
	Memory:      256 * 1024,
	Iterations:  16,
	Parallelism: 4,
}

func IsSet() (bool, error) {
	var d data.Data

	if err := d.FromFile(config.PasumanDataFile); err != nil {
		return false, err
	}

	if d.MasterPassword == "" {
		return false, nil
	}

	return true, nil
}

func IsCorrect(masterPassword string) (bool, error) {
	var data data.Data

	if err := data.FromFile(config.PasumanDataFile); err != nil {
		return false, err
	}

	match, _, err := argon2id.CheckHash(masterPassword, data.MasterPassword)
	if err != nil {
		return false, err
	}

	return match, nil
}

// SetMasterPassword - change master password and re-encrypt all encrypted values
// with the new master password.
// nolint: funlen
func SetMasterPassword(oldMasterPassword, newMasterPassword string) error {
	if newMasterPassword == "" {
		return util.ErrMasterPasswordMustNotBeEmpty
	}

	var err error

	if oldMasterPassword != "" {
		oldMasterPasswordCorrect, err := IsCorrect(oldMasterPassword)
		if !oldMasterPasswordCorrect {
			return util.ErrMasterPasswordIncorrect
		}

		if err != nil {
			return err
		}
	}

	var data data.Data

	if err := data.FromFile(config.PasumanDataFile); err != nil {
		return err
	}

	if data.MasterPassword, err = argon2id.CreateHash(newMasterPassword, &argon2idParams); err != nil {
		return err
	}

	if oldMasterPassword != "" {
		fmt.Println("Re-encrypting all entries with new master password, please wait...")

		for idx := range data.Entries {
			encryptedEntryID := data.Entries[idx].ID

			entryID, err := encrypt.Decrypt(oldMasterPassword, encryptedEntryID)
			if err != nil {
				return err
			}

			if data.Entries[idx].ID, err = encrypt.Encrypt(newMasterPassword, entryID); err != nil {
				return err
			}

			encryptedEntryPassword := data.Entries[idx].Password

			entryPassword, err := encrypt.Decrypt(oldMasterPassword, encryptedEntryPassword)
			if err != nil {
				return err
			}

			if data.Entries[idx].Password, err = encrypt.Encrypt(newMasterPassword, entryPassword); err != nil {
				return err
			}
		}

		fmt.Println("Done!")
	}

	if err = data.ToFile(config.PasumanDataFile); err != nil {
		return err
	}

	return nil
}
