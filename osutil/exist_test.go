package osutil

import (
	"testing"

	"os"
	"path"
)

func TestExist(t *testing.T) {
	d := t.TempDir()

	ok, err := Exist(d)
	if err != nil {
		t.Fatal("check exists: ", err)
	}
	if !ok {
		t.Errorf("dir %q should exist", d)
	}

	f := path.Join(d, "post")
	if err := os.WriteFile(f, []byte("post"), 0600); err != nil {
		t.Fatal("write file: ", err)
	}

	ok, err = Exist(f)
	if err != nil {
		t.Fatal("check file exists: ", err)
	}
	if !ok {
		t.Errorf("file %q should exist", f)
	}

	ghost := path.Join(d, "ghost")
	ok, err = Exist(ghost)
	if err != nil {
		t.Fatal("check ghost file exists: ", err)
	}
	if ok {
		t.Errorf("file %q should not exist", f)
	}
}
