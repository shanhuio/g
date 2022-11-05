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

package jsonx

import (
	"testing"

	"reflect"
	"strings"
)

func TestDecoder(t *testing.T) {
	input := strings.NewReader(`"a""b";"c"`)

	dec := NewDecoder(input)
	var got []string
	for dec.More() {
		var s string
		if err := dec.Decode(&s); err != nil {
			t.Fatal(err)
		}
		got = append(got, s)
	}

	want := []string{"a", "b", "c"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestUnmarshal_error(t *testing.T) {
	var s string
	if err := Unmarshal([]byte(`"missing`), &s); err == nil {
		t.Errorf("parse incomplete string passed")
	}
}

func TestUnmarshal(t *testing.T) {
	var v int
	if err := Unmarshal([]byte("1234"), &v); err != nil {
		t.Fatal(err)
	}
	if v != 1234 {
		t.Errorf("got %d, want 1234", v)
	}
}

func TestDecoder_series(t *testing.T) {
	input := strings.NewReader(strings.Join([]string{
		`str "string"`,
		`num 3`,
		`struct {Field: "value"}`,
	}, "\n"))

	type structType struct {
		Field string
	}

	dec := NewDecoder(input)
	tm := func(t string) interface{} {
		switch t {
		case "str":
			return new(string)
		case "num":
			return new(int)
		case "struct":
			return new(structType)
		}
		return nil
	}
	list, errs := dec.DecodeSeries(tm)
	if errs != nil {
		for _, err := range errs {
			t.Error(err)
		}
	}

	strVal := "string"
	numVal := 3
	want := []*Typed{
		{Type: "str", V: &strVal},
		{Type: "num", V: &numVal},
		{Type: "struct", V: &structType{Field: "value"}},
	}
	if len(list) != len(want) {
		t.Errorf("got %d entries, want %d", len(list), len(want))
	} else {
		for i, got := range list {
			w := want[i]
			if got.Type != w.Type {
				t.Errorf(
					"entry #%d, got type %q, want %q",
					i, got.Type, w.Type,
				)
			}
			if !reflect.DeepEqual(got.V, w.V) {
				t.Errorf(
					"entry value #%d, got %+v, want %+v",
					i, got.V, w.V,
				)
			}
		}
	}
}

func TestDecoder_series_error(t *testing.T) {
	s := strings.Join([]string{
		`t {`,
		`	a:"x/**`,
		`}`,
	}, "\n")
	input := strings.NewReader(s)
	dec := NewDecoder(input)
	type structType struct{}
	tm := func(t string) interface{} {
		return new(structType)
	}
	if _, errs := dec.DecodeSeries(tm); errs == nil {
		t.Errorf("decode %q got no error", s)
	} else {
		for _, err := range errs {
			t.Log(err)
		}
	}
}
