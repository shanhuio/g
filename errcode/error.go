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

package errcode

import (
	"errors"
	"fmt"
)

// Error is a generic error with a string error code.
type Error struct {
	Code    string // code is the type of the error.
	Err     error  // error is the error message, human friendly.
	Message string // alternative message.
}

// Unwrap returns the unwrapped error.
func (e *Error) Unwrap() error { return e.Err }

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Err.Error()
}

// Common general error codes
const (
	NotFound     = "not-found"
	InvalidArg   = "invalid-arg"
	Internal     = "internal"
	Unauthorized = "unauthorized"
	TimeOut      = "time-out"
)

// Add creates a new error with code as the error code.
func Add(code string, err error) *Error {
	return &Error{
		Code: code,
		Err:  err,
	}
}

// Of returns the code of the error. For errors that
// do not have a code, it returns empty string.
func Of(err error) string {
	codeErr := new(Error)
	if errors.As(err, &codeErr) {
		return codeErr.Code
	}
	return ""
}

// IsNotFound checks if it is a not-found error.
func IsNotFound(err error) bool {
	return Of(err) == NotFound
}

// IsInvalidArg checks if it is an invalid argument error.
func IsInvalidArg(err error) bool {
	return Of(err) == InvalidArg
}

// IsInternal checks if it is an internal error.
func IsInternal(err error) bool {
	return Of(err) == Internal
}

// IsUnauthorized checks if it is an unauthorized error.
func IsUnauthorized(err error) bool {
	return Of(err) == Unauthorized
}

// IsTimeOut checks if it is a time-out error.
func IsTimeOut(err error) bool {
	return Of(err) == TimeOut
}

// Errorf creates an Error with the given error code.
func Errorf(code string, f string, args ...interface{}) *Error {
	return Add(code, fmt.Errorf(f, args...))
}

// AltErrorf replaces the message of err to be the formatted message, but keeps
// the error code.
func AltErrorf(err error, f string, args ...interface{}) error {
	msg := fmt.Sprintf(f, args...)
	return &Error{
		Err:     err,
		Message: msg,
		Code:    Of(err),
	}
}

// Annotate annotates an error but keeps the error code.
func Annotate(err error, msg string) error {
	return fmt.Errorf("%s: %w", msg, err)
}

// Annotatef annotates an error with a formatted message but keeps the error
// code.
func Annotatef(err error, f string, args ...interface{}) error {
	return Annotate(err, fmt.Sprintf(f, args...))
}

// NotFoundf creates a new not-found error.
func NotFoundf(f string, args ...interface{}) *Error {
	return Errorf(NotFound, f, args...)
}

// InvalidArgf creates a new invalid arugment error.
func InvalidArgf(f string, args ...interface{}) *Error {
	return Errorf(InvalidArg, f, args...)
}

// Internalf creates a new internal error.
func Internalf(f string, args ...interface{}) *Error {
	return Errorf(Internal, f, args...)
}

// Unauthorizedf returns an error caused by an unauthrozied request.
func Unauthorizedf(f string, args ...interface{}) *Error {
	return Errorf(Unauthorized, f, args...)
}

// TimeOutf returns a new time-out error.
func TimeOutf(f string, args ...interface{}) *Error {
	return Errorf(TimeOut, f, args...)
}
