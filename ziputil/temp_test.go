package ziputil

import (
	"testing"

	"bytes"
	"os"

	"shanhu.io/g/tempfile"
)

func TestOpenInTemp(t *testing.T) {
	ne := func(err error) {
		if err != nil {
			t.Fatal(err)
		}
	}

	bs, err := os.ReadFile("testdata/testfile.zip")
	ne(err)

	f, err := tempfile.NewFile("", "ziputil")
	ne(err)
	defer f.CleanUp()

	r, err := OpenInTemp(bytes.NewReader(bs), f)
	ne(err)

	if len(r.File) != 1 {
		t.Fatal("want 1 file in testfile.zip")
	}

	got := r.File[0].Name
	want := "testfile"
	if got != want {
		t.Errorf("file name want %q, got %q", want, got)
	}
}
