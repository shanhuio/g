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

package aries

import (
	"context"
	"testing"

	"net"
	"net/http"
	"net/url"
	"path/filepath"

	"shanhu.io/g/httputil"
)

func TestDefaultAddr(t *testing.T) {
	addr := DefaultAddr("aries")
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		t.Fatalf("bad addr %q", addr)
	}

	if host != "localhost" {
		t.Errorf("got %q, want localhost", host)
	}

	addr2 := DefaultAddr("aries")
	if addr2 != addr {
		t.Errorf("got differnt %q, want %q", addr2, addr)
	}

	addr3 := DefaultAddr("aries2")
	if addr3 == addr {
		t.Errorf("got same %q as %q, want different", addr3, addr)
	}
}

func TestListen_unix(t *testing.T) {
	tmp := t.TempDir()
	sock := filepath.Join(tmp, "server.sock")
	lis, err := Listen(sock)
	if err != nil {
		t.Fatal("listen:", err)
	}
	defer lis.Close()

	s := &http.Server{Handler: Serve(StringFunc("hello"))}
	serveErr := make(chan error)
	go func() { serveErr <- s.Serve(lis) }()

	ctx := context.Background()

	client := httputil.NewUnixClient(sock)
	res, err := client.GetString("/")
	if err != nil {
		t.Fatal("get:", err)
	}

	if res != "hello" {
		t.Errorf("got %q, want hello", res)
	}

	if err := s.Shutdown(ctx); err != nil {
		t.Fatal("shutdown:", err)
	}
	if err := <-serveErr; err != http.ErrServerClosed {
		t.Errorf("got %v, want %v", err, http.ErrServerClosed)
	}
}

func TestListen_http(t *testing.T) {
	lis, err := Listen("localhost:0")
	if err != nil {
		t.Fatal("listen:", err)
	}
	defer lis.Close()

	s := &http.Server{Handler: Serve(StringFunc("hello"))}
	serveErr := make(chan error)
	go func() { serveErr <- s.Serve(lis) }()

	ctx := context.Background()

	client := &httputil.Client{
		Server: &url.URL{
			Scheme: "http",
			Host:   lis.Addr().String(),
		},
	}

	res, err := client.GetString("/")
	if err != nil {
		t.Fatal("get:", err)
	}

	if res != "hello" {
		t.Errorf("got %q, want hello", res)
	}

	if err := s.Shutdown(ctx); err != nil {
		t.Fatal("shutdown:", err)
	}
	if err := <-serveErr; err != http.ErrServerClosed {
		t.Errorf("got %v, want %v", err, http.ErrServerClosed)
	}
}
