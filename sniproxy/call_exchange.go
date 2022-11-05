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
	"github.com/gorilla/websocket"
)

// callExchange captures the state of a single call at the client
// transport.
type callExchange struct {
	typ  uint8
	id   uint64
	req  encoderTo
	resp decoderFrom
	err  error
	done func()
}

func newCallExchange(c *transportCall) *callExchange {
	x := &callExchange{
		typ:  c.typ,
		req:  c.req,
		resp: c.resp,
	}
	x.done = func() { c.done(x.err) }
	return x
}

// pendingFetch fetches a pending callExchange based on the call id.
// if no call is found
type pendingFetch struct {
	id   uint64
	call chan *callExchange
}

func sendExchangeReq(conn *websocket.Conn, c *callExchange) error {
	w, err := conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}
	defer w.Close()

	enc := newEncoder(w)
	enc.u64(c.id)
	enc.u8(c.typ)

	if c.req != nil {
		c.req.encodeTo(enc)
	}
	if enc.hasErr() {
		return enc.Err()
	}
	return w.Close()
}
