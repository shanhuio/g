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

package https

import (
	"testing"

	"shanhu.io/g/errcode"
)

func TestNewCACert(t *testing.T) {
	cert, err := NewCACert("test.shanhu.io")
	if err != nil {
		t.Fatalf("NewCACert() got error: %s", err)
	}

	if _, err := cert.X509KeyPair(); err != nil {
		t.Fatalf("convert to tls cert got error: %s", err)
	}
}

func TestMakeRSACertWithNoHost(t *testing.T) {
	_, err := MakeRSACert(&CertConfig{}, 0)
	if !errcode.IsInvalidArg(err) {
		t.Errorf("expect invalid arg error without hosts, got %s", err)
	}
}
