package aries

import (
	"context"
	"testing"

	"net"
	"net/http"
	"net/url"
	"path/filepath"

	"shanhu.io/pub/httputil"
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
