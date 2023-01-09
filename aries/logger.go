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

package aries

import (
	"errors"
	"fmt"
	"log"
	"os"

	"shanhu.io/pub/errcode"
)

// LogPrinter is the interface for printing server logs.
type LogPrinter interface {
	Print(s string)
}

// Logger is a logger for logging server logs
type Logger struct {
	p LogPrinter
}

// NewLogger creates a new logger that prints server
// logs to the given printer.
func NewLogger(p LogPrinter) *Logger {
	return &Logger{p: p}
}

type stdLog struct{}

func (*stdLog) Print(s string) { log.Println(s) }

var theStdLog = new(stdLog)

// StdLogger returns the logger that logs to the default std log.
func StdLogger() *Logger {
	return &Logger{p: theStdLog}
}

// AltError logs the error and returns an alternative error with code.
func (l *Logger) AltError(err error, code, s string) error {
	if err == nil {
		return nil
	}
	l.p.Print(fmt.Sprintf("%s: %s", s, err))
	return errcode.Add(code, errors.New(s))
}

// AltInternal logs the error and returns an alternative internal error.
func (l *Logger) AltInternal(err error, s string) error {
	return l.AltError(err, errcode.Internal, s)
}

// AltInvalidArg logs the error and returns an alternative invalid arg error.
func (l *Logger) AltInvalidArg(err error, s string) error {
	return l.AltError(err, errcode.InvalidArg, s)
}

// Exit prints the error and exit the service.
func (l *Logger) Exit(err error) error {
	l.p.Print(err.Error())
	os.Exit(1)
	panic("unreachable")
}

// Print prints a message to the logger.
func (l *Logger) Print(args ...interface{}) {
	l.p.Print(fmt.Sprint(args...))
}

// Printf prints a formatted message to the logger.
func (l *Logger) Printf(f string, args ...interface{}) {
	l.p.Print(fmt.Sprintf(f, args...))
}

// Log logs the message to the logger if the logger is not nil.
func Log(l *Logger, s string) {
	if l == nil {
		return
	}
	l.Print(s)
}
