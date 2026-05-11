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
