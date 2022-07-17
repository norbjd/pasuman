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

package cmd

import (
	"fmt"
)

var (
	GitVersion = "unknown"
	GitCommit  = "unknown"
	BuildDate  = "unknown"
)

func Version() string {
	message := "pasuman"

	if GitVersion != "unknown" && GitVersion != "" {
		message += fmt.Sprintf(" %s -", GitVersion)
	}

	message += fmt.Sprintf(" %s - built %s\n\n", GitCommit, BuildDate)

	message += "Copyright (C) 2022 norbjd\n" +
		"License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>\n" +
		"This is free software: you are free to change and redistribute it.\n" +
		"There is NO WARRANTY, to the extent permitted by law.\n"

	return message
}
