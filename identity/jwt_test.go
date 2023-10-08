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

package identity

import (
	"context"
	"testing"

	"time"

	"shanhu.io/g/jwt"
)

func TestJWT(t *testing.T) {
	ctx := context.Background()

	now := time.Now()
	core := NewMemCore(func() time.Time { return now })

	coreConfig := SingleKeyCoreConfig(now.Add(time.Hour))
	if _, err := core.Init(coreConfig); err != nil {
		t.Fatal("init core: ", err)
	}

	signer := newJWTSigner(core)

	claim := &jwt.ClaimSet{
		Iss: "shanhu.io",
		Aud: "doorway.homedrv",
		Iat: now.Unix(),
		Exp: now.Add(time.Hour).Unix(),
		Sub: "core.homedrv",
	}

	encoded, err := jwt.EncodeAndSign(ctx, claim, signer)
	if err != nil {
		t.Fatal("generate token: ", err)
	}

	t.Log("token: ", encoded)

	v := newJWTVerifier(core)
	decoded, err := jwt.DecodeAndVerify(ctx, encoded, v, now)
	if err != nil {
		t.Fatal("decode and verify token: ", err)
	}

	keyID := decoded.Header.KeyID

	pub, err := signer.rsaPublicKeyPEM(ctx, keyID)
	if err != nil {
		t.Fatal("read public key: ", err)
	}

	t.Logf("public key:\n%s", pub)
}
