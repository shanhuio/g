// Copyright (C) 2023  Shanhu Tech Inc.
//
// This program is free software: you can redistribute it and/or modify it
// under the terms of the GNU Affero General Public License as published by the
// Free Software Foundation, either version 3 of the License, or (at your
// option) any later version.
//
// This program is distributed in the hope that it will be useful, but WITHOUT
// ANY WARRANTY; without even the implied warranty of MERCHANTABILITY or
// FITNESS FOR A PARTICULAR PURPOSE.  See the GNU Affero General Public License
// for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package ziputil

import (
	"testing"

	"archive/zip"
	"io"
	"os"
	"path"
	"reflect"

	"shanhu.io/g/osutil"
	"shanhu.io/g/tempfile"
)

func testDiffFile(t *testing.T, f1, f2 string) bool {
	t.Helper()

	bs1, err := os.ReadFile(f1)
	if err != nil {
		t.Fatal("read f1: ", err)
	}

	bs2, err := os.ReadFile(f2)
	if err != nil {
		t.Fatal("read f2: ", err)
	}

	if !reflect.DeepEqual(bs1, bs2) {
		return false
	}

	s1, err := os.Stat(f1)
	if err != nil {
		t.Fatal("stat f1: ", err)
	}

	s2, err := os.Stat(f2)
	if err != nil {
		t.Fatal("stat f2: ", err)
	}

	if s1.Mode() != s2.Mode() {
		return false
	}
	return true
}

func TestZipFile(t *testing.T) {
	temp, err := tempfile.NewFile("", "ziputil")
	if err != nil {
		t.Fatal("new temp file: ", err)
	}
	defer temp.CleanUp()

	const p = "testdata/testfile"
	if err := ZipFile(p, temp); err != nil {
		t.Fatal("zip file: ", err)
	}

	size, err := temp.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Fatal("seek file: ", err)
	}

	if err := temp.Reset(); err != nil {
		t.Fatal("reset: ", err)
	}

	output := t.TempDir()

	z, err := zip.NewReader(temp, size)
	if err != nil {
		t.Fatal("new reader: ", err)
	}

	if err := UnzipDir(output, z, true); err != nil {
		t.Fatal("unzip: ", err)
	}

	outPath := path.Join(output, "testfile")
	if !testDiffFile(t, outPath, p) {
		t.Error("zip loop back failed")
	}
}

func TestZipDir(t *testing.T) {
	temp, err := tempfile.NewFile("", "ziputil")
	if err != nil {
		t.Error("new temp file: ", err)
	}
	defer temp.CleanUp()

	const p = "testdata/testdir"
	if err := ZipDir(p, temp); err != nil {
		t.Error("zip dir: ", err)
	}

	size, err := temp.Seek(0, io.SeekCurrent)
	if err != nil {
		t.Error("seek: ", err)
	}

	if err := temp.Reset(); err != nil {
		t.Error("reset: ", err)
	}

	output := t.TempDir()

	z, err := zip.NewReader(temp, size)
	if err != nil {
		t.Error("new reader: ", err)
	}

	if err := UnzipDir(output, z, true); err != nil {
		t.Error("unzip: ", err)
	}

	for _, name := range []string{
		"bin-file", "private-file", "text-file",
	} {
		outPath := path.Join(output, name)
		target := path.Join(p, name)
		if !testDiffFile(t, outPath, target) {
			t.Errorf("zip loop back failed for file %q", name)
		}
	}
}

func testClearDir(t *testing.T, clear bool) {
	z, err := zip.OpenReader("testdata/testfile.zip")
	if err != nil {
		t.Fatal("open: ", err)
	}
	defer z.Close()

	output := t.TempDir()

	ob := path.Join(output, "native-file")
	msg := []byte("lived here long time ago")
	if err := os.WriteFile(ob, msg, 0600); err != nil {
		t.Fatal("write file: ", err)
	}

	if err := UnzipDir(output, &z.Reader, clear); err != nil {
		t.Fatal("unzip dir: ", err)
	}

	exist, err := osutil.Exist(ob)
	if err != nil {
		t.Fatal("exist: ", err)
	}

	if clear && exist {
		t.Error("should clear directory, but still see the file")
	}
	if !clear && !exist {
		t.Error("should preserve the file, but lost")
	}
}

func TestClearDir(t *testing.T) { testClearDir(t, true) }

func TestNoClearDir(t *testing.T) { testClearDir(t, false) }
