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

package states

import (
	"strings"

	"shanhu.io/pub/errcode"
)

// GetToken gets a token string from the given key. The fetched value is
// treated as a string and whitespaces are trimmed.
func GetToken(ctx C, s States, key string) (string, error) {
	bs, err := s.Get(ctx, key)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bs)), nil
}

// GetTokenDefault gets a token string from the given keky. The fetched value
// is treated as a string and whitespaces are trimmed. If the key does not
// exist, v is returned.
func GetTokenDefault(ctx C, s States, key, v string) (string, error) {
	tok, err := GetToken(ctx, s, key)
	if errcode.IsNotFound(err) {
		return v, nil
	}
	return tok, err
}
