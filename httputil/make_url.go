package httputil

import (
	"net/url"
	"path"
)

func makeURL(base *url.URL, p string) (string, error) {
	u := *base
	up, err := url.Parse(p)
	if err != nil {
		return "", err
	}

	// append two paths
	u.Path = path.Join(u.Path, up.Path)
	u.RawQuery = up.RawQuery
	u.Fragment = up.Fragment
	return u.String(), nil
}
