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

package httpstest

import (
	"context"
	"net"
)

// SinkDialFunc returns a dialing function that always dials to the same
// address.
func SinkDialFunc(sinkAddr string) func(
	ctx context.Context, net, addr string,
) (net.Conn, error) {
	return sink(sinkAddr)
}

func sink(sinkAddr string) func(
	ctx context.Context, net, addr string,
) (net.Conn, error) {
	d := new(net.Dialer)
	return func(ctx context.Context, net, addr string) (net.Conn, error) {
		return d.DialContext(ctx, net, sinkAddr)
	}
}

func sinkHTTPS(httpAddr, httpsAddr string) func(
	ctx context.Context, net, addr string,
) (net.Conn, error) {
	d := new(net.Dialer)
	return func(ctx context.Context, netStr, addr string) (net.Conn, error) {
		_, port, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, err
		}
		sinkAddr := httpAddr
		if port == "443" || port == "https" {
			sinkAddr = httpsAddr
		}
		return d.DialContext(ctx, netStr, sinkAddr)
	}
}
