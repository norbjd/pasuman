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

package util

import (
	"bufio"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"golang.org/x/term"
)

var errInvalidTerminal = errors.New("stdin and stdout should be terminal")

// ReadLine - read from stdin and handles backspace.
// Inspired from https://gist.github.com/artyom/a59e2707976124f387f5
func ReadLine() (string, error) {
	if !term.IsTerminal(int(os.Stdin.Fd())) || !term.IsTerminal(int(os.Stdout.Fd())) {
		return "", errInvalidTerminal
	}

	oldState, err := term.MakeRaw(0)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = term.Restore(0, oldState)
	}()

	screen := struct {
		io.Reader
		io.Writer
	}{os.Stdin, os.Stdout}
	termScreen := term.NewTerminal(screen, "")

	line, err := termScreen.ReadLine()

	if errors.Is(err, io.EOF) {
		return "", nil
	}

	if err != nil {
		return "", err
	}

	return line, nil
}

// ReadPassword - read password from stdin and handle CTRL+C to restore.
// Inspired from https://groups.google.com/g/golang-nuts/c/DCl8xUJMJJ0.
func ReadPassword() (string, error) {
	// nolint: unconvert
	// int conversion is necessary to build for windows
	fd := int(syscall.Stdin)

	if !term.IsTerminal(fd) {
		// if not a terminal, use environment variable
		// warning: it should not be used outside of tests
		if password := os.Getenv("PASUMAN_MASTER_PASSWORD"); password != "" {
			return password, nil
		}
	}

	oldState, err := term.GetState(fd)
	if err != nil {
		return "", err
	}

	defer func() {
		_ = term.Restore(fd, oldState)
	}()

	sigch := make(chan os.Signal, 1)
	signal.Notify(sigch, os.Interrupt)

	go func() {
		for range sigch {
			_ = term.Restore(fd, oldState)

			fmt.Println()
			os.Exit(1)
		}
	}()

	password, err := term.ReadPassword(fd)
	if err != nil {
		return "", err
	}

	return string(password), nil
}

// CopyToClipboard - uses ANSI OSC 52 to copy text to clipboard,
// it may not work on some terminals (see https://github.com/ojroques/vim-oscyank/blob/main/README.md).
func CopyToClipboard(s string) {
	out := bufio.NewWriter(os.Stdout)
	fmt.Fprintf(out, "\033]52;c;%s\a", base64.StdEncoding.EncodeToString([]byte(s)))
	out.Flush()
}
