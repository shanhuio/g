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

package timeutil

import (
	"encoding/base64"
	"io"
	"time"

	"shanhu.io/pub/errcode"
)

// Challenge is a timestamp with a crypto random nonce. A server can provide
// an HMAC signed challenge for clients via RPC, and use it as a structure
// to restrict the client into a time window that is defined by the server.
// It also provides a reliable clock source and a nonce source. The nonce
// can be used as a one-time token to defend replay attacks.
type Challenge struct {
	N string // Nonce.
	T *Timestamp
}

// NewChallenge creates a challenge with the given timestamp t.
// rand is used to generate the nonce.
func NewChallenge(t time.Time, rand io.Reader) (*Challenge, error) {
	nonce := make([]byte, 16)
	if _, err := rand.Read(nonce); err != nil {
		return nil, errcode.Annotate(err, "read nonce")
	}

	nonceStr := base64.RawStdEncoding.EncodeToString(nonce)
	return &Challenge{
		N: nonceStr,
		T: NewTimestamp(t),
	}, nil
}
