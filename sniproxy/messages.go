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

// Message types. For backwards compatibility, new message types must be added
// in the back.
const (
	msgShutdown = iota
	msgHello
	msgDial
	msgWrite
	msgRead
	msgStatus
	msgClose
	msgShutdownHint

	msgDialSide
	msgDialSide2
)

type encoderTo interface {
	encodeTo(enc *encoder)
}

type decoderFrom interface {
	decodeFrom(dec *decoder)
}

type message interface {
	encoderTo
	decoderFrom
}

const (
	errUnknown = iota + 1
	errUnknownType
	errBug
	errAccept
	errSessionNotFound
	errRead
	errWrite
	errClose
	errInternal
	errEOF
	errSiding
)
