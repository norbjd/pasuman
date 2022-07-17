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

package data

import (
	"encoding/json"
	"os"
)

type Data struct {
	MasterPassword string  `json:"master_password"`
	Entries        []Entry `json:"entries"`
}

type Entry struct {
	UniqueID    string   `json:"unique_id"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Site        string   `json:"site"`
	ID          string   `json:"id"`
	Password    string   `json:"password"`
}

func (data *Data) FromFile(file string) error {
	byteContents, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	return json.Unmarshal(byteContents, data)
}

func (data *Data) ToFile(file string) error {
	byteContents, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(file, byteContents, os.ModeAppend)
}
