package aries

import (
	"net/url"
)

// DiscardURLServerParts discards the server parts of an URL,
// including scheme, host, port and user info.
func DiscardURLServerParts(u *url.URL) *url.URL {
	cp := *u
	cp.Scheme = ""
	cp.Opaque = ""
	cp.User = nil
	cp.Host = ""
	return &cp
}
