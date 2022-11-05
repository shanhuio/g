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

package sniproxy

import (
	"context"
	"net"

	"shanhu.io/pub/errcode"
)

type dialer interface {
	dial(ctx context.Context, hello *TLSHelloInfo, asAddr string) (
		net.Conn, error,
	)
}

type tcpDialer struct {
	raddrs map[string]*net.TCPAddr
}

func (d *tcpDialer) dial(_ context.Context, hello *TLSHelloInfo, _ string) (
	net.Conn, error,
) {
	domain := hello.ServerName
	raddr, ok := d.raddrs[domain]
	if !ok {
		return nil, errcode.NotFoundf("no connection for domain %q", domain)
	}

	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		return nil, err
	}
	return conn, nil
}
