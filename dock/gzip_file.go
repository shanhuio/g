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
	"compress/gzip"
	"os"
)

type gzipReader struct {
	f *os.File
	*gzip.Reader
}

func gzipOpen(p string) (*gzipReader, error) {
	f, err := os.Open(p)
	if err != nil {
		return nil, err
	}

	r, err := gzip.NewReader(f)
	if err != nil {
		f.Close()
		return nil, err
	}

	return &gzipReader{f: f, Reader: r}, nil
}

func (r *gzipReader) Close() error {
	if err := r.Reader.Close(); err != nil {
		r.f.Close()
		return err
	}
	return r.f.Close()
}

type gzipWriter struct {
	f *os.File
	*gzip.Writer
}

func gzipCreate(p string) (*gzipWriter, error) {
	f, err := os.Create(p)
	if err != nil {
		return nil, err
	}

	w := gzip.NewWriter(f)
	return &gzipWriter{f: f, Writer: w}, nil
}

func (w *gzipWriter) Close() error {
	if err := w.Writer.Close(); err != nil {
		w.f.Close()
		return err
	}
	return w.f.Close()
}
