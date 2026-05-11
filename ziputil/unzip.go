package ziputil

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// UnzipDir unzips a zip file into a directory.
// If the directory already exists, it removes all stuff in the directory
// if the clear flag is set to true.
func UnzipDir(dir string, r *zip.Reader, clear bool) error {
	if clear {
		err := os.RemoveAll(dir)
		if err != nil {
			return err
		}
	}

	for _, f := range r.File {
		mod := f.Mode()

		name := filepath.Join(dir, f.Name)
		if mod.IsDir() {
			if err := os.MkdirAll(name, mod); err != nil {
				return err
			}
			continue
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		if err := os.MkdirAll(filepath.Dir(name), 0700); err != nil {
			return err
		}

		fout, err := os.Create(name)
		if err != nil {
			return err
		}
		defer fout.Close()

		if err := fout.Chmod(mod); err != nil {
			return err
		}

		if _, err := io.Copy(fout, rc); err != nil {
			return err
		}

		if err := fout.Close(); err != nil {
			return err
		}
	}
	return nil
}
