// Copyright (C) 2022  Shanhu Tech Inc.
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

package dock

import (
	"archive/tar"
	"bytes"
	"io"
	"io/ioutil"

	"shanhu.io/pub/errcode"
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

		bs, err := ioutil.ReadAll(r)
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
