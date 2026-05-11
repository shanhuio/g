package dock

import (
	"net/url"
	"path"
)

func urlQuery(p string, q url.Values) string {
	u := &url.URL{Path: p}
	if len(q) != 0 {
		u.RawQuery = q.Encode()
	}
	return u.String()
}

func apiURLQuery(p string, q url.Values) string {
	return urlQuery(path.Join(apiVersion, p), q)
}

func singleQuery(k, v string) url.Values {
	q := make(url.Values)
	q.Add(k, v)
	return q
}
