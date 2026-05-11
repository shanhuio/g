package tempfile

import (
	"io"
	"os"
)

// File is a temp file.
type File struct {
	*os.File
	Name        string
	SkipCleanUp bool
}

// NewFile creates a new temp file.
func NewFile(dir, prefix string) (*File, error) {
	f, err := os.CreateTemp(dir, prefix)
	if err != nil {
		return nil, err
	}
	return &File{
		File: f,
		Name: f.Name(),
	}, nil
}

// Reset is an alias for Seek(0, 0) on the file.
func (f *File) Reset() error {
	_, err := f.File.Seek(0, io.SeekStart)
	return err
}

// Remove removes the temp file.
func (f *File) Remove() error {
	return os.Remove(f.Name)
}

// Rename the file to another name.
func (f *File) Rename(p string) error {
	return os.Rename(f.Name, p)
}

// CleanUp closes the underlying temp file and removes it.
func (f *File) CleanUp() error {
	if f.SkipCleanUp {
		return nil
	}
	if err := f.File.Close(); err != nil {
		f.Remove() // remove anyways
		return err
	}
	return f.Remove()
}
