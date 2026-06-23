package signer

import (
	"testing"

	"reflect"

	"shanhu.io/g/rand"
)

func testSigner(t *testing.T, k []byte) {
	s := New(k)
	o := func(bs []byte) {
		signed := s.Sign(bs)
		ok, dat := s.Check(signed)
		if !ok {
			t.Error("check failed")
		} else if !reflect.DeepEqual(dat, bs) {
			t.Errorf("got %v, want %v", dat, bs)
		}

		h := s.SignHex(bs)
		ok, dat = s.CheckHex(h)
		if !ok {
			t.Error("check failed")
		} else if !reflect.DeepEqual(dat, bs) {
			t.Errorf("got %v, want %v", dat, bs)
		}
	}

	os := func(s string) { o([]byte(s)) }
	os("")
	os("something")
	os("            ")

	for range 5 {
		o(rand.Bytes(10))
	}
}

func TestSigner(t *testing.T) {
	testSigner(t, nil)
	testSigner(t, []byte{})
	for range 3 {
		testSigner(t, rand.Bytes(8))
	}
}
