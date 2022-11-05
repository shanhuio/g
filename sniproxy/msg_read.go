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

type readRequest struct {
	session uint64
	maxRead int
}

func (m *readRequest) encodeTo(enc *encoder) {
	enc.u64(m.session)
	enc.u64(uint64(m.maxRead))
}

func (m *readRequest) decodeFrom(dec *decoder) {
	m.session = dec.u64()
	m.maxRead = int(dec.u64())
}

type readResponse struct {
	bytes []byte
	err   *remoteErr
}

func (m *readResponse) encodeTo(enc *encoder) {
	enc.bytes(m.bytes)
	encodeRemoteErr(enc, m.err)
}

func (m *readResponse) decodeFrom(dec *decoder) {
	m.bytes = dec.bytes(m.bytes)
	m.err = decodeRemoteErr(dec)
}
