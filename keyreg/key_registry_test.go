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

package keyreg

import (
	"testing"

	"net/http/httptest"
	"net/url"
	"os"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/rsautil"
)

func TestMemKeyRegistry(t *testing.T) {
	keyBytes, err := os.ReadFile("testdata/keys/yumuzi")
	if err != nil {
		t.Fatal(err)
	}

	keys, err := rsautil.ParsePublicKeys(keyBytes)
	if err != nil {
		t.Fatal(err)
	}

	s := NewMemKeyRegistry()
	s.Set("h8liu", keys)

	got, err := s.Keys("h8liu")
	if err != nil {
		t.Fatal(err)
	}
	if len(keys) != len(got) {
		t.Errorf("want %d keys, got %d", len(keys), len(got))
	}

	for i, k := range got {
		h := k.HashStr()
		want := keys[i].HashStr()
		if h != want {
			t.Errorf("for key %d, want hash %q, got %q", i, want, h)
		}
	}
}

func testFileKeyRegistry(t *testing.T, ks KeyRegistry) {
	t.Helper()

	for _, test := range []struct {
		user   string
		hashes []string
	}{
		{"h8liu", []string{"zFet8qN1eNMvCQQqZRLy9Yxe-smJa8jmu30rOvBMeXw"}},
		{"yumuzi", []string{
			"Rxf8wK9cdKA6Zhn6KtVjSF3WUPLfnjbRlHuduSiOMsg",
			"zUcyOLg7_GzRTo4MDpyTnIxh6gqgGemUq0si_NjRXc4",
		}},
		{"xuduoduo", nil},
	} {
		t.Logf("test key for: %s", test.user)

		got, err := ks.Keys(test.user)
		if err != nil {
			if len(test.hashes) == 0 && errcode.IsNotFound(err) {
				continue
			}
			t.Fatal(err)
		}
		if len(got) != len(test.hashes) {
			t.Errorf("want %d keys, got %d", len(test.hashes), len(got))
			continue
		}

		for i, want := range test.hashes {
			if gotHash := got[i].HashStr(); want != gotHash {
				t.Errorf("key %d, want hash %q, got %q", i, want, gotHash)
			}
		}
	}
}

func TestFileKeyRegistry(t *testing.T) {
	s := NewFileKeyRegistry(map[string]string{
		"h8liu":  "testdata/keys/h8liu",
		"yumuzi": "testdata/keys/yumuzi",
	})

	testFileKeyRegistry(t, s)
}

func TestWebKeyRegistry(t *testing.T) {
	static := aries.NewStaticFiles("testdata/keys")
	s := httptest.NewServer(aries.Serve(static))
	defer s.Close()

	t.Log(s.URL)
	u, err := url.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	ks := NewWebKeyRegistry(u)
	testFileKeyRegistry(t, ks)
}
