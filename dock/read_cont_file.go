package dock

import (
	"archive/tar"
	"bytes"
	"io"

	"shanhu.io/g/errcode"
)

// ReadContFile reads out a single file from a container.
func ReadContFile(c *Cont, f string) ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := c.CopyOutTar(f, buf); err != nil {
		return nil, errcode.Annotate(err, "read from container")
	}

	var content []byte
	got := false
	r := tar.NewReader(buf)
	for {
		h, err := r.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, errcode.Annotate(err, "read tar")
		}

		if got {
			return nil, errcode.InvalidArgf("not one single file")
		}
		if h.Typeflag != tar.TypeReg {
			return nil, errcode.InvalidArgf("not a regular file")
		}

		bs, err := io.ReadAll(r)
		if err != nil {
			return nil, errcode.Annotate(err, "read file content")
		}
		content = bs
		got = true
	}

	if !got {
		return nil, errcode.NotFoundf("file %q not found", f)
	}

	return content, nil
}
