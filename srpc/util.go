package srpc

import (
	"net/url"
	"path"
	"strings"
)

const contentTypeJSON = "application/json"

func urlJoin(server *url.URL, p string) *url.URL {
	u := *server
	u.Path = path.Join("/", server.Path, p)
	if strings.HasSuffix(p, "/") && !strings.HasSuffix(u.Path, "/") {
		u.Path += "/"
	}
	return &u
}
