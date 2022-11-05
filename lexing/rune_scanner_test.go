// Copyright (C) 2022  Shanhu Tech Inc.
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

package lexing

import (
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestRuneScanner(t *testing.T) {
	testCase := []struct {
		r    rune
		line int
		col  int
	}{
		{'a', 1, 1},
		{'~', 1, 2},
		{' ', 1, 3},
		{'\n', 1, 4},
		{'\n', 2, 1},
		{'1', 3, 1},
		{'A', 3, 2},
	}

	r := strings.NewReader("a~ \n\n1A")
	file := "a.txt"
	s := newRuneScanner(file, r)
	for _, tc := range testCase {
		if !s.scan() {
			t.Fatal("scan failed")
		}
		p := s.pos()
		want := &Pos{
			Col:  tc.col,
			Line: tc.line,
			File: file,
		}
		if !reflect.DeepEqual(p, want) {
			t.Errorf("pos got %v, want %v", p, want)
		}
		if s.Rune != tc.r {
			t.Errorf("rune got %c, want %c", s.Rune, tc.r)
		}
	}
	if s.scan() {
		t.Error("s.scan() got false, want true")
	}
	if !s.closed {
		t.Error("s close got false, want true")
	}
	if s.Err != io.EOF {
		t.Errorf("unexpected error %v", s.Err)
	}
}

type errorReader struct {
	r io.Reader
	n int
}

var errTest = errors.New("timeout")

func (r *errorReader) Read(bs []byte) (int, error) {
	if r.n == 0 {
		r.n++
		return r.r.Read(bs)
	}
	return 0, errTest
}

func TestRuneScannerError(t *testing.T) {
	r := &errorReader{r: strings.NewReader("x")}
	s := newRuneScanner("a.txt", r)
	if !s.scan() {
		t.Error("first scan should succeed")
	}
	if s.scan() {
		t.Error("second scan should fail")
	}
	if s.Err != errTest {
		t.Errorf("got err %v, want errTest", s.Err)
	}
}
