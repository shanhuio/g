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

package webgen

import (
	"golang.org/x/net/html"
	"shanhu.io/pub/errcode"
)

// Node wraps around an html node.
type Node struct{ *html.Node }

// Add appends more stuff into the node.
func (n *Node) Add(children ...interface{}) error {
	return addChildren(n.Node, children...)
}

func text(s string) *html.Node {
	return &html.Node{
		Type: html.TextNode,
		Data: s,
	}
}

// Text creates a text node.
func Text(s string) *Node { return &Node{text(s)} }

func addChildren(n *html.Node, children ...interface{}) error {
	for _, child := range children {
		switch c := child.(type) {
		case Class:
			setClass(n, c)
		case Attrs:
			setAttrs(n, c)
		case string:
			n.AppendChild(text(c))
		case *html.Node:
			n.AppendChild(c)
		case *Node:
			n.AppendChild(c.Node)
		default:
			return errcode.Internalf("unknown child type: %T", child)
		}
	}
	return nil
}
