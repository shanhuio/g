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

package jwt

import (
	"bytes"
	"context"
	"io"

	"shanhu.io/g/errcode"
)

// Signer signs the token, returns the signature and the header.
type Signer interface {
	Header(ctx context.Context) (*Header, error)
	Sign(ctx context.Context, h *Header, data []byte) ([]byte, error)
}

// EncodeAndSign signs and encodes a claim set and signs it.
func EncodeAndSign(ctx context.Context, c *ClaimSet, s Signer) (string, error) {
	h, err := s.Header(ctx)
	if err != nil {
		return "", errcode.Annotate(err, "get header")
	}
	hb, err := h.encode()
	if err != nil {
		return "", errcode.Annotate(err, "encode header")
	}

	cb, err := c.encode()
	if err != nil {
		return "", errcode.Annotate(err, "encode claims")
	}
	buf := new(bytes.Buffer)
	io.WriteString(buf, hb)
	io.WriteString(buf, ".")
	io.WriteString(buf, cb)
	sig, err := s.Sign(ctx, h, buf.Bytes())
	if err != nil {
		return "", errcode.Annotate(err, "signing token")
	}
	io.WriteString(buf, ".")
	io.WriteString(buf, encodeSegmentBytes(sig))
	return buf.String(), nil
}
