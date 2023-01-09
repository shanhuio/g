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
	"fmt"
	"strings"
)

// Mux is a router for a given context
type Mux struct {
	exacts   map[string]Func
	prefixes map[string]Func
	t        *trieNode
}

// NewMux creates a new mux for the incoming request.
func NewMux() *Mux {
	return &Mux{
		t:        newTrieRoot(),
		prefixes: make(map[string]Func),
		exacts:   make(map[string]Func),
	}
}

// Prefix adds a prefix matching rule.
func (m *Mux) Prefix(s string, f Func) error {
	if !m.t.add(s) {
		return fmt.Errorf("duplicate prefix %q", s)
	}
	m.prefixes[s] = f
	return nil
}

// Exact adds an exact matching rule.
func (m *Mux) Exact(s string, f Func) error {
	_, ok := m.exacts[s]
	if ok {
		return fmt.Errorf("duplicate exact %q", s)
	}
	m.exacts[s] = f
	return nil
}

// Dir add is a shortcut of Exact(s) and Prefix(s + "/").
func (m *Mux) Dir(s string, f Func) error {
	if s == "/" {
		if err := m.Exact(s, f); err != nil {
			return err
		}
		return m.Prefix(s, f)
	}

	s = strings.TrimSuffix(s, "/")
	if err := m.Exact(s, f); err != nil {
		return err
	}
	return m.Prefix(s+"/", f)
}

// Route returns the serving function for the given context.
func (m *Mux) Route(c *C) Func {
	if f, ok := m.exacts[c.Path]; ok {
		return f
	}
	s, _ := trieFind(m.t, c.Path)
	if f, ok := m.prefixes[s]; ok {
		return f
	}
	return nil
}

// Serve serves an incoming request based on c.Path.
// It returns true when it hits something.
// And it returns false when it hits nothing.
func (m *Mux) Serve(c *C) error {
	f := m.Route(c)
	if f == nil {
		return Miss
	}
	return f(c)
}
