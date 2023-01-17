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

package authgate

import (
	"testing"

	"net/http/httptest"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/keyreg"
	"shanhu.io/pub/keyreg/testkeys"
	"shanhu.io/pub/rsautil"
	"shanhu.io/pub/signer"
	"shanhu.io/pub/signin/signinapi"
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
