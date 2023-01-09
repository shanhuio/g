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

package creds

import (
	"crypto/rsa"

	"shanhu.io/pub/rsautil"
)

func parsePrivateKey(name string, bs []byte, tty bool) (
	*rsa.PrivateKey, error,
) {
	if tty {
		return rsautil.ParsePrivateKeyTTY(name, bs)
	}
	return rsautil.ParsePrivateKey(bs)
}

func readPrivateKey(pemFile string, tty bool) (*rsa.PrivateKey, error) {
	if tty {
		return rsautil.ReadPrivateKeyTTY(pemFile)
	}
	return rsautil.ReadPrivateKey(pemFile)
}

func readEndpointKey(p *Endpoint) (*rsa.PrivateKey, error) {
	tty := !p.NoTTY
	if p.Key != nil {
		return parsePrivateKey("key", p.Key, tty)
	}
	return readPrivateKey(p.PemFile, tty)
}
