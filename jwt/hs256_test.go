package jwt

import (
	"context"
	"testing"

	"time"

	"shanhu.io/g/rand"
)

func TestHS256(t *testing.T) {
	key := rand.Bytes(32) // 256 bits
	h := NewHS256(key, "")
	now := time.Now()
	c := &ClaimSet{
		Iss: "shanhu.io",
		Aud: "nextcloud",
		Iat: now.Unix(),
		Exp: now.Add(time.Hour).Unix(),
		Sub: "h8liu",
	}

	ctx := context.Background()

	tokStr, err := EncodeAndSign(ctx, c, h)
	if err != nil {
		t.Fatal("encode: ", err)
	}
	t.Log(tokStr)

	tok, err := DecodeAndVerify(ctx, tokStr, h, now)
	if err != nil {
		t.Fatal("decode: ", err)
	}

	if got, want := tok.ClaimSet.Iss, c.Iss; got != want {
		t.Errorf("got issuer %q, want %q", got, want)
	}
}
