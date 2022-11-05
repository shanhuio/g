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
