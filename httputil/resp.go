package httputil

import (
	"io"
	"net/http"
)

func copyRespBody(resp *http.Response, w io.Writer) error {
	defer resp.Body.Close()
	if w == nil {
		return nil
	}
	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return resp.Body.Close()
}
