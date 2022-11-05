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

package dock

import (
	"io"
	"net/http"
	"net/url"

	"shanhu.io/pub/httputil"
)

// Socket is the default socket location.
const Socket = "/var/run/docker.sock"

type emptyReader struct{}

func (emptyReader) Read([]byte) (int, error) { return 0, io.EOF }

// Client is a docker daemon client that can be used to issue
// docker commands.
type Client struct {
	client *httputil.Client
}

// NewClient creates a new client using the given httputil.Client
func NewClient(c *httputil.Client) *Client {
	return &Client{client: c}
}

// NewUnixClient creates a new unix domain socket client.
// When sock is empty, "/var/run/docker.sock" is used.
func NewUnixClient(sock string) *Client {
	if sock == "" {
		sock = Socket
	}
	return NewClient(httputil.NewUnixClient(sock))
}

func (c *Client) call(
	p string, q url.Values, req, resp interface{},
) error {
	return c.jsonCall(p, q, req, resp)
}

func (c *Client) jsonCall(
	p string, q url.Values, req, resp interface{},
) error {
	u := apiURLQuery(p, q)
	return c.client.Call(u, req, resp)
}

func (c *Client) jsonPost(
	p string, q url.Values, req interface{}, w io.Writer,
) error {
	u := apiURLQuery(p, q)
	return c.client.JSONPost(u, req, w)
}

func (c *Client) jsonGet(p string, q url.Values, resp interface{}) error {
	u := apiURLQuery(p, q)
	return c.client.JSONGet(u, resp)
}

func (c *Client) post(
	p string, q url.Values, r io.Reader, w io.Writer,
) error {
	u := apiURLQuery(p, q)
	if r == nil {
		r = emptyReader{}
	}
	return c.client.Post(u, r, w)
}

func (c *Client) del(p string, q url.Values) error {
	return c.client.Delete(apiURLQuery(p, q))
}

func (c *Client) poke(p string, q url.Values) error {
	return c.client.Poke(apiURLQuery(p, q))
}

func (c *Client) put(p string, q url.Values, r io.Reader) error {
	u := apiURLQuery(p, q)
	return c.client.Put(u, io.NopCloser(r))
}

func (c *Client) get(p string, q url.Values) (*http.Response, error) {
	return c.client.Get(apiURLQuery(p, q))
}

func (c *Client) getInto(p string, q url.Values, w io.Writer) (int64, error) {
	return c.client.GetInto(apiURLQuery(p, q), w)
}
