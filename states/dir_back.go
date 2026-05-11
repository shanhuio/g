package states

import (
	"net/url"
	"os"
	"path/filepath"

	"shanhu.io/g/errcode"
)

type dirBack struct {
	dir string
}

func newDirBack(dir string) *dirBack {
	return &dirBack{dir: dir}
}

func (b *dirBack) filepath(key string) string {
	return filepath.Join(b.dir, filepath.ToSlash(key))
}

func (b *dirBack) Get(_ C, key string) ([]byte, error) {
	bs, err := os.ReadFile(b.filepath(key))
	return bs, errcode.FromOS(err)
}

func writeFile(p string, data []byte) error {
	const flag = os.O_RDWR | os.O_CREATE | os.O_TRUNC
	f, err := os.OpenFile(p, flag, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(data); err != nil {
		return err
	}
	return f.Sync() // Flush to stable storage.
}

func (b *dirBack) Put(_ C, key string, data []byte) error {
	p := b.filepath(key)
	if err := os.MkdirAll(filepath.Dir(p), 0700); err != nil {
		return err
	}
	return errcode.FromOS(writeFile(p, data))
}

func (b *dirBack) Del(_ C, key string) error {
	return errcode.FromOS(os.Remove(b.filepath(key)))
}

func (b *dirBack) URL() *url.URL {
	return &url.URL{
		Scheme: "file",
		Path:   b.dir,
	}
}
