package netutil

import (
	"testing"

	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
)

func TestKeepAliveListener(t *testing.T) {
	lis, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatal("create tcp listener: ", err)
	}

	wrap := WrapKeepAlive(lis)

	handler := func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "hello")
	}

	s := &http.Server{
		Handler: http.HandlerFunc(handler),
	}

	serveErr := make(chan error)
	go func() {
		serveErr <- s.Serve(wrap)
	}()

	addr := wrap.Addr().String()
	client := new(http.Client)

	u := &url.URL{
		Scheme: "http",
		Host:   addr,
	}
	resp, err := client.Get(u.String())
	if err != nil {
		t.Fatal("http get: ", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal("read body: ", err)
	}

	if got := string(body); got != "hello" {
		t.Errorf("body got %q, want `hello`", got)
	}

	resp.Body.Close()

	ctx := context.Background()
	if err := s.Shutdown(ctx); err != nil {
		t.Fatal("shutdown server: ", err)
	}

	if err := <-serveErr; err != http.ErrServerClosed {
		t.Fatalf("serve not returning ErrServerClosed, but %s", err)
	}
}
