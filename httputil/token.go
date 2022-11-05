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

package httputil

import (
	"context"
	"net/http"
)

// SetAuthToken sets authorization header token.
func SetAuthToken(h http.Header, tok string) {
	if tok == "" {
		return
	}
	h.Set("Authorization", "Bearer "+tok)
}

// TokenSource is an interface that can provides a bearer token for
// authentication.
type TokenSource interface {
	Token(ctx context.Context, tr http.RoundTripper) (string, error)
}

// StaticToken is a token source that provides a fixed,
// static token.
type StaticToken struct {
	T string
}

// NewStaticToken creates a new static token provider.
func NewStaticToken(tok string) *StaticToken {
	return &StaticToken{T: tok}
}

// Token always returns the fixed, static token T.
func (s *StaticToken) Token(
	ctx context.Context, tr http.RoundTripper,
) (string, error) {
	return s.T, nil
}
