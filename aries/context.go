// Copyright (C) 2023  Shanhu Tech Inc.
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

package aries

import (
	"context"
	"net/http"
	"strings"
	"time"

	"shanhu.io/pub/errcode"
)

// C provides the request context for a web application.
type C struct {
	Path string

	User      string
	UserLevel int // 0 for normal user. 0 with empty User is anonymous.
	UserData  interface{}

	Req     *http.Request
	Resp    http.ResponseWriter
	Context context.Context

	HTTPS bool

	Data map[string]interface{}

	route    *route
	routePos int
}

// NewContext creates a new context from the incomming request.
func NewContext(w http.ResponseWriter, req *http.Request) *C {
	isHTTPS := false
	u := req.URL
	if req.TLS != nil {
		isHTTPS = true
	} else if strings.ToLower(req.Header.Get("X-Forwarded-Proto")) == "https" {
		isHTTPS = true
	}

	return &C{
		Path:    u.Path,
		Resp:    w,
		Req:     req,
		Context: req.Context(),
		HTTPS:   isHTTPS,
		Data:    make(map[string]interface{}),

		route: newRoute(u.Path),
	}
}

// Redirect redirects the request to another URL.
func (c *C) Redirect(url string) {
	http.Redirect(c.Resp, c.Req, url, http.StatusFound)
}

// Rel returns the current relative route. The return value changes if the
// routing is using a router, otherwise, it will always return the full routing
// path.
func (c *C) Rel() string { return c.route.rel(c.routePos) }

// RelRoute returns the current relative route string array.
func (c *C) RelRoute() []string {
	ret := c.route.relRoute(c.routePos)
	cp := make([]string, len(ret))
	copy(cp, ret)
	return cp
}

// ShiftRoute shift the routing pointer by inc.
func (c *C) ShiftRoute(inc int) {
	c.routePos += inc
	if c.routePos >= c.route.size() {
		c.routePos = c.route.size()
	}
}

// PathIsDir return true if the path ends with a slash.
func (c *C) PathIsDir() bool { return c.route.isDir }

// Current returns the next part in the current relative route.
func (c *C) Current() string { return c.route.current(c.routePos) }

// ReadCookie reads the cookie from the context.
func (c *C) ReadCookie(name string) string {
	cookie, err := c.Req.Cookie(name)
	if err != nil || cookie == nil {
		return ""
	}
	return cookie.Value
}

// WriteCookie sets a cookie.
func (c *C) WriteCookie(name, v string, expires time.Time) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    v,
		Path:     "/",
		Expires:  expires,
		Secure:   c.HTTPS,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Resp, cookie)
}

// ClearCookie clears a cookie.
func (c *C) ClearCookie(name string) {
	cookie := &http.Cookie{
		Name:     name,
		Value:    "",
		Path:     "/",
		Secure:   c.HTTPS,
		MaxAge:   5, // Delete in 5 seconds.
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Resp, cookie)
}

// ErrCode returns an error based on its error code.
func (c *C) ErrCode(err error) bool {
	if err == nil {
		return false
	}
	code := errcode.Of(err)
	switch code {
	case errcode.NotFound:
		return c.replyError(404, err)
	case errcode.Internal:
		return c.replyError(500, err)
	case errcode.Unauthorized:
		return c.replyError(403, err)
	case errcode.InvalidArg:
		return c.replyError(400, err)
	}
	return c.replyError(500, err)
}

func (c *C) replyError(code int, err error) bool {
	if err == nil {
		return false
	}
	http.Error(c.Resp, err.Error(), code)
	return true
}

// IsMobile checks if the user agent of the request is mobile or not.
func (c *C) IsMobile() bool {
	return isMobile(c.Req.UserAgent())
}

// Method returns the c.Req.Method.
func (c *C) Method() string { return c.Req.Method }
