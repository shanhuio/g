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

type endpointExchange struct {
	id      uint64
	t       uint8
	errcode uint8
	req     decoderFrom
	resp    encoderTo
}

func (x *endpointExchange) encodeTo(enc *encoder) {
	enc.u64(x.id)
	enc.u8(x.t)
	enc.u8(x.errcode)

	if x.resp != nil {
		x.resp.encodeTo(enc)
	}
}

func writeExchangeResp(conn *websocket.Conn, x *endpointExchange) error {
	w, err := conn.NextWriter(websocket.BinaryMessage)
	if err != nil {
		return err
	}

	enc := newEncoder(w)
	x.encodeTo(enc)
	if enc.hasErr() {
		return enc.Err()
	}
	return w.Close()
}
