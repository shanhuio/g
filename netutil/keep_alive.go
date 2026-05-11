package netutil

import (
	"net"
	"time"
)

type keepAliveListener struct {
	*net.TCPListener
}

func (ln keepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}

	tc.SetKeepAlive(true)
	tc.SetKeepAlivePeriod(3 * time.Minute)
	return tc, nil
}

// WrapKeepAlive wraps the listener. If the listener is a TCP listener, it
// sets keep alive to 3 minute.
func WrapKeepAlive(ln net.Listener) net.Listener {
	tcpLis, ok := ln.(*net.TCPListener)
	if !ok {
		return ln
	}
	return keepAliveListener{tcpLis}
}
