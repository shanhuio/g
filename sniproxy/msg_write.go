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

type writeRequest struct {
	session uint64
	bytes   []byte
}

func (m *writeRequest) encodeTo(enc *encoder) {
	enc.u64(m.session)
	enc.bytes(m.bytes)
}

func (m *writeRequest) decodeFrom(dec *decoder) {
	m.session = dec.u64()
	m.bytes = dec.bytes(m.bytes)
}

type writeResponse struct {
	written int
	err     *remoteErr
}

func (m *writeResponse) encodeTo(enc *encoder) {
	enc.u64(uint64(m.written))
	encodeRemoteErr(enc, m.err)
}

func (m *writeResponse) decodeFrom(dec *decoder) {
	m.written = int(dec.u64())
	m.err = decodeRemoteErr(dec)
}
