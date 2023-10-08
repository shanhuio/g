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

package caco3

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"shanhu.io/g/errcode"
)

// buildAction is a structure for creating the digest of the execution of a
// rule.
type buildAction struct {
	Rule     string `json:",omitempty"` // Digest of the rule functor.
	RuleType string `json:",omitempty"`

	// Map of dependency names to their digiests.
	Deps map[string]string `json:",omitempty"`

	Outs      []string `json:",omitempty"`
	DockerOut bool     `json:",omitempty"`

	OutputOf string `json:",omitempty"` // Get the output from a rule.
}

func makeDigest(t, name string, v interface{}) (string, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintln(buf, t)
	fmt.Fprintln(buf, name)
	bs, err := json.Marshal(v)
	if err != nil {
		return "", errcode.Annotate(err, "json marshal")
	}
	buf.Write(bs)
	sum := sha256.Sum256(buf.Bytes())
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}
