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
	"testing"
)

func TestIsLetter(t *testing.T) {
	for _, r := range "abzdATZ" {
		if !IsLetter(r) {
			t.Errorf("%v should be a letter", r)
		}
	}

	for _, r := range "013_%~-" {
		if IsLetter(r) {
			t.Errorf("%v should not be a letter", r)
		}
	}
}

func TestIsDigit(t *testing.T) {
	for _, r := range "0123456789" {
		if !IsDigit(r) {
			t.Errorf("%v should be a digit", r)
		}
	}

	for _, r := range "abzATZ#%~" {
		if IsDigit(r) {
			t.Errorf("%v should not be a digit", r)
		}
	}
}

func TestIsHexDigit(t *testing.T) {
	for _, r := range "0123456789" {
		if !IsHexDigit(r) {
			t.Errorf("%v should be a hexdigit", r)
		}
	}

	for _, r := range "abcdefABCDEF" {
		if !IsHexDigit(r) {
			t.Errorf("%v should be a hexdigit", r)
		}
	}

	for _, r := range "gJmXY!@*" {
		if IsHexDigit(r) {
			t.Errorf("%v should not be a hexdigit", r)
		}
	}
}
