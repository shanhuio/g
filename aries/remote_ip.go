package aries

import (
	"net"
	"strings"
)

// RemoteIP returns the remote IP address.
func RemoteIP(c *C) net.IP {
	forwardedFor := c.Req.Header.Get("X-Forwarded-For")
	ips := strings.SplitSeq(forwardedFor, ",")
	for ip := range ips {
		if parsed := net.ParseIP(ip); parsed != nil {
			return parsed
		}
	}

	remoteAddr := c.Req.RemoteAddr
	if remoteAddr == "@" || remoteAddr == "" {
		return nil
	}
	host, _, err := net.SplitHostPort(remoteAddr)
	if err != nil {
		return nil
	}
	return net.ParseIP(host)
}

// RemoteIPString returns the string form of the remote IP address.
// It returns empty string when IP cannot be determined.
func RemoteIPString(c *C) string {
	ip := RemoteIP(c)
	if ip == nil {
		return ""
	}
	return ip.String()
}
