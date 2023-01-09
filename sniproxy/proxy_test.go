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

package sniproxy

import (
	"testing"

	"context"
	"net"
	"net/http"
	"net/http/httptest"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/https/httpstest"
	"shanhu.io/pub/httputil"
)

func TestProxy(t *testing.T) {
	tlsConfigs, err := httpstest.NewTLSConfigs([]string{
		"site1.com",
		"site2.com",
	})
	if err != nil {
		t.Fatal(err)
	}

	site1 := httptest.NewUnstartedServer(aries.StringFunc("site1"))
	site1.TLS = tlsConfigs.Server
	site1.StartTLS()
	defer site1.Close()

	site2 := httptest.NewUnstartedServer(aries.StringFunc("site2"))
	site2.TLS = tlsConfigs.Server
	site2.StartTLS()
	defer site2.Close()

	httpAddr := func(s *httptest.Server) *net.TCPAddr {
		return s.Listener.Addr().(*net.TCPAddr)
	}

	dialer := &tcpDialer{
		raddrs: map[string]*net.TCPAddr{
			"site1.com": httpAddr(site1),
			"site2.com": httpAddr(site2),
		},
	}

	localhost := &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)}
	lis, err := net.ListenTCP("tcp", localhost)
	if err != nil {
		t.Fatal(err)
	}
	defer lis.Close()

	p := newProxy(dialer)
	ctx := context.Background()
	go p.serve(ctx, lis)

	client := &http.Client{
		Transport: tlsConfigs.Sink(lis.Addr().String()),
	}

	got, err := httputil.GetString(client, "https://site1.com")
	if err != nil {
		t.Fatal("get site1:", err)
	}
	if got != "site1" {
		t.Errorf("want site1, got %q", got)
	}

	got, err = httputil.GetString(client, "https://site2.com")
	if err != nil {
		t.Fatal("get site2:", err)
	}
	if got != "site2" {
		t.Errorf("want site2, got %q", got)
	}
}
