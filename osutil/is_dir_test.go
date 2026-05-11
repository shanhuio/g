package osutil

import (
	"testing"

	"os"
	"path"
)

func TestIsDir(t *testing.T) {
	d := t.TempDir()

	ok, err := IsDir(d)
	if err != nil {
		t.Fatal("check is dir: ", err)
	}
	if !ok {
		t.Errorf("IsDir(%q) should return true", d)
	}

	f := path.Join(d, "post")
	if err := os.WriteFile(f, []byte("post"), 0600); err != nil {
		t.Fatal("write file: ", err)
	}

	ok, err = IsDir(f)
	if err != nil {
		t.Fatal("check file is dir: ", err)
	}
	if ok {
		t.Errorf("IsDir(%q) should return false", f)
	}

	ok, err = IsDir(path.Join(d, "ghost"))
	if err != nil {
		t.Fatal("check ghost dir: ", err)
	}
	if ok {
		t.Errorf("IsDir(%q) should return false", f)
	}
}
