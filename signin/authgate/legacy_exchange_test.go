package authgate

import (
	"testing"

	"net/http/httptest"

	"shanhu.io/g/aries"
	"shanhu.io/g/httputil"
	"shanhu.io/g/keyreg"
	"shanhu.io/g/keyreg/testkeys"
	"shanhu.io/g/rsautil"
	"shanhu.io/g/signer"
	"shanhu.io/g/signin/signinapi"
)

func TestLegacyExchange(t *testing.T) {
	const user = "h8liu"

	pubKey, err := rsautil.NewPublicKey([]byte(testkeys.Pub1))
	if err != nil {
		t.Fatal("compile public key: ", err)
	}
	priKey, err := rsautil.ParsePrivateKey([]byte(testkeys.Pem1))
	if err != nil {
		t.Fatal("parse private key: ", err)
	}

	kr := keyreg.NewMemKeyRegistry()
	kr.Set(user, []*rsautil.PublicKey{pubKey})

	gate := New(&Config{SessionKey: []byte("test-key")})
	ex := NewLegacyExchange(gate, kr)

	r := aries.NewRouter()
	r.Index(aries.StringFunc("index"))
	r.Call("/signin", ex.Exchange)

	const secret = "mikima-kibubabu"
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
