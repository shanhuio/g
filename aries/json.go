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

package aries

import (
	"encoding/json"
	"log"

	"shanhu.io/pub/errcode"
)

// ReplyJSON replies a JSON marshable object over the response.
func ReplyJSON(c *C, v interface{}) error {
	c.Resp.Header().Set("Content-Type", "application/json")

	bs, err := json.Marshal(v)
	if err != nil {
		return errcode.Internalf("response encode error: %s", err)
	}

	if _, err := c.Resp.Write(bs); err != nil {
		log.Println(err)
	}
	return nil
}

// UnmarshalJSONBody unmarshals the body into a JSON object.
func UnmarshalJSONBody(c *C, v interface{}) error {
	dec := json.NewDecoder(c.Req.Body)
	if err := dec.Decode(v); err != nil {
		return errcode.Add(errcode.InvalidArg, err)
	}
	return nil
}

// PrintJSON replies a JSON marshaable object over the reponse with
// pretty printing.
func PrintJSON(c *C, v interface{}) error {
	bs, err := json.MarshalIndent(v, "", "    ")
	if err != nil {
		return errcode.Internalf("response encode error: %s", err)
	}

	if _, err := c.Resp.Write(bs); err != nil {
		log.Println(err)
	}
	return nil
}
