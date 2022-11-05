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
	"html/template"
)

// Template makes an HTML template with the given body.
func Template(p *Page, body *Node) (*template.Template, error) {
	s, err := RenderString(p, body)
	if err != nil {
		return nil, err
	}
	return template.New("index").Parse(s)
}

// TemplateBody makes an HTML template with the given elements as the body.
func TemplateBody(children ...interface{}) (*template.Template, error) {
	return Template(nil, Body(children...))
}
