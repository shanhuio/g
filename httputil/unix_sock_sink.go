package httputil

import (
	"context"
	"net"
	"net/http"
)

func unixSockSink(sinkAddr string) func(
	ctx context.Context, net, addr string,
) (net.Conn, error) {
	d := new(net.Dialer)
	return func(ctx context.Context, net, addr string) (net.Conn, error) {
		return d.DialContext(ctx, "unix", sinkAddr)
	}
}

func unixSockTransport(sockAddr string) *http.Transport {
	return &http.Transport{
		DialContext: unixSockSink(sockAddr),
	}
}
