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

func newRequestMessage(t uint8) (decoderFrom, bool) {
	switch t {
	case msgShutdown:
		return nil, true
	case msgHello:
		return new(helloRequest), true
	case msgDial:
		return new(dialRequest), true
	case msgDialSide:
		return new(dialSideRequest), true
	case msgDialSide2:
		return new(dialSide2Request), true
	case msgRead:
		return new(readRequest), true
	case msgWrite:
		return new(writeRequest), true
	case msgClose:
		return new(closeRequest), true
	default:
		return nil, false
	}
}
