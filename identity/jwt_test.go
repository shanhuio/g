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
