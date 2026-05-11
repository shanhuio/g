package dock

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"
)

func TestGzipRoundTrip(t *testing.T) {
	for _, test := range []struct {
		name    string
		content []byte
	}{
		{name: "small", content: []byte("hello, gzip")},
		{name: "empty", content: nil},
		{name: "binary", content: bytes.Repeat([]byte{0, 1, 2, 3}, 1024)},
	} {
		t.Run(test.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "out.gz")

			w, err := gzipCreate(path)
			if err != nil {
				t.Fatalf("gzipCreate: %v", err)
			}
			if _, err := w.Write(test.content); err != nil {
				w.Close()
				t.Fatalf("Write: %v", err)
			}
			if err := w.Close(); err != nil {
				t.Fatalf("Close: %v", err)
			}

			r, err := gzipOpen(path)
			if err != nil {
				t.Fatalf("gzipOpen: %v", err)
			}
			got, err := io.ReadAll(r)
			if err != nil {
				r.Close()
				t.Fatalf("ReadAll: %v", err)
			}
			if err := r.Close(); err != nil {
				t.Fatalf("Close: %v", err)
			}

			if !bytes.Equal(got, test.content) {
				t.Errorf(
					"round-trip mismatch: got %d bytes, want %d",
					len(got), len(test.content),
				)
			}
		})
	}
}

func TestGzipOpenMissing(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nonexistent.gz")
	r, err := gzipOpen(path)
	if err == nil {
		r.Close()
		t.Fatal("gzipOpen, got nil, want error")
	}
}

func TestGzipOpenNotGzip(t *testing.T) {
	path := filepath.Join(t.TempDir(), "plain.txt")
	if err := os.WriteFile(path, []byte("not gzipped"), 0644); err != nil {
		t.Fatalf("WriteFile: %v", err)
	}
	r, err := gzipOpen(path)
	if err == nil {
		r.Close()
		t.Fatal("gzipOpen, got nil, want error")
	}
}

func TestGzipCreateBadPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "no_such_dir", "out.gz")
	w, err := gzipCreate(path)
	if err == nil {
		w.Close()
		t.Fatal("gzipCreate, got nil, want error")
	}
}
