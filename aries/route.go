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

package aries

import (
	"bytes"
	"strings"
)

type routePart struct {
	start, end int
}

type route struct {
	p      string
	parts  []*routePart
	routes []string
	isDir  bool
}

func newRoute(p string) *route {
	if p == "" {
		return new(route)
	}
	w := new(bytes.Buffer)
	n := len(p)
	isDir := p[n-1] == '/'

	splits := strings.Split(p, "/")
	var parts []*routePart
	var routes []string
	for _, s := range splits {
		if len(s) == 0 {
			continue
		}
		w.WriteString("/")
		start := w.Len()
		w.WriteString(s)
		end := w.Len()
		parts = append(parts, &routePart{
			start: start,
			end:   end,
		})
		routes = append(routes, s)
	}

	return &route{
		p:      w.String(),
		parts:  parts,
		routes: routes,
		isDir:  isDir,
	}
}

func (r *route) path() string { return r.p }

func (r *route) size() int { return len(r.routes) }

func (r *route) dir(i int) string {
	if i >= len(r.parts) {
		return r.p
	}
	return r.p[:r.parts[i].start]
}

func (r *route) current(i int) string {
	if i >= len(r.parts) {
		return ""
	}
	part := r.parts[i]
	return r.p[part.start:part.end]
}

func (r *route) rel(i int) string {
	if i >= len(r.parts) {
		return ""
	}
	return r.p[r.parts[i].start:]
}

func (r *route) relRoute(i int) []string {
	if i >= len(r.routes) {
		return nil
	}
	return r.routes[i:]
}
