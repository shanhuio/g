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

package sniproxy

import (
	"fmt"
	"sync"

	"shanhu.io/pub/errcode"
)

type connections struct {
	closed bool
	m      map[uint64]*connection
	mu     sync.Mutex
}

func newConnections() *connections {
	return &connections{
		m: make(map[uint64]*connection),
	}
}

func (cs *connections) add(c *connection) error {
	id := c.session()

	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.closed {
		return errAlreadyShutdown
	}

	if _, found := cs.m[id]; found {
		return fmt.Errorf("session id conflict: %d", id)
	}
	cs.m[id] = c
	return nil
}

func (cs *connections) get(id uint64) (*connection, error) {
	cs.mu.Lock()
	defer cs.mu.Unlock()
	if cs.closed {
		return nil, errAlreadyShutdown
	}
	c, ok := cs.m[id]
	if !ok {
		return nil, errcode.NotFoundf("session not found: %d", id)
	}
	return c, nil
}

func (cs *connections) remove(id uint64) error {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.closed {
		return errAlreadyShutdown
	}

	if _, found := cs.m[id]; !found {
		return errcode.NotFoundf("session not found: %d", id)
	}
	delete(cs.m, id)
	return nil
}

func (cs *connections) shutdown() map[uint64]*connection {
	cs.mu.Lock()
	defer cs.mu.Unlock()

	if cs.closed {
		return nil
	}
	cs.closed = true
	return cs.m
}
