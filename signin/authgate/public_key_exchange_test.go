package authgate

import (
	"net/http/httptest"
	"testing"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/keyreg"
	"shanhu.io/pub/rsautil"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/signin/signinapi"
)

func TestPublicKeyExchange(t *testing.T) {
	const user = "h8liu"

	pri, pub, err := rsautil.GenerateKey(nil, 1024)
	if err != nil {
		t.Fatal("generate rsa key: ", err)
	}

	pubKey, err := rsautil.NewPublicKey(pub)
	if err != nil {
		t.Fatal("compile public key: ", err)
	}
	priKey, err := rsautil.ParsePrivateKey(pri)
	if err != nil {
		t.Fatal("parse private key: ", err)
	}

	kr := keyreg.NewMemKeyRegistry()
	kr.Set(user, []*rsautil.PublicKey{pubKey})

	gate := New(&Config{SessionKey: []byte("test-key")})
	ex := NewPublicKeyExchange(gate, kr)

	r := aries.NewRouter()
	r.Index(aries.StringFunc("index"))
	r.Call("/signin", ex.Exchange)

	const secret = "mikimakibubabu"
	u := aries.NewRouter()
	r.Get("/secret", aries.StringFunc(secret))

	set := &aries.ServiceSet{
		Auth:  gate,
		Guest: r,
		User:  u,
	}

	s := httptest.NewServer(aries.Serve(set))
	client := httputil.NewClientMust(s.URL)

	sig, err := signer.RSASignTime(priKey)
	if err != nil {
		t.Fatal("sign time: ", err)
	}
	req := &signinapi.Request{
		User:       user,
		SignedTime: sig,
	}
	creds := new(signinapi.Creds)
	if err := client.Call("/signin", req, creds); err != nil {
		t.Fatal("signin: ", err)
	}

	client.TokenSource = httputil.NewStaticToken(creds.Token)
	got, err := client.GetString("/secret")
	if err != nil {
		t.Fatal("get secret: ", err)
	}
	if got != secret {
		t.Errorf("got secret %q, want %q", got, secret)
	}
}
