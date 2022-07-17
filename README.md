# pasuman: a command-line password manager

![](./assets/banner.svg)

![CI](https://github.com/norbjd/pasuman/actions/workflows/main.yml/badge.svg?branch=main)

pasuman is a command-line password manager designed to be:

- easy to use
- secure
- lightweight
- small:
    - minimum lines of code so it is easy to review
    - minimal number of dependencies
- cross-platform (GNU/Linux, MacOS)

## 📹 Demo (in 3 minutes!)

https://user-images.githubusercontent.com/26850303/178991840-25686d44-de03-4d38-be79-e9c82e3230be.mp4

## 🏃 Getting started

Download latest pasuman from [releases page](https://github.com/norbjd/pasuman/releases). For example, from a Linux machine, run:

```shell
VERSION=v1.0.0
curl -Lo pasuman https://github.com/norbjd/pasuman/releases/download/$VERSION/pasuman-linux-amd64
chmod u+x pasuman
mv pasuman /usr/local/bin
```

From a MacOS machine, replace `pasuman-linux-amd64` with the right OS and architecture.

On first use, run `pasuman master-password` to set a strong master password (see [Security > Choosing a strong master password](#choosing-a-strong-master-password)).

When asked to enter data directory:

```
Enter data directory (leave empty for default: /home/norbjd/.pasuman)
```

Choose a directory where your passwords (encrypted) will be stored. Leave default if you don't know, otherwise you can for example set the path to a USB key (`/run/media/norbjd/my-key/.pasuman`) that you can use if you have multiple devices and want to share your passwords between these devices.

Then, choose a **strong** master password. You will always be able to change it by running again the command `pasuman master-password`.

You can now use pasuman to add entries, search them, etc. See the demo video above for a preview of what you can do with pasuman.

You can also enable autocompletion for a better experience, for example by adding `<(source pasuman completion bash)` in your `.bashrc` if you are using `bash` (same command exists for `fish`, `powershell`, and `zsh`).

## ℹ️ Help

```
Copyright (C) 2022 norbjd
License GPLv3: GNU GPL version 3 <https://gnu.org/licenses/gpl.html>
This is free software: you are free to change and redistribute it.
There is NO WARRANTY, to the extent permitted by law.

A command-line password manager

Usage:
  pasuman [command]

Available Commands:
  add             Add an entry
  completion      Generate the autocompletion script for the specified shell
  generate        Generate a random password
  get             Get an entry
  help            Help about any command
  list            List entries
  list-profiles   List different profiles
  man             Generate man pages
  master-password Set or change master password
  remove          Remove an entry
  remove-lock     Remove lock
  search          Search an entry by a term
  update          Update an entry

Flags:
  -h, --help             help for pasuman
      --profile string   Profile (default "default")
  -v, --version          version for pasuman
```

## 💡 Model

pasuman stores entries.

An entry is composed of:

- an unique ID: if unspecified, a UUID v4 is generated
- a description (optional)
- some tags (optional): to help classify entries
- a site (optional)
- an ID (optional)
- a password

Only IDs and passwords are considered sensitive data.

## 💽 Storage

All entries are stored on disk, in simple JSON file(s). Sensitive data is stored securely (see [Security > ID and password storage](#id-and-password-storage)).

```json
{
  "master_password": "$argon2id$v=19$m=262144,t=16,p=4$YvSoJz5jjGQWiflI0pP1bW4R+b/FmMOYoypEp8eHHaKeasv2ikt/PpQQUrOXyFB0uKiHOUEc6gSG9SyqtqFTfw$AG/SFTkMBycYb7R0Q0b/me31G2EmAvoa8i7vRgAFI+k",
  "entries": [
    {
      "unique_id": "8a1c5d7f-8c6c-4d7f-a1a8-756e357d2d5e",
      "description": "This is a description",
      "tags": [
        "tag1",
        "tag2"
      ],
      "site": "https://mysupersite.pasuman",
      "id": "6RydipYrXpV+UdLms37CxfgEt9wl+IE17Tt9TmMjtzR3eGNe8eE3spohg5wmdjCIffqzRFLb0LIjc0z30UjYQA==*sH59oYEASE97QTTL*WhrPBXkxA1T8Q6d5...CPqo",
      "password": "DfHhyZ1NFdGILKgtdKFFjK8r2xSrC73vhWWhCih/i0jL2pLe9qGUFDD3tPHGzx4lanIYfY86JkZ9+ClYc9BpQg==*g4HIQQqaJRkaWyK3*POIpZyhDKphfZ4Vq...BJK8"
    }
  ]
}
```

## 🔒 Security

### Master password

Master password is hashed using [Argon2id](https://www.ietf.org/rfc/rfc9106.html), using a random salt. Parameters used can be seen [in the code](pkg/masterpassword/master_password.go).

The unhashed master password stays in memory only during the `pasuman` process life.

### ID and password storage

Encryption of an ID or a password is done in two steps:

- derive a key from the master password using [Argon2id](https://www.ietf.org/rfc/rfc9106.html) and a random salt. Parameters used can be seen [in the code](pkg/encrypt/encrypt.go)
- use AES with Galois/Counter Mode (AES-GCM) to encrypt the ID or password using the previously generated key

The unencrypted ID or password stays in memory only during the `pasuman` process life.

### Choosing a strong master password

[EFF Dice-Generated Passphrases](https://www.eff.org/dice) can be used as strong master passwords. Just be sure that you can remember your master password; otherwise, access to all your passwords stored by pasuman will be lost forever.

### Report a security issue

To report a critical security or vulnerability issue, contact me at `<my_github_username>+pasuman <at> googlemail <dot> com`.

## ❓ FAQ

**Q**: Why did I get ``Error: file is locked: use `pasuman remove-lock` command to remove it (only if no other pasuman process is running!)``?

**A**: This can happen:
- after you quit pasuman brutally (example: CTRL+C during the password prompt)
- when another pasuman process is running
- if pasuman exits brutally (for example if it panics; this should not happen, please open [an issue](https://github.com/norbjd/pasuman/issues) if so)

A lock have been implemented to avoid concurrent executions of pasuman, that can lead to data loss or an unexpected state.

To solve the issue, check that no other pasuman process is running, and run `pasuman remove-lock`.

## ✉️ Contact

If you have any question, feel free to create [a new discussion](https://github.com/norbjd/pasuman/discussions) or [an issue](https://github.com/norbjd/pasuman/issues) if you noticed a bug. You can also [reach me on Twitter (@norbjd)](https://twitter.com/norbjd).

## ⚖️ License

pasuman is licensed under the GNU General Public License (version 3), see [COPYING](COPYING).

```
Copyright (C) 2022 norbjd

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, version 3 of the License.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
```

**Note**: `pasuman --help` and `pasuman --version` shows a copyright notice, as defined by [GNU Coding Standards](https://www.gnu.org/prep/standards/standards.html#g_t_002d_002dversion).
