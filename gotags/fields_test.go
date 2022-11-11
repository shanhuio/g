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

package gotags

import (
	"testing"
)

func TestParseFieldsEmpty(t *testing.T) {
	_, err := parseFields("")
	if err != nil {
		t.Fatalf("unexpected error from parseFields: %s", err)
	}
}

func TestParseFieldsLanguage(t *testing.T) {
	set, err := parseFields("+l")
	if err != nil {
		t.Fatalf("unexpected error from parseFields: %s", err)
	}
	if !set.Includes(Language) {
		t.Fatal("expected set to include Language")
	}
}

func TestParseFieldsInvalid(t *testing.T) {
	_, err := parseFields("junk")
	if err == nil {
		t.Fatal("expected parseFields to return error")
	}
	if _, ok := err.(ErrInvalidFields); !ok {
		t.Fatalf("expected parseFields to return error of type ErrInvalidFields, got %T", err)
	}
}
