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

// ConnCounters groups a pair of read/write counters.
type ConnCounters struct {
	Read  *Counter
	Write *Counter
}

// NewConnCounters creates an ConnCounters instance.
func NewConnCounters() *ConnCounters {
	return &ConnCounters{
		Read:  NewCounter(),
		Write: NewCounter(),
	}
}

// Conn wraps net.Conn and counts the number of bytes read/written.
type Conn struct {
	net.Conn
	Counters *ConnCounters
}

// Read reads bytes from underlying net.Conn and count the bytes read.
func (c *Conn) Read(b []byte) (n int, err error) {
	n, err = c.Conn.Read(b)
	c.Counters.Read.Add(int64(n))
	return n, err
}

// Write writes bytes to underlying net.Conn and count the bytes written.
func (c *Conn) Write(b []byte) (n int, err error) {
	n, err = c.Conn.Write(b)
	c.Counters.Write.Add(int64(n))
	return n, err
}

// NewConn creates instances of CountingConn.
func NewConn(conn net.Conn, counters *ConnCounters) net.Conn {
	return &Conn{
		Conn:     conn,
		Counters: counters,
	}
}
