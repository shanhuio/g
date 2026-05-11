package ziputil

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ZipFile creates a zip file of a single file.
func ZipFile(file string, w io.Writer) error {
	abs, err := filepath.Abs(file)
	if err != nil {
		return err
	}
	name := filepath.Base(abs)
	if name == "" {
		return fmt.Errorf("missing name for for the file")
	}

	info, err := os.Stat(file)
	if err != nil {
		return err
	}

	fin, err := os.Open(file)
	if err != nil {
		return err
	}
	defer fin.Close()

	ar := zip.NewWriter(w)
	h := &zip.FileHeader{Name: name}
	h.SetMode(info.Mode())
	h.SetModTime(info.ModTime())

	zipFile, err := ar.CreateHeader(h)
	if err != nil {
		return err
	}

	if _, err := io.Copy(zipFile, fin); err != nil {
		return err
	}

	if err := fin.Close(); err != nil {
		return err
	}
	return ar.Close()
}

// ZipDir creates a zip file of a directory with all file in it.
func ZipDir(dir string, w io.Writer) error {
	ar := zip.NewWriter(w)
	walk := func(p string, info os.FileInfo, err error) error {
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}

		mod := info.Mode()
		t := info.ModTime()

		if info.IsDir() {
			h := &zip.FileHeader{Name: rel + "/"}
			h.SetMode(mod)
			h.SetModTime(t)
			_, err := ar.CreateHeader(h)
			return err
		}

		fin, err := os.Open(p)
		if err != nil {
			return err
		}
		defer fin.Close()

		h := &zip.FileHeader{Name: rel}
		h.SetMode(mod)
		h.SetModTime(t)

		w, err := ar.CreateHeader(h)
		if err != nil {
			return err
		}

		if _, err = io.Copy(w, fin); err != nil {
			return err
		}
		return fin.Close()
	}

	if err := filepath.Walk(dir, walk); err != nil {
		return err
	}

	return ar.Close()
}
