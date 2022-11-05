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

package caco3

type buildContext struct {
	nodes map[string]*buildNode

	built map[string]string // mapping to digests

	cache *buildCache
}

func (c *buildContext) nodeType(n string) string {
	node, ok := c.nodes[n]
	if !ok {
		return ""
	}
	return node.typ
}

func (c *buildContext) ruleType(n string) string {
	node, ok := c.nodes[n]
	if !ok {
		return ""
	}
	return node.ruleType
}
