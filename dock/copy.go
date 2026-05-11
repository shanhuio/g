package dock

import (
	"os"

	"shanhu.io/g/errcode"
	"shanhu.io/g/tarutil"
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
