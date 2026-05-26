package errcode

import (
	"errors"
	"fmt"

	stderrcode "shanhu.io/std/errcode"
)

// Error is a generic error with a string error code.
type Error = stderrcode.Error

// Common general error codes
const (
	NotFound     = stderrcode.NotFound
	InvalidArg   = stderrcode.InvalidArg
	Internal     = stderrcode.Internal
	Unauthorized = stderrcode.Unauthorized
	TimeOut      = stderrcode.TimeOut
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
func Errorf(code string, f string, args ...any) *Error {
	return Add(code, fmt.Errorf(f, args...))
}

// AltErrorf replaces the message of err to be the formatted message, but keeps
// the error code.
func AltErrorf(err error, f string, args ...any) error {
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
func Annotatef(err error, f string, args ...any) error {
	return Annotate(err, fmt.Sprintf(f, args...))
}

// NotFoundf creates a new not-found error.
func NotFoundf(f string, args ...any) *Error {
	return Errorf(NotFound, f, args...)
}

// InvalidArgf creates a new invalid arugment error.
func InvalidArgf(f string, args ...any) *Error {
	return Errorf(InvalidArg, f, args...)
}

// Internalf creates a new internal error.
func Internalf(f string, args ...any) *Error {
	return Errorf(Internal, f, args...)
}

// Unauthorizedf returns an error caused by an unauthrozied request.
func Unauthorizedf(f string, args ...any) *Error {
	return Errorf(Unauthorized, f, args...)
}

// TimeOutf returns a new time-out error.
func TimeOutf(f string, args ...any) *Error {
	return Errorf(TimeOut, f, args...)
}
