package httputil

import (
	"net/http"
)

func setHeader(h http.Header, k, v string) {
	if v == "" {
		return
	}
	h.Set(k, v)
}
