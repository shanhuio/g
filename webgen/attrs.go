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

package webgen

import (
	"sort"
	"strings"

	"golang.org/x/net/html"
)

// Attrs is an attribute map.
type Attrs map[string]string

func setAttrs(node *html.Node, attrs Attrs) {
	var keys []string
	for k := range attrs {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		node.Attr = append(node.Attr, html.Attribute{
			Key: k,
			Val: attrs[k],
		})
	}
	return
}

// Class is the class attribute for a div.
type Class []string

func setClass(node *html.Node, cls Class) {
	if len(cls) == 0 {
		return
	}

	node.Attr = append(node.Attr, html.Attribute{
		Key: "class",
		Val: strings.Join([]string(cls), " "),
	})
}
