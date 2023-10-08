// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package creds

import (
	"bytes"
	"encoding/json"
	"os"
	"os/user"
	"path/filepath"

	"shanhu.io/g/jsonutil"
	"shanhu.io/g/osutil"
)

const homeDir = ".shanhu"

// Home returns the directory for saving the credentials and config files.
func Home() (string, error) {
	u, err := user.Current()
	if err != nil {
		return "", err
	}
	return filepath.Join(u.HomeDir, homeDir), nil
}

// MakeHome creates the home directory if it does not exist.
func MakeHome() error {
	h, err := Home()
	if err != nil {
		return err
	}
	return os.MkdirAll(h, 700)
}

// HomeFile returns the path of a file under the home directory.
func HomeFile(f string) (string, error) {
	h, err := Home()
	if err != nil {
		return "", err
	}
	return filepath.Join(h, f), nil
}

// ReadHomeFile reads the content of a file under the home directory.
// The file must be mode 0600.
func ReadHomeFile(f string) ([]byte, error) {
	p, err := HomeFile(f)
	if err != nil {
		return nil, err
	}
	return osutil.ReadPrivateFile(p)
}

// WriteHomeFile updates a file under the home directory.
func WriteHomeFile(f string, bs []byte) error {
	p, err := HomeFile(f)
	if err != nil {
		return err
	}
	return os.WriteFile(p, bs, 0600)
}

// WriteHomeJSONFile updates a file under the home directory with a
// JSON marshallable blob.
func WriteHomeJSONFile(f string, v interface{}) error {
	buf := new(bytes.Buffer)
	if err := jsonutil.Fprint(buf, v); err != nil {
		return err
	}
	return WriteHomeFile(f, buf.Bytes())
}

// ReadHomeJSONFile reads a file under the home directory into a JSON
// marshallable structure.
func ReadHomeJSONFile(f string, v interface{}) error {
	bs, err := ReadHomeFile(f)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}
