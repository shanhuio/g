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

// Package srpc provides the Shanhu RPC caller for performing
// program-to-program, service-to-service or CLI-to-service interactions.
// It can also be used for Web-to-service AJAX-based RPC calls.
// The transport interface is intentionally limtied, and it is not RESTful.
// So this is not to be used to interact with a typical RESTful interface
// that is provided by third-party or other services.
package srpc

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"shanhu.io/pub/errcode"
)

// Caller is an RPC caller that can call
type Caller struct {
	server *url.URL
	client *http.Client

	tokener Tokener
}

// NewCaller returns a caller that calls to the specific URL.
func NewCaller(server *url.URL) *Caller {
	return NewTokenCaller(server, nil)
}

// NewTokenCaller returns a caller that calls to the specific URL with the given
// tokener to provide auth tokens.
func NewTokenCaller(server *url.URL, tokener Tokener) *Caller {
	return &Caller{
		server:  server,
		client:  http.DefaultClient,
		tokener: tokener,
	}
}

// Call performs a JSON RPC call on the specified method path.
func (c *Caller) Call(
	ctx context.Context, p string, req, resp interface{},
) error {
	var token *string
	if c.tokener != nil {
		t, err := c.tokener.Token(ctx)
		if err != nil {
			return errcode.Annotate(err, "get auth token")
		}
		token = &t
	}

	reqBytes, err := json.Marshal(req)
	if err != nil {
		return errcode.Annotate(err, "marshal request")
	}

	u := urlJoin(c.server, p)
	getBody := func() (io.ReadCloser, error) {
		return io.NopCloser(bytes.NewReader(reqBytes)), nil
	}
	reqBody, _ := getBody()

	httpReq := (&http.Request{
		Method:     http.MethodPost,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		URL:        u,
		Header:     make(http.Header),
		Body:       reqBody,
		GetBody:    getBody,

		ContentLength: int64(len(reqBytes)),
	}).WithContext(ctx)

	httpReq.Header.Set("Content-Type", contentTypeJSON)
	httpReq.Header.Set("Accept", contentTypeJSON)
	if token != nil {
		httpReq.Header.Set("Authorization", "Bearer "+*token)
	}

	client := http.DefaultClient
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if !isSuccessStatus(httpResp.StatusCode) {
		return respError(httpResp)
	}

	if t := httpResp.Header.Get("Content-Type"); t != contentTypeJSON {
		return fmt.Errorf("unexpected content type: %q", t)
	}

	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(resp); err != nil {
		return errcode.Annotate(err, "unmarshal response")
	}

	return nil
}
