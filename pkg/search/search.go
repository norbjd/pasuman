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

package search

import (
	"strings"

	"github.com/norbjd/pasuman/pkg/data"
	"github.com/norbjd/pasuman/pkg/list"
)

func Search(profile string, term string, caseSensitive bool) ([]data.Entry, error) {
	entries, err := list.List(profile)
	if err != nil {
		return nil, err
	}

	if !caseSensitive {
		term = strings.ToLower(term)
	}

	searchResults := make([]data.Entry, 0)

	for _, entry := range entries {
		uniqueID := entry.UniqueID
		if !caseSensitive {
			uniqueID = strings.ToLower(uniqueID)
		}

		if strings.Contains(uniqueID, term) {
			searchResults = append(searchResults, entry)

			continue
		}

		description := entry.Description
		if !caseSensitive {
			description = strings.ToLower(description)
		}

		if strings.Contains(description, term) {
			searchResults = append(searchResults, entry)

			continue
		}

		tags := entry.Tags
		if !caseSensitive {
			for idx := range tags {
				tags[idx] = strings.ToLower(tags[idx])
			}
		}

		for _, tag := range tags {
			if strings.Contains(tag, term) {
				searchResults = append(searchResults, entry)
			}
		}

		site := entry.Site
		if !caseSensitive {
			site = strings.ToLower(site)
		}

		if strings.Contains(site, term) {
			searchResults = append(searchResults, entry)

			continue
		}
	}

	return searchResults, nil
}
