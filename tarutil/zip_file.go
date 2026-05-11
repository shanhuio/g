package tarutil

import (
	"archive/tar"
	"archive/zip"
	"io"
	"path"

	"shanhu.io/g/errcode"
)

func copyZipFile(w io.Writer, f *zip.File) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	if _, err := io.Copy(w, rc); err != nil {
		rc.Close()
		return err
	}
	return rc.Close()
}

// TarZipFile puts all files from a zip file into a tar stream.
func TarZipFile(tw *tar.Writer, p string, dir string) error {
	z, err := zip.OpenReader(p)
	if err != nil {
		return errcode.Annotate(err, "open zip file")
	}

	for _, f := range z.File {
		stat := f.FileInfo()
		tarStat, err := tar.FileInfoHeader(stat, "")
		if err != nil {
			return errcode.Annotatef(err, "tar stat for: %q", f.Name)
		}
		name := f.Name
		if dir != "" {
			name = path.Join(dir, name)
		}
		tarStat.Name = name
		if err := tw.WriteHeader(tarStat); err != nil {
			return errcode.Annotatef(err, "write header: %q", f.Name)
		}
		if err := copyZipFile(tw, f); err != nil {
			return errcode.Annotatef(err, "copy zip file: %q", f.Name)
		}
	}

	return nil
}
