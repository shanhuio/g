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
	"net/http"
	"net/url"
	"path"

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
		client:  &http.Client{},
		tokener: tokener,
	}
}

const contentTypeJSON = "application/json"

// Call performs a JSON RPC call on the specified method path.
func (c *Caller) Call(
	ctx context.Context, p string, req, resp interface{},
) error {
	reqBytes, err := json.Marshal(req)
	if err != nil {
		return errcode.Annotate(err, "marshal request")
	}
	reqBody := bytes.NewReader(reqBytes)

	u := *c.server
	u.Path = path.Join(u.Path, p)

	httpReq, err := http.NewRequestWithContext(
		ctx, http.MethodPost, u.String(), reqBody,
	)
	if err != nil {
		return errcode.Annotate(err, "make http request")
	}

	// Content-Length will be already set by NewRequestWithContext.

	httpReq.Header.Set("Content-Type", contentTypeJSON)
	httpReq.Header.Set("Accept", contentTypeJSON)

	httpResp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer httpResp.Body.Close()

	if !isSuccessStatus(httpResp.StatusCode) {
		return respError(httpResp)
	}

	dec := json.NewDecoder(httpResp.Body)
	if err := dec.Decode(resp); err != nil {
		return errcode.Annotate(err, "unmarshal response")
	}

	return nil
}
