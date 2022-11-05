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
	"fmt"

	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

// Element create a new element.
func Element(n interface{}, children ...interface{}) *Node {
	var ret *html.Node
	switch n := n.(type) {
	case atom.Atom:
		ret = &html.Node{
			Type:     html.ElementNode,
			DataAtom: n,
			Data:     n.String(),
		}
	case string:
		ret = &html.Node{
			Type: html.ElementNode,
			Data: n,
		}
	default:
		panic(fmt.Sprintf("unknown element %T", n))
	}

	if err := addChildren(ret, children...); err != nil {
		panic(err)
	}
	return &Node{ret}
}

func bind(a atom.Atom) func(c ...interface{}) *Node {
	return func(c ...interface{}) *Node { return Element(a, c...) }
}

// Shorthand element creators.
var (
	HTML = bind(atom.Html)

	Head  = bind(atom.Head)
	Title = bind(atom.Title)
	Meta  = bind(atom.Meta)
	Link  = bind(atom.Link)

	Body = bind(atom.Body)
	Div  = bind(atom.Div)
	Span = bind(atom.Span)

	P          = bind(atom.P)
	Pre        = bind(atom.Pre)
	Blockquote = bind(atom.Blockquote)
	Strong     = bind(atom.Strong)
	Em         = bind(atom.Em)

	H1 = bind(atom.H1)
	H2 = bind(atom.H2)
	H3 = bind(atom.H3)
	H4 = bind(atom.H3)
	H5 = bind(atom.H3)
	H6 = bind(atom.H3)

	A = bind(atom.A)

	Ul = bind(atom.Ul)
	Ol = bind(atom.Ol)
	Li = bind(atom.Li)

	Br = bind(atom.Br)
)
