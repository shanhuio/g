package oauth2

import (
	"net/url"
	"strings"

	"shanhu.io/g/aries"
	"shanhu.io/std/errcode"
)

// ParseRedirect parses an in-site redirection URL.
// The server parts (scheme, host, port, user info) are discarded.
func ParseRedirect(redirect string) (string, error) {
	if redirect == "" {
		return "", nil
	}

	u, err := url.Parse(redirect)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(u.Path, "/") {
		return "", errcode.InvalidArgf(
			"redirect path part %q is not absolute", u.Path,
		)
	}

	return aries.DiscardURLServerParts(u).String(), nil
}
