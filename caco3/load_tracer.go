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

type loadTracer struct {
	trace []string
	m     map[string]bool
}

func newLoadTracer() *loadTracer {
	return &loadTracer{
		m: make(map[string]bool),
	}
}

func (t *loadTracer) push(name string) bool {
	if t.m[name] {
		return false
	}
	t.trace = append(t.trace, name)
	t.m[name] = true
	return true
}

func (t *loadTracer) pop() {
	n := len(t.trace)
	if n == 0 {
		return
	}
	last := t.trace[n-1]
	delete(t.m, last)
	t.trace = t.trace[:n-1]
}

func (t *loadTracer) stack() []string {
	return t.trace
}
