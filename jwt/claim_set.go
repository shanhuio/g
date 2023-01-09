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

package jwt

import (
	"encoding/json"
	"strings"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/strutil"
)

// ClaimSet contains the JWT claims
type ClaimSet struct {
	Iss   string `json:"iss"`   // Issuer.
	Scope string `json:"scope"` // Scope, space-delimited list.
	Aud   string `json:"aud"`   // Audiance. Intended target.
	Exp   int64  `json:"exp"`   // Expiration time (Unix timestamp seconds)
	Iat   int64  `json:"iat"`   // Asserstion time (Unix timestamp seconds)
	Typ   string `json:"typ"`   // Token type.

	Sub string `json:"sub"`

	Extra map[string]interface{} `json:"-"`
}

// ExtraString reads an extra string field from the claim set.
func (c *ClaimSet) ExtraString(k string) (string, bool) {
	if len(c.Extra) == 0 {
		return "", false
	}
	v, ok := c.Extra[k]
	if !ok {
		return "", false
	}
	s, ok := v.(string)
	if !ok {
		return "", false
	}
	return s, true
}

func (c *ClaimSet) encode() (string, error) {
	m := make(map[string]interface{})

	for _, entry := range []struct {
		k, v     string
		mustHave bool
	}{
		{k: "iss", v: c.Iss, mustHave: true},
		{k: "scope", v: c.Scope},
		{k: "aud", v: c.Aud, mustHave: true},
		{k: "typ", v: c.Typ},
		{k: "sub", v: c.Sub},
	} {
		if entry.mustHave || entry.v != "" {
			m[entry.k] = entry.v
		}
	}

	m["exp"] = c.Exp
	m["iat"] = c.Iat

	for k, v := range c.Extra {
		m[k] = v
	}

	return encodeSegment(m)
}

func decodeClaimSet(s string) (*ClaimSet, error) {
	bs, err := decodeSegmentBytes(s)
	if err != nil {
		return nil, err
	}

	c := new(ClaimSet)
	if err := json.Unmarshal(bs, c); err != nil {
		return nil, err
	}
	m := make(map[string]interface{})
	if err := json.Unmarshal(bs, &m); err != nil {
		return nil, err
	}

	for _, k := range []string{
		"iss", "scope", "aud", "exp", "iat", "typ", "sub",
	} {
		delete(m, k)
	}
	if len(m) > 0 {
		c.Extra = m
	}
	return c, nil
}

// CheckClaimSet checks claims in claim set, see if it matches the values
// in the template.
func CheckClaimSet(claims, tmpl *ClaimSet) error {
	if claims == nil {
		return errcode.Unauthorizedf("claims missing")
	}
	if tmpl == nil { // nothing to check.
		return nil
	}
	if tmpl.Iss != "" {
		if claims.Iss != tmpl.Iss {
			return errcode.Unauthorizedf("wrong issuer")
		}
	}
	if tmpl.Aud != "" {
		if claims.Aud != tmpl.Aud {
			return errcode.Unauthorizedf("wrong audiance")
		}
	}
	if tmpl.Typ != "" {
		if claims.Typ != tmpl.Typ {
			return errcode.Unauthorizedf("wrong type")
		}
	}
	if tmpl.Sub != "" {
		if claims.Sub != tmpl.Sub {
			return errcode.Unauthorizedf("wrong subject")
		}
	}
	if tmpl.Scope != "" {
		tmplScopes := strings.Fields(tmpl.Scope)
		claimScopesSet := strutil.MakeSet(strings.Fields(claims.Scope))
		for _, s := range tmplScopes {
			if _, ok := claimScopesSet[s]; !ok {
				return errcode.Unauthorizedf("scope %q missing", s)
			}
		}
	}

	return nil
}
