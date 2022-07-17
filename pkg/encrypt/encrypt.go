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

package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"
	"strings"

	"github.com/norbjd/pasuman/pkg/constants"
	"github.com/norbjd/pasuman/pkg/util"
	"golang.org/x/crypto/argon2"
)

const (
	encryptedMessageStringSeparator   = "*"
	encryptedMessageSplitStringLength = 3

	saltLength  = 64
	nonceLength = 12

	argon2Time      = 16
	argon2Memory    = 256 * 1024
	argon2Threads   = 4
	argon2KeyLength = 32

	encryptedStringWithPaddingLength = 512
	paddingSeparator                 = '\x00'
)

var errEncryptedMessageStringInvalid = errors.New("encrypted message string is invalid")

type EncryptedMessage struct {
	base64Salt    string
	base64Nonce   string
	base64Message string
}

func (e *EncryptedMessage) String() string {
	return fmt.Sprintf("%s%s%s%s%s",
		e.base64Salt, encryptedMessageStringSeparator,
		e.base64Nonce, encryptedMessageStringSeparator,
		e.base64Message)
}

func (e *EncryptedMessage) FromString(s string) error {
	split := strings.Split(s, encryptedMessageStringSeparator)
	if len(split) != encryptedMessageSplitStringLength {
		return fmt.Errorf("%w: %s", errEncryptedMessageStringInvalid, s)
	}

	e.base64Salt = split[0]
	e.base64Nonce = split[1]
	e.base64Message = split[2]

	return nil
}

func Encrypt(masterPassword, stringToEncrypt string) (string, error) {
	if masterPassword == "" {
		return "", util.ErrMasterPasswordMustNotBeEmpty
	}

	salt := make([]byte, saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(masterPassword), salt,
		argon2Time, argon2Memory, argon2Threads, argon2KeyLength)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	rightPaddingLength := encryptedStringWithPaddingLength - len(stringToEncrypt)

	if rightPaddingLength > 0 {
		rightPadding := make([]byte, rightPaddingLength)
		rightPadding[0] = paddingSeparator

		for i := 1; i < rightPaddingLength; i++ {
			num, err := rand.Int(rand.Reader, big.NewInt(int64(len(constants.Alphabet))))
			if err != nil {
				return "", err
			}

			rightPadding[i] = constants.Alphabet[num.Int64()]
		}

		stringToEncrypt += string(rightPadding)
	}

	plaintext := []byte(stringToEncrypt)

	nonce := make([]byte, nonceLength)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := aesgcm.Seal(nil, nonce, plaintext, nil)

	encryptedMessage := EncryptedMessage{
		base64Salt:    base64.StdEncoding.EncodeToString(salt),
		base64Nonce:   base64.StdEncoding.EncodeToString(nonce),
		base64Message: base64.StdEncoding.EncodeToString(ciphertext),
	}

	return encryptedMessage.String(), nil
}

func Decrypt(masterPassword, stringToDecrypt string) (string, error) {
	var encryptedMessage EncryptedMessage
	if err := encryptedMessage.FromString(stringToDecrypt); err != nil {
		return "", err
	}

	salt, err := base64.StdEncoding.DecodeString(encryptedMessage.base64Salt)
	if err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(masterPassword), salt,
		argon2Time, argon2Memory, argon2Threads, argon2KeyLength)

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedMessage.base64Message)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce, err := base64.StdEncoding.DecodeString(encryptedMessage.base64Nonce)
	if err != nil {
		return "", err
	}

	plaintext, err := aesgcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	decrypted := strings.Split(string(plaintext), string(paddingSeparator))[0]

	return decrypted, nil
}
