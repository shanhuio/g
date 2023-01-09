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

package dock

import (
	"archive/tar"
	"io"
	"os"
	"path/filepath"

	"shanhu.io/pub/errcode"
)

func createFile(r io.Reader, name string, mod os.FileMode) error {
	const fileCreateFlag = os.O_RDWR | os.O_CREATE | os.O_TRUNC

	f, err := os.OpenFile(name, fileCreateFlag, mod)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := io.Copy(f, r); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}

func writeTarToDir(r io.Reader, destDir string) error {
	tr := tar.NewReader(r)
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		dest := filepath.Join(destDir, filepath.FromSlash(header.Name))

		switch typ := header.Typeflag; typ {
		case tar.TypeReg:
			dir := filepath.Dir(dest)
			if dir != "" && dir != "." {
				if err := os.MkdirAll(dir, 0700); err != nil {
					return err
				}
			}

			mod := header.FileInfo().Mode()
			if err := createFile(tr, dest, mod); err != nil {
				return err
			}
		case tar.TypeDir:
			if err := os.MkdirAll(dest, header.FileInfo().Mode()); err != nil {
				return err
			}
		default:
			return errcode.Internalf("type %s not supported", string(typ))
		}
	}
	return nil
}

func writeFirstFileAs(r io.Reader, file string) error {
	tr := tar.NewReader(r)
	written := false
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if written {
			continue
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}

		mod := header.FileInfo().Mode()
		if err := createFile(tr, file, mod); err != nil {
			return err
		}
		written = true
	}

	if !written {
		return errcode.NotFoundf("no file found")
	}
	return nil
}
