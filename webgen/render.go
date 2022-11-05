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
	"bytes"
	"io"

	"golang.org/x/net/html"
)

// Page contains the configuration of a page.
type Page struct {
	NoDocType bool
	Title     string
}

// HTMLDocString is the doc type string for HTML.
const HTMLDocString = "<!doctype html>"

// Render renders a page with the given HTML
func Render(w io.Writer, n *Node) error {
	return html.Render(w, n.Node)
}

// RenderBody renders a page with the given Body
func RenderBody(w io.Writer, page *Page, body *Node) error {
	if page == nil {
		page = new(Page) // just use empty value for default.
	}

	if !page.NoDocType {
		if _, err := io.WriteString(w, HTMLDocString+"\n"); err != nil {
			return err
		}
	}

	doc := NewHTMLEnglish()
	head := Head(NewMeta("charset", "UTF-8"))
	if page.Title != "" {
		if err := head.Add(Title(page.Title)); err != nil {
			return err
		}
	}
	if err := doc.Add(head, body); err != nil {
		return err
	}

	if err := html.Render(w, doc.Node); err != nil {
		return err
	}
	_, err := io.WriteString(w, "\n")
	return err
}

// RenderString renders a page into a string.
func RenderString(page *Page, body *Node) (string, error) {
	buf := new(bytes.Buffer)
	if err := RenderBody(buf, page, body); err != nil {
		return "", err
	}
	return buf.String(), nil
}
