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
	"net/url"

	"encoding/json"
	"shanhu.io/pub/errcode"
)

// Options is a the JSON marshalable options for dialing an endpoint.
type Options struct {
	// Using a new websocket connection for each new incoming
	// new connection.
	Siding bool `json:",omitempty"`

	// Remote enables sending remote address.
	DialWithAddr bool `json:",omitempty"`
}

func decodeOptions(s string) (*Options, error) {
	opt := new(Options)
	if s == "" {
		return opt, nil
	}
	if err := json.Unmarshal([]byte(s), opt); err != nil {
		return nil, errcode.InvalidArgf("invalid options: %s", err)
	}

	return opt, nil
}

func optionsFromQuery(q url.Values) (*Options, error) {
	return decodeOptions(q.Get("opt"))
}
