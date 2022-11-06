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
	"fmt"
	"regexp"
)

// FieldSet is a set of extension fields to include in a tag.
type FieldSet map[TagField]bool

// Includes tests whether the given field is included in the set.
func (f FieldSet) Includes(field TagField) bool {
	b, ok := f[field]
	return ok && b
}

// ErrInvalidFields is an error returned when attempting to parse invalid
// fields.
type ErrInvalidFields struct {
	Fields string
}

func (e ErrInvalidFields) Error() string {
	return fmt.Sprintf("invalid fields: %s", e.Fields)
}

var (
	// currently only "+l" is supported
	fieldsPattern  = regexp.MustCompile(`^\+l$`)
	symbolsPattern = regexp.MustCompile(`^\+q$`)
)

func parseFields(fields string) (FieldSet, error) {
	if fields == "" {
		return FieldSet{}, nil
	}
	if fieldsPattern.MatchString(fields) {
		return FieldSet{Language: true}, nil
	}
	return FieldSet{}, ErrInvalidFields{fields}
}

func parseExtraSymbols(symbols string) (FieldSet, error) {
	if symbols == "" {
		return FieldSet{}, nil
	}
	if symbolsPattern.MatchString(symbols) {
		return FieldSet{ExtraTags: true}, nil
	}
	return FieldSet{}, ErrInvalidFields{symbols}
}
