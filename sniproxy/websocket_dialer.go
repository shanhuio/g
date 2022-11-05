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
	"encoding/json"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
)

type websocketDialer struct {
	url    *url.URL
	token  string
	dialer *websocket.Dialer
}

func (d *websocketDialer) address() string {
	return d.url.String()
}

func (d *websocketDialer) getDialer() *websocket.Dialer {
	if d.dialer != nil {
		return d.dialer
	}
	return &websocket.Dialer{
		ReadBufferSize:  DefaultReadBufferSize,
		WriteBufferSize: DefaultWriteBufferSize,
	}
}

func (d *websocketDialer) dial(ctx context.Context, opt *Options) (
	*websocket.Conn, error,
) {
	header := make(http.Header)
	if d.token != "" {
		httputil.SetAuthToken(header, d.token)
	}

	u := *d.url
	q := u.Query()
	if opt == nil {
		opt = new(Options)
	}
	optBytes, err := json.Marshal(opt)
	if err != nil {
		return nil, errcode.Annotate(err, "marshal options")
	}
	q.Set("opt", string(optBytes))

	u.RawQuery = q.Encode()

	dialer := d.getDialer()
	conn, _, err := dialer.DialContext(ctx, u.String(), header)
	return conn, err
}

func (d *websocketDialer) dialSide(
	ctx context.Context, tok string, k *sessionKey,
) (*websocket.Conn, error) {
	side, err := k.encode()
	if err != nil {
		return nil, errcode.Annotate(err, "encode session key")
	}

	header := make(http.Header)
	if tok != "" {
		httputil.SetAuthToken(header, tok)
	}

	u := *d.url
	q := u.Query()
	q.Set("side", side)
	u.RawQuery = q.Encode()

	dialer := d.getDialer()
	conn, _, err := dialer.DialContext(ctx, u.String(), header)
	return conn, err
}
