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

package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"

	"github.com/norbjd/pasuman/pkg/util"
)

const (
	defaultDirMode  = os.FileMode(0o700)
	defaultFileMode = os.FileMode(0o600)
)

var (
	ConfigDir         string
	PasumanConfigFile string
	PasumanDataFile   string
)

var errCannotInitConfigFile = errors.New("cannot init config file")

func Init(profile string, createDataFile bool) {
	if ConfigDir == "" {
		configDir, err := os.UserConfigDir()
		if err != nil {
			log.Fatal("No config dir")
		}

		ConfigDir = configDir
	}

	pasumanConfigDir := ConfigDir + string(os.PathSeparator) + "pasuman"
	PasumanConfigFile = pasumanConfigDir + string(os.PathSeparator) + "config.json"

	config := initConfigFileIfNecessary(pasumanConfigDir)
	initDataFile(config, profile, createDataFile)
}

type Config struct {
	DataDirectory string `json:"data_directory"`
}

func GetConfig() Config {
	pasumanConfigFileHandler, err := os.Open(PasumanConfigFile)
	if err != nil {
		log.Fatal(err)
	}

	var config Config

	jsonDecoder := json.NewDecoder(pasumanConfigFileHandler)
	if err = jsonDecoder.Decode(&config); err != nil {
		log.Fatal(err)
	}

	return config
}

func GetDataFile(config Config, profile string) string {
	return fmt.Sprintf("%s%s%s.json", config.DataDirectory, string(os.PathSeparator), profile)
}

func initConfigFileIfNecessary(pasumanConfigDir string) Config {
	dirInfo, err := os.Stat(pasumanConfigDir)
	if errors.Is(err, os.ErrNotExist) {
		if err := os.MkdirAll(pasumanConfigDir, defaultDirMode); err != nil {
			log.Fatal(err)
		}

		dirInfo, err = os.Stat(pasumanConfigDir)
	}

	if err != nil {
		log.Fatal(err)
	}

	if !dirInfo.IsDir() {
		log.Fatalf("%s should be a directory", pasumanConfigDir)
	}

	fileInfo, err := os.Stat(PasumanConfigFile)

	if errors.Is(err, os.ErrNotExist) {
		fileInfo, err = initPasuman()
	}

	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		log.Fatalf("%s should be a file, not a directory", PasumanConfigFile)
	}

	config := GetConfig()

	if _, err := os.Stat(config.DataDirectory); errors.Is(err, os.ErrNotExist) {
		writeErr := os.MkdirAll(config.DataDirectory, defaultDirMode)
		if writeErr != nil {
			log.Fatal("Cannot create data directory")
		}
	}

	return config
}

func initDataFile(config Config, profile string, create bool) {
	PasumanDataFile = GetDataFile(config, profile)

	if create {
		if _, err := os.Stat(PasumanDataFile); errors.Is(err, os.ErrNotExist) {
			writeErr := os.WriteFile(PasumanDataFile, []byte("{}"), defaultFileMode)
			if writeErr != nil {
				log.Fatalf("Cannot init profile %s file", profile)
			}
		}
	}
}

func initPasuman() (fs.FileInfo, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	defaultDataDirectory := homeDir + string(os.PathSeparator) + ".pasuman"

	fmt.Printf("Enter data directory (leave empty for default: %s): ", defaultDataDirectory)

	dataDirectory, err := util.ReadLine()
	if err != nil {
		return nil, err
	}

	dataDirectory = strings.TrimSpace(dataDirectory)

	if dataDirectory == "" {
		dataDirectory = defaultDataDirectory
	}

	config := Config{DataDirectory: dataDirectory}

	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return nil, err
	}

	writeErr := os.WriteFile(PasumanConfigFile, configJSON, defaultFileMode)
	if writeErr != nil {
		return nil, errCannotInitConfigFile
	}

	fileInfo, err := os.Stat(PasumanConfigFile)
	if err != nil {
		return nil, err
	}

	return fileInfo, err
}
