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

// Package tasks provides a way to implement a set of simple subcommand-like
// API for user to invoke. Although a task is a POST request, the request often
// has no body content.
package tasks

import (
	"net/http"
	"sort"

	"shanhu.io/pub/aries"
	"shanhu.io/pub/errcode"
)

type tasks struct {
	lst []string
	m   map[string]aries.Service
}

func newTasks(m map[string]aries.Service) *tasks {
	var lst []string
	for name := range m {
		lst = append(lst, name)
	}
	sort.Strings(lst)
	return &tasks{
		lst: lst,
		m:   m,
	}
}

func (t *tasks) serve(c *aries.C) error {
	name := c.Rel()

	if c.Req.Method != http.MethodPost {
		return errcode.NotFoundf("task %q must use POST", name)
	}

	f, found := t.m[name]
	if !found {
		if name == "help" {
			return aries.ReplyJSON(c, t.lst)
		}
		return errcode.InvalidArgf("unknown task: %q", name)
	}
	return f.Serve(c)
}

// Serve returns the serving function for a task list.
func Serve(tasks map[string]aries.Service) aries.Func {
	return newTasks(tasks).serve
}
