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
	"os"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/tarutil"
)

// CopyInTarGz copies a gzipped tarball into the container.
func CopyInTarGz(c *Cont, f, target string) error {
	r, err := gzipOpen(f)
	if err != nil {
		return err
	}
	defer r.Close()

	// copy in the source tarbal
	if err := c.CopyInTar(r, target); err != nil {
		return err
	}
	return r.Close()
}

// CopyOutTarGz copies out a directory of file into a gzipped
// tarball.
func CopyOutTarGz(c *Cont, target, f string) error {
	w, err := gzipCreate(f)
	if err != nil {
		return err
	}
	defer w.Close()

	if err := c.CopyOutTar(target, w); err != nil {
		return err
	}
	return w.Close()
}

// CopyInTarStream copies in a tarutil.Stream into the container.
func CopyInTarStream(c *Cont, files *tarutil.Stream, target string) error {
	r := newWriteToReader(files)
	defer r.Close()

	if err := c.CopyInTar(r, target); err != nil {
		return err
	}
	if err := r.Join(); err != nil {
		return errcode.Annotate(err, "copy files")
	}
	return nil
}

// CopyInTarFile copies a tar file into the container.
func CopyInTarFile(c *Cont, file, toPath string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	return c.CopyInTar(f, toPath)
}

// CopyOutTarFile copies a file or a directory out of the container
// into a tarball file.
func CopyOutTarFile(c *Cont, fromPath, file string) error {
	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	if err := c.CopyOutTar(fromPath, f); err != nil {
		return err
	}
	return f.Close()
}
