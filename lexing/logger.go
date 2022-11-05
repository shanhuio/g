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

// Logger is an error logging interface
type Logger interface {
	Errorf(p *Pos, fmt string, args ...interface{})
}

// LogError adds a error to the logger if the error is not nil and returns
// true.  If the error is nil, it returns false.
func LogError(log Logger, e error) bool {
	if e == nil {
		return false
	}

	log.Errorf(nil, "%s", e.Error())
	return true
}
