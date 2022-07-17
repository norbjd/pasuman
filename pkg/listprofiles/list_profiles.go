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

package listprofiles

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	"github.com/norbjd/pasuman/pkg/config"
)

func ListProfiles() ([]string, error) {
	configFileHandler, err := os.Open(config.PasumanConfigFile)
	if err != nil {
		return nil, err
	}

	var config config.Config

	jsonDecoder := json.NewDecoder(configFileHandler)
	if err = jsonDecoder.Decode(&config); err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(config.DataDirectory)
	if err != nil {
		return nil, err
	}

	var profiles []string

	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			split := strings.Split(f.Name(), ".")
			profiles = append(profiles, strings.Join(split[:len(split)-1], "."))
		}
	}

	return profiles, nil
}
