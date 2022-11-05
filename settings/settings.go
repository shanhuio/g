// Copyright (C) 2022  Shanhu Tech Inc.
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

// Package settings provides a generic, simple key-value pair like interface
// for saving application settings.
package settings

// Settings is an interface for saving simple JSON object based settings.
type Settings interface {
	// Get gets a setting. Returns errcode.NotFound error when the
	// key is missing.
	Get(key string, v interface{}) error

	// Set sets a setting.
	Set(key string, v interface{}) error

	// Has checks if a key exists. It does not have to read the key.
	Has(key string) (bool, error)
}

// String gets a string-type value from the settings.
func String(b Settings, key string) (string, error) {
	var s string
	if err := b.Get(key, &s); err != nil {
		return "", err
	}
	return s, nil
}

// Bytes gets a []byte type value from the settings.
func Bytes(b Settings, key string) ([]byte, error) {
	var bs []byte
	if err := b.Get(key, &bs); err != nil {
		return nil, err
	}
	return bs, nil
}
