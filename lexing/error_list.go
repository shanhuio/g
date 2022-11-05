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
	"fmt"
	"io"
)

// ErrorList saves a list of error
type ErrorList struct {
	errs []*Error

	Max int

	inJail bool
}

var _ Logger = new(ErrorList)

// NewErrorList creates a new error list with default (20) maximum
// lines of errors.
func NewErrorList() *ErrorList {
	ret := new(ErrorList)
	ret.Max = 20

	return ret
}

// Add appends the error to the list. Change the state to "in jail".
func (lst *ErrorList) Add(e *Error) {
	if e == nil {
		panic("nil error")
	}

	lst.inJail = true
	if len(lst.errs) >= lst.Max {
		return
	}

	lst.errs = append(lst.errs, e)
}

// AddAll adds a list of errors into the list.
func (lst *ErrorList) AddAll(es []*Error) {
	for _, e := range es {
		lst.Add(e)
	}
}

// Jail puts it in jail without generating a new error message
func (lst *ErrorList) Jail() { lst.inJail = true }

// InJail checks if a new error has been added since created or last bail out
func (lst *ErrorList) InJail() bool { return lst.inJail }

// BailOut clears the "in jail" state.
func (lst *ErrorList) BailOut() { lst.inJail = false }

// Errorf appends a new error with particular position and format.
func (lst *ErrorList) Errorf(p *Pos, f string, args ...interface{}) {
	lst.Add(&Error{p, fmt.Errorf(f, args...), ""})
}

// CodeErrorf appends a new error with a ErrorCode.
func (lst *ErrorList) CodeErrorf(p *Pos, c, f string, args ...interface{}) {
	lst.Add(&Error{p, fmt.Errorf(f, args...), c})
}

// Print prints to the writer (maximume lst.MaxPrint errors).
func (lst *ErrorList) Print(w io.Writer) error {
	for _, e := range lst.errs {
		_, pe := fmt.Fprintln(w, e)
		if pe != nil {
			return pe
		}
	}

	return nil
}

// Errs retunrs the errors in the list
func (lst *ErrorList) Errs() []*Error {
	ret := lst.errs
	if len(ret) == 0 {
		return nil
	}
	return ret
}

// SingleErr returns an error array with one error.
func SingleErr(err error) []*Error {
	return []*Error{{Err: err}}
}

// SingleCodeErr returns an error array with one error with ErrorCode.
func SingleCodeErr(code string, err error) []*Error {
	return []*Error{{Err: err, Code: code}}
}
