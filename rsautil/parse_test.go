// Copyright (C) 2022  Shanhu Tech Inc.
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

package rsautil

import (
	"testing"

	"os"
	"path/filepath"
	"reflect"
)

func TestReadKey(t *testing.T) {
	tmp, err := os.MkdirTemp("", "rsautil")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

	// Recreate the test.pem key file to make sure it has the right permission
	// bits.
	testPem, err := os.ReadFile("testdata/test.pem")
	if err != nil {
		t.Fatal("read test key content: ", err)
	}

	privateKeyFile := filepath.Join(tmp, "test.pem")
	if err := os.WriteFile(privateKeyFile, testPem, 0600); err != nil {
		t.Fatal("create test key file: ", err)
	}

	privateKey, err := ReadPrivateKey(privateKeyFile)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := ReadPublicKey("testdata/test.pub")
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(&privateKey.PublicKey, publicKey) {
		t.Error("public/private key pair not matching")
	}
}
