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

package generate

import (
	"crypto/rand"
	"math/big"

	"github.com/norbjd/pasuman/pkg/constants"
)

func Generate(length uint) (string, error) {
	generated := make([]byte, length)

	for i := uint(0); i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(constants.Alphabet))))
		if err != nil {
			return "", err
		}

		generated[i] = constants.Alphabet[num.Int64()]
	}

	return string(generated), nil
}
