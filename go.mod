// This file is part of pasuman (https://github.com/norbjd/pasuman).
//
// pasuman is a command-line password manager.
// Copyright (c) 2022 norbjd
//
// This program is free software; you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, version 3.
//
// This program is distributed in the hope that it will be useful, but
// WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the GNU
// General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

module github.com/norbjd/pasuman

go 1.18

require (
	github.com/alexedwards/argon2id v0.0.0-20211130144151-3585854a6387
	github.com/spf13/cobra v1.5.0
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d
	golang.org/x/term v0.0.0-20220526004731-065cf7ba2467
)

// tests
require (
	github.com/google/uuid v1.3.0
	github.com/spf13/pflag v1.0.5
	github.com/stretchr/testify v1.7.5
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/inconshreveable/mousetrap v1.0.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	golang.org/x/sys v0.0.0-20220702020025-31831981b65f // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
