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

package sniproxy

import (
	"io"
)

type remoteErr struct {
	code    int
	message string
}

func newRemoteErr(code int, err error) *remoteErr {
	return newRemoteErrString(code, err.Error())
}

func newRemoteErrString(code int, s string) *remoteErr {
	return &remoteErr{
		code:    code,
		message: s,
	}
}

func (e *remoteErr) Error() string { return e.message }

func (e *remoteErr) toError() error {
	if e == nil {
		return nil
	}
	if e.code == errEOF {
		return io.EOF
	}
	return e
}

func encodeRemoteErr(enc *encoder, err *remoteErr) {
	if err == nil {
		var empty remoteErr
		empty.encodeTo(enc)
		return
	}
	err.encodeTo(enc)
}

func (e *remoteErr) encodeTo(enc *encoder) {
	if e.code == 0 {
		enc.u64(0)
		return
	}
	enc.u64(uint64(e.code))
	enc.str(e.message)
}

func (e *remoteErr) decodeFrom(dec *decoder) {
	e.code = int(dec.u64())
	if e.code != 0 {
		e.message = dec.str()
	}
}

func decodeRemoteErr(dec *decoder) *remoteErr {
	err := new(remoteErr)
	err.decodeFrom(dec)
	if err.code == 0 {
		return nil
	}
	return err
}

var (
	remoteErrNotAccepting = newRemoteErrString(errAccept, "not accepting")
	remoteErrSiding       = newRemoteErrString(errSiding, "tunnel is siding")

	remoteErrSessionNotFound = newRemoteErrString(
		errSessionNotFound, "session not found",
	)
)
