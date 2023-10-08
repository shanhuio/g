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
	mrand "math/rand"
	"net"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"shanhu.io/g/errcode"
	"shanhu.io/g/rand"
)

type endpointClient struct {
	conn     *websocket.Conn
	options  *Options
	tr       *transport
	sessions *sync.Map
	ids      *sessionID // only used in side dialing
	rand     *mrand.Rand
	office   *connMailOffice

	token func() (string, error)
}

func newEndpointClient(conn *websocket.Conn, opt *Options) *endpointClient {
	return &endpointClient{
		conn:     conn,
		options:  opt,
		tr:       newTransport(conn, opt),
		sessions: new(sync.Map),
		ids:      newSessionID(),
		rand:     rand.New(),
		office:   newConnMailOffice(),
	}
}

func (c *endpointClient) serve() error {
	return c.tr.serve()
}

func (c *endpointClient) setToken(f func() (string, error)) {
	c.token = f
}

func (c *endpointClient) Hello(ctx context.Context, msg string) (
	string, error,
) {
	req := &helloResponse{msg: msg}
	resp := new(helloResponse)
	if err := c.tr.call(ctx, msgHello, req, resp); err != nil {
		return "", err
	}
	return resp.msg, nil
}

func (c *endpointClient) Dial(
	ctx context.Context, asAddr string,
) (net.Conn, error) {
	if !c.options.Siding {
		req := &dialRequest{}
		resp := new(dialResponse)
		if err := c.tr.call(ctx, msgDial, req, resp); err != nil {
			return nil, err
		}
		if resp.err != nil {
			return nil, resp.err
		}
		return newTunnel(c.tr, resp.session), nil
	}

	// side connection does not create a tunnel.
	token, err := c.token()
	if err != nil {
		return nil, errcode.Annotate(err, "get side token")
	}
	k := &sessionKey{
		ID:  c.ids.next(),
		Key: c.rand.Uint64(),
	}
	box := c.office.newBox(k)
	defer box.cleanUp()

	resp := new(dialResponse)
	if c.options.DialWithAddr {
		req := &dialSide2Request{
			session: k.ID,
			key:     k.Key,
			token:   token,
			tcpAddr: asAddr,
		}
		if err := c.tr.call(ctx, msgDialSide2, req, resp); err != nil {
			return nil, err
		}
	} else {
		req := &dialSideRequest{
			session: k.ID,
			key:     k.Key,
			token:   token,
		}
		if err := c.tr.call(ctx, msgDialSide, req, resp); err != nil {
			return nil, err
		}
	}
	if resp.err != nil {
		return nil, resp.err
	}
	return box.receive(ctx)
}

func (c *endpointClient) deliverSideConn(k *sessionKey, conn net.Conn) error {
	return c.office.deliver(k, conn)
}

func (c *endpointClient) Close() error {
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()
	err := c.tr.shutdown(ctx)
	c.conn.Close() // Always close the conn.
	return err
}
