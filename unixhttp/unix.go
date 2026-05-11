package unixhttp

import (
	"net"
	"net/http"
	"os"

	"shanhu.io/g/osutil"
)

// Listen creates a listener at the given unix domain socket address.
func Listen(p string) (net.Listener, error) {
	exist, err := osutil.IsSock(p)
	if err != nil {
		return nil, err
	}

	if exist {
		if err := os.Remove(p); err != nil {
			return nil, err
		}
	}
	lis, err := net.ListenUnix("unix", &net.UnixAddr{
		Name: p,
		Net:  "unix",
	})
	if err != nil {
		return nil, err
	}
	return lis, nil
}

// ListenAndServe listens and serves at the given unix domain socket
// path.
func ListenAndServe(p string, h http.Handler) error {
	lis, err := Listen(p)
	if err != nil {
		return err
	}
	return http.Serve(lis, h)
}
