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

package gocheck

import (
	"os"

	"shanhu.io/pub/lexing"
	"shanhu.io/pub/textbox"
)

// CheckRect checks if all the files are within the given rectangle.
func CheckRect(files []string, h, w int) []*lexing.Error {
	errs := lexing.NewErrorList()
	for _, f := range files {
		fin, err := os.Open(f)
		if lexing.LogError(errs, err) {
			continue
		}

		errs.AddAll(textbox.CheckRect(f, fin, h, w))
		if lexing.LogError(errs, fin.Close()) {
			continue
		}
	}

	return errs.Errs()
}
