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

package keyreg

import (
	"net/url"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/rsautil"
)

// WebKeyRegistry is a storage of public keys backed by a web site.
type WebKeyRegistry struct {
	client *httputil.Client
}

// NewWebKeyRegistry creates a new key store backed by a web site
// at the given base URL.
func NewWebKeyRegistry(base *url.URL) *WebKeyRegistry {
	client := &httputil.Client{Server: base}
	return &WebKeyRegistry{client: client}
}

// Keys returns the public keys of the given user.
func (s *WebKeyRegistry) Keys(user string) ([]*rsautil.PublicKey, error) {
	if !IsSimpleName(user) {
		return nil, errcode.InvalidArgf("unsupported user name: %q", user)
	}
	bs, err := s.client.GetBytes(user)
	if err != nil {
		return nil, err
	}
	return rsautil.ParsePublicKeys(bs)
}
