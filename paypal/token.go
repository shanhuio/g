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

package paypal

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"shanhu.io/pub/errcode"
)

func makeTokenRequest(host, id, secret string) *http.Request {
	u := &url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "/v1/oauth2/token",
		User:   url.UserPassword(id, secret),
	}

	form := make(url.Values)
	form.Set("grant_type", "client_credentials")

	header := make(http.Header)
	header.Add("Accept", "application/json")
	header.Add("Accept-Language", "eu_US")

	return &http.Request{
		Method: http.MethodPost,
		URL:    u,
		Header: header,
		Body:   io.NopCloser(strings.NewReader(form.Encode())),
	}
}

func tokenFromResponse(r io.Reader) (string, error) {
	body, err := io.ReadAll(r)
	if err != nil {
		return "", errcode.Annotate(err, "read token response")
	}

	var dat struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.Unmarshal(body, &dat); err != nil {
		return "", errcode.Annotate(err, "parse token from response")
	}

	return dat.AccessToken, nil
}
