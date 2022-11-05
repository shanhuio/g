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

package counting

import (
	"net"
)

// Listener is a net.Listener that tracks the number of bytes flown over
// the underlying listener.
type Listener struct {
	net.Listener
	Counters *ConnCounters
}

// NewListener creates a counting.Listener instance.
func NewListener(lis net.Listener, c *ConnCounters) *Listener {
	return &Listener{
		Listener: lis,
		Counters: c,
	}
}

// Accept wrapa incoming connection in a counting.Conn.
func (l *Listener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return nil, err
	}
	return NewConn(conn, l.Counters), nil
}

// WrapListener wraps the listener with a pair of ConnCounters.
// If c is nil, it returns lis directly.
func WrapListener(lis net.Listener, c *ConnCounters) net.Listener {
	if c == nil {
		return lis
	}
	return NewListener(lis, c)
}
