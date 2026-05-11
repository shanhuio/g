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
