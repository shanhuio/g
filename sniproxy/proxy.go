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
	"context"
	"errors"
	"io"
	"log"
	"net"
	"strings"
	"sync"

	"shanhu.io/pub/netutil"
)

type proxy struct {
	dialer dialer
}

func newProxy(d dialer) *proxy {
	return &proxy{dialer: d}
}

var errNameRejected = errors.New("name rejected")

func isRejectedDomain(name string) bool {
	if name == "" {
		return true
	}
	if ip := net.ParseIP(name); ip != nil {
		return true
	}
	for _, suf := range []string{
		// Telegram proxy domains. Not real ones, but uses port 443 with
		// TLS-like protocol.
		".iproxy.cloud",
		".after.blue",
		".spothot.online",
		".speedy.red",
	} {
		if strings.HasSuffix(name, suf) {
			return true
		}
	}
	return false
}

func (p *proxy) hostConn(ctx context.Context, conn net.Conn) error {
	defer conn.Close()
	bc := NewTLSHelloConn(conn)
	hello, err := bc.HelloInfo()
	if err != nil {
		return err
	}
	if isRejectedDomain(hello.ServerName) {
		return errNameRejected
	}

	addr := conn.RemoteAddr().String()
	remote, err := p.dialer.dial(ctx, hello, addr)
	if err != nil {
		return err
	}
	closer := &closerOnce{Closer: remote}
	defer closer.Close()

	return netutil.JoinConn(ctx, remote, bc)
}

// IsClosedConnError checks if the error is a Closed connection error.
func IsClosedConnError(err error) bool {
	return errors.Is(err, net.ErrClosed)
}

func isCommonProxyError(err error) bool {
	if IsClosedConnError(err) {
		return true
	}
	for _, this := range []error{
		context.Canceled,
		io.ErrUnexpectedEOF,
		io.ErrClosedPipe,
		errNameRejected,
		errAlreadyShutdown,
		errAlreadyClosed,
	} {
		if this == err {
			return true
		}
	}

	return false
}

func (p *proxy) serve(ctx context.Context, lis net.Listener) error {
	var wg sync.WaitGroup
	defer wg.Wait()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		lis.Close() // to unblock lis.Accept
	}()

	for {
		conn, err := lis.Accept()
		if err != nil {
			if IsClosedConnError(err) {
				return nil
			}
			return err
		}
		wg.Add(1)
		go func(ctx context.Context) {
			defer wg.Done()
			if err := p.hostConn(ctx, conn); err != nil {
				if !isCommonProxyError(err) {
					log.Print("proxy connection error: ", err)
				}
			}
		}(ctx)
	}
}
