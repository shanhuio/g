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
	"context"
	"strings"
	"time"

	"shanhu.io/pub/errcode"
)

// Verifier verifies the token.
type Verifier interface {
	Verify(ctx context.Context, h *Header, data, sig []byte, t time.Time) error
}

// DecodeAndVerify decodes and verifies a token.
func DecodeAndVerify(
	ctx context.Context, token string, v Verifier, t time.Time,
) (*Token, error) {
	decoded, err := Decode(token)
	if err != nil {
		return nil, errcode.Annotate(err, "decode token")
	}
	if err := Verify(ctx, decoded, v, t); err != nil {
		return nil, err
	}
	return decoded, nil
}

// Verify verifies if a decoded token has the valid signature.
func Verify(ctx context.Context, tok *Token, v Verifier, t time.Time) error {
	if v != nil {
		if err := v.Verify(
			ctx, tok.Header, tok.Payload, tok.Signature, t,
		); err != nil {
			return errcode.Annotate(err, "verify signature")
		}
	}

	_, err := CheckTime(tok.ClaimSet, t)
	return err
}

// Decode decodes the token without verifying it.
func Decode(token string) (*Token, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errcode.InvalidArgf(
			"invalid token: %d parts", len(parts),
		)
	}

	h, c, sig := parts[0], parts[1], parts[2]
	header := new(Header)
	if err := decodeSegment(h, header); err != nil {
		return nil, errcode.InvalidArgf("decode header: %s", err)
	}

	payload := []byte(token[:len(h)+1+len(c)])
	sigBytes, err := decodeSegmentBytes(sig)
	if err != nil {
		return nil, errcode.InvalidArgf("decode signature: %s", err)
	}
	claims, err := decodeClaimSet(c)
	if err != nil {
		return nil, errcode.InvalidArgf("decode claims: %s", err)
	}

	return &Token{
		Header:    header,
		ClaimSet:  claims,
		Payload:   payload,
		Signature: sigBytes,
	}, nil
}
