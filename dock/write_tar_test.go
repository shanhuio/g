package dock

import (
	"archive/tar"
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"shanhu.io/g/errcode"
)

type tarEntry struct {
	name    string
	mode    int64
	typ     byte
	content string
}

func buildTar(t *testing.T, entries []tarEntry) []byte {
	t.Helper()
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	for _, e := range entries {
		h := &tar.Header{
			Name:     e.name,
			Mode:     e.mode,
			Typeflag: e.typ,
			Size:     int64(len(e.content)),
		}
		if err := tw.WriteHeader(h); err != nil {
			t.Fatalf("WriteHeader %q: %v", e.name, err)
		}
		if e.content != "" {
			if _, err := tw.Write([]byte(e.content)); err != nil {
				t.Fatalf("Write %q: %v", e.name, err)
			}
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("tar Close: %v", err)
	}
	return buf.Bytes()
}

func TestWriteTarToDir(t *testing.T) {
	tarBytes := buildTar(t, []tarEntry{
		{name: "top.txt", mode: 0644, typ: tar.TypeReg, content: "top"},
		{name: "sub/", mode: 0755, typ: tar.TypeDir},
		{name: "sub/inner.txt", mode: 0644, typ: tar.TypeReg, content: "inner"},
		{
			name:    "auto/nested/deep.txt",
			mode:    0644,
			typ:     tar.TypeReg,
			content: "deep",
		},
	})

	dest := t.TempDir()
	if err := writeTarToDir(bytes.NewReader(tarBytes), dest); err != nil {
		t.Fatalf("writeTarToDir: %v", err)
	}

	for _, c := range []struct{ path, want string }{
		{path: "top.txt", want: "top"},
		{path: "sub/inner.txt", want: "inner"},
		{path: "auto/nested/deep.txt", want: "deep"},
	} {
		got, err := os.ReadFile(filepath.Join(dest, c.path))
		if err != nil {
			t.Errorf("ReadFile %q: %v", c.path, err)
			continue
		}
		if string(got) != c.want {
			t.Errorf(
				"%q content, got %q, want %q", c.path, got, c.want,
			)
		}
	}

	info, err := os.Stat(filepath.Join(dest, "sub"))
	if err != nil {
		t.Fatalf("stat sub: %v", err)
	}
	if !info.IsDir() {
		t.Errorf("sub: not a directory")
	}
}

func TestWriteTarToDirEmpty(t *testing.T) {
	tarBytes := buildTar(t, nil)
	dest := t.TempDir()
	if err := writeTarToDir(bytes.NewReader(tarBytes), dest); err != nil {
		t.Errorf("writeTarToDir on empty tar: %v", err)
	}
}

func TestWriteTarToDirUnsupportedType(t *testing.T) {
	buf := new(bytes.Buffer)
	tw := tar.NewWriter(buf)
	if err := tw.WriteHeader(&tar.Header{
		Name:     "link",
		Typeflag: tar.TypeSymlink,
		Linkname: "target",
	}); err != nil {
		t.Fatalf("WriteHeader: %v", err)
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("tar Close: %v", err)
	}

	err := writeTarToDir(buf, t.TempDir())
	if err == nil {
		t.Fatal("writeTarToDir, got nil, want error")
	}
	if !errcode.IsInternal(err) {
		t.Errorf("got %v, want internal error", err)
	}
}

func TestWriteFirstFileAs(t *testing.T) {
	tarBytes := buildTar(t, []tarEntry{
		{name: "first.txt", mode: 0644, typ: tar.TypeReg, content: "first"},
		{name: "second.txt", mode: 0644, typ: tar.TypeReg, content: "second"},
	})

	dest := filepath.Join(t.TempDir(), "out.txt")
	if err := writeFirstFileAs(bytes.NewReader(tarBytes), dest); err != nil {
		t.Fatalf("writeFirstFileAs: %v", err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "first" {
		t.Errorf("got %q, want %q", got, "first")
	}
}

func TestWriteFirstFileAsSkipsDirs(t *testing.T) {
	tarBytes := buildTar(t, []tarEntry{
		{name: "dir/", mode: 0755, typ: tar.TypeDir},
		{name: "dir/file.txt", mode: 0644, typ: tar.TypeReg, content: "found"},
	})

	dest := filepath.Join(t.TempDir(), "out.txt")
	if err := writeFirstFileAs(bytes.NewReader(tarBytes), dest); err != nil {
		t.Fatalf("writeFirstFileAs: %v", err)
	}
	got, err := os.ReadFile(dest)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(got) != "found" {
		t.Errorf("got %q, want %q", got, "found")
	}
}

func TestWriteFirstFileAsNotFound(t *testing.T) {
	for _, test := range []struct {
		name    string
		entries []tarEntry
	}{
		{name: "empty", entries: nil},
		{name: "dirs only", entries: []tarEntry{
			{name: "a/", mode: 0755, typ: tar.TypeDir},
			{name: "b/", mode: 0755, typ: tar.TypeDir},
		}},
	} {
		t.Run(test.name, func(t *testing.T) {
			tarBytes := buildTar(t, test.entries)
			dest := filepath.Join(t.TempDir(), "out.txt")
			err := writeFirstFileAs(bytes.NewReader(tarBytes), dest)
			if !errcode.IsNotFound(err) {
				t.Errorf("got %v, want NotFound", err)
			}
		})
	}
}
