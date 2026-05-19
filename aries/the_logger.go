package aries

import (
	"fmt"
)

// TheLogger is the default logger that logs to default golang log.
var TheLogger = StdLogger()

// AltInternal prints error to TheLogger. It is an alias to
// TheLogger.AltInteral
func AltInternal(err error, s string) error {
	return TheLogger.AltInternal(err, s)
}

// AltInternalf printe the formatted error to TheLogger.
func AltInternalf(err error, f string, args ...any) error {
	return TheLogger.AltInternal(err, fmt.Sprintf(f, args...))
}
