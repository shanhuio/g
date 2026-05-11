package ziputil

import (
	"archive/zip"
	"io"

	"shanhu.io/g/tempfile"
)

// OpenInTemp copies all bytes of r into a temp file, and then opens the
// creates a zip reader that is backed by f.
func OpenInTemp(r io.Reader, tmp *tempfile.File) (*zip.Reader, error) {
	n, err := io.Copy(tmp, r)
	if err != nil {
		return nil, err
	}
	if err := tmp.Reset(); err != nil {
		return nil, err
	}
	return zip.NewReader(tmp, n)
}
