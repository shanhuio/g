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
