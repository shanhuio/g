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

package certutil

import (
	"crypto/tls"
	"crypto/x509"
	"math/big"
	"testing"
	"time"
)

func TestGetter(t *testing.T) {
	now := time.Now()
	sleep := func(d time.Duration) {
		now = now.Add(d)
	}
	readNow := func() time.Time {
		return now
	}

	config := &getterConfig{
		now:   readNow,
		sleep: sleep,
	}

	f := func(hello *tls.ClientHelloInfo) (*tls.Certificate, error) {
		// TODO(h8liu): use a better certificate
		cert := new(tls.Certificate)
		cert.Leaf = &x509.Certificate{
			SerialNumber: big.NewInt(12345),
			NotBefore:    now,
			NotAfter:     now.Add(time.Hour),
		}
		return cert, nil
	}

	start := now
	get := wrapAutoCert(f, config)
	hello := &tls.ClientHelloInfo{ServerName: "example.com"}
	if _, err := get(hello); err != nil {
		t.Fatal("get certificate: ", err)
	}

	atLeast := start.Add(getterDelay)
	if now.Before(atLeast) {
		t.Errorf(
			"should have delayed to +%s, but at +%s",
			getterDelay, now.Sub(start),
		)
	}

	// get again
	atLeast = now.Add(getterDelay)
	before := now
	if _, err := get(hello); err != nil {
		t.Fatal("get certificate: ", err)
	}
	if now.Before(atLeast) {
		t.Errorf(
			"should have delayed to +%s, but at +%s",
			getterDelay, now.Sub(before),
		)
	}

	mature := start.Add(getterMature)
	if now.Before(mature) {
		now = mature
	}

	// get again, this time after it is matured.
	before = now
	if _, err := get(hello); err != nil {
		t.Fatal("get certificate: ", err)
	}
	if now.After(before) {
		t.Errorf("time moved by %s on second get", now.Sub(before))
	}
}
