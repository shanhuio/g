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
	"html/template"
	"path/filepath"

	"shanhu.io/pub/strutil"
)

// Templates is a collection of templates.
type Templates struct {
	path   string
	logger *Logger
}

// DefaultTemplatePath is the default template path.
const DefaultTemplatePath = "_/tmpl"

// NewTemplates creates a collection of templates in a particular folder.
func NewTemplates(p string, logger *Logger) *Templates {
	if logger == nil {
		logger = StdLogger()
	}
	p = strutil.Default(p, DefaultTemplatePath)
	return &Templates{path: p, logger: logger}
}

func (ts *Templates) tmpl(f string) string {
	return filepath.Join(ts.path, f)
}

// TemplatesJSON tells the Templates to print JSON data rather than
// render the template.
const TemplatesJSON = "!JSON"

// Serve serves a data page using a particular template.
func (ts *Templates) Serve(c *C, p string, dat interface{}) error {
	if ts.path == TemplatesJSON {
		return ReplyJSON(c, dat)
	}
	t, err := template.ParseFiles(ts.tmpl(p))
	if err != nil {
		Log(ts.logger, err.Error())
		return NotFound
	}
	if err := t.Execute(c.Resp, dat); err != nil {
		Log(ts.logger, err.Error())
	}
	return nil
}
