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

package httputil

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"shanhu.io/pub/errcode"
)

// Client performs client that calls to a remote server with an optional token.
type Client struct {
	Server *url.URL

	// TokenSource is an optional token source to proivde bearer token.
	TokenSource TokenSource

	UserAgent string // Optional User-Agent for each request.
	Accept    string // Optional Accept header.

	Transport http.RoundTripper
}

func (c *Client) addAuth(req *http.Request) error {
	if c.TokenSource == nil {
		return nil
	}
	ctx := req.Context()
	tok, err := c.TokenSource.Token(ctx, c.Transport)
	if err != nil {
		return errcode.Annotate(err, "get token")
	}
	SetAuthToken(req.Header, tok)
	return nil
}

func (c *Client) doRaw(ctx context.Context, req *http.Request) (
	*http.Response, error,
) {
	return (&http.Client{Transport: c.Transport}).Do(req.WithContext(ctx))
}

func (c *Client) do(ctx context.Context, req *http.Request) (
	*http.Response, error,
) {
	resp, err := c.doRaw(ctx, req)
	if err != nil {
		return nil, err
	}
	if !isSuccess(resp) {
		defer resp.Body.Close()
		return nil, RespError(resp)
	}
	return resp, nil
}

func (c *Client) req(m, p string, r io.Reader) (*http.Request, error) {
	u, err := makeURL(c.Server, p)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(m, u, r)
	if err != nil {
		return nil, err
	}
	if err := c.addAuth(req); err != nil {
		return nil, errcode.Annotate(err, "add auth to request")
	}
	setHeader(req.Header, "User-Agent", c.UserAgent)
	setHeader(req.Header, "Accept", c.Accept)
	return req, nil
}

func (c *Client) reqJSON(m, p string, r io.Reader) (*http.Request, error) {
	req, err := c.req(m, p, r)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Put puts a stream to a path on the server.
func (c *Client) Put(p string, r io.Reader) error {
	return c.PutN(p, r, -1)
}

// PutN puts a stream to a path on the server with content length
// set to n.
func (c *Client) PutN(p string, r io.Reader, n int64) error {
	req, err := c.req(http.MethodPut, p, r)
	if err != nil {
		return err
	}
	if n >= 0 {
		req.ContentLength = n
	}
	resp, err := c.do(context.TODO(), req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// PutBytes puts bytes to a path on the server.
func (c *Client) PutBytes(p string, bs []byte) error {
	return c.PutN(p, bytes.NewBuffer(bs), int64(len(bs)))
}

// JSONPut puts an object in JSON encoding.
func (c *Client) JSONPut(p string, v interface{}) error {
	bs, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.PutBytes(p, bs)
}

func (c *Client) poke(ctx context.Context, m, p string) error {
	req, err := c.req(m, p, nil)
	if err != nil {
		return err
	}
	resp, err := c.do(ctx, req)
	if err != nil {
		return err
	}
	return resp.Body.Close()
}

// GetCode gets a response from a route and returns the
// status code.
func (c *Client) GetCode(p string) (int, error) {
	req, err := c.req(http.MethodGet, p, nil)
	if err != nil {
		return 0, err
	}
	resp, err := c.doRaw(context.TODO(), req)
	if err != nil {
		return 0, err
	}
	code := resp.StatusCode
	resp.Body.Close()
	return code, nil
}

// Poke posts a signal to the given route on the server.
func (c *Client) Poke(p string) error {
	return c.poke(context.TODO(), http.MethodPost, p)
}

// Get gets a response from a route on the server.
func (c *Client) Get(p string) (*http.Response, error) {
	req, err := c.req(http.MethodGet, p, nil)
	if err != nil {
		return nil, err
	}
	return c.do(context.TODO(), req)
}

// GetString gets the string response from a route on the server.
func (c *Client) GetString(p string) (string, error) {
	resp, err := c.Get(p)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	return respString(resp)
}

// GetInto gets the specified path and writes everything from the body to the
// given writer.
func (c *Client) GetInto(p string, w io.Writer) (int64, error) {
	resp, err := c.Get(p)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()
	return io.Copy(w, resp.Body)
}

// GetBytes gets the byte array from a route on the server.
func (c *Client) GetBytes(p string) ([]byte, error) {
	resp, err := c.Get(p)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

// JSONGet gets the content of a path and decodes the response
// into resp as JSON.
func (c *Client) JSONGet(p string, resp interface{}) error {
	req, err := c.reqJSON(http.MethodGet, p, nil)
	if err != nil {
		return nil
	}
	httpResp, err := c.do(context.TODO(), req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(resp); err != nil {
		return err
	}
	return httpResp.Body.Close()
}

// Post posts with request body from r, and copies the response body
// to w.
func (c *Client) Post(p string, r io.Reader, w io.Writer) error {
	if r != nil {
		r = io.NopCloser(r)
	}
	req, err := c.req(http.MethodPost, p, r)
	if err != nil {
		return err
	}
	resp, err := c.do(context.TODO(), req)
	if err != nil {
		return err
	}
	return copyRespBody(resp, w)
}

func (c *Client) jsonPost(ctx context.Context, p string, req interface{}) (
	*http.Response, error,
) {
	bs, err := json.Marshal(req)
	if err != nil {
		return nil, errcode.Annotate(err, "marshal request")
	}
	httpReq, err := c.reqJSON(http.MethodPost, p, bytes.NewBuffer(bs))
	if err != nil {
		return nil, err
	}
	return c.do(ctx, httpReq)
}

// JSONPost posts a JSON object as the request body and writes the body
// into the given writer.
func (c *Client) JSONPost(p string, req interface{}, w io.Writer) error {
	resp, err := c.jsonPost(context.TODO(), p, req)
	if err != nil {
		return err
	}
	return copyRespBody(resp, w)
}

// CallContext performs a call with the request as a marshalled JSON object,
// and the response unmarshalled as a JSON object.
func (c *Client) CallContext(
	ctx context.Context, p string, req, resp interface{},
) error {
	httpResp, err := c.jsonPost(ctx, p, req)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if resp == nil {
		return nil
	}
	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(resp); err != nil {
		return errcode.Annotate(err, "decode response")
	}
	return httpResp.Body.Close()
}

// Call performs a CallContext with context.TODO().
func (c *Client) Call(p string, req, resp interface{}) error {
	return c.CallContext(context.TODO(), p, req, resp)
}

// JSONCall is an alias to Call.
func (c *Client) JSONCall(p string, req, resp interface{}) error {
	return c.Call(p, req, resp)
}

// Delete sends a delete message to the particular path.
func (c *Client) Delete(p string) error {
	return c.poke(context.TODO(), http.MethodDelete, p)
}
