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

type dialRequest struct {
}

func (m *dialRequest) encodeTo(*encoder)   {}
func (m *dialRequest) decodeFrom(*decoder) {}

type dialResponse struct {
	session uint64
	err     *remoteErr
}

func (m *dialResponse) encodeTo(enc *encoder) {
	enc.u64(m.session)
	encodeRemoteErr(enc, m.err)
}

func (m *dialResponse) decodeFrom(dec *decoder) {
	m.session = dec.u64()
	m.err = decodeRemoteErr(dec)
}

type dialSideRequest struct {
	session uint64
	key     uint64
	token   string
}

func (m *dialSideRequest) encodeTo(enc *encoder) {
	enc.u64(m.session)
	enc.u64(m.key)
	enc.str(m.token)
}

func (m *dialSideRequest) decodeFrom(dec *decoder) {
	m.session = dec.u64()
	m.key = dec.u64()
	m.token = dec.str()
}

type dialSide2Request struct {
	session uint64
	key     uint64
	token   string
	tcpAddr string
}

func (m *dialSide2Request) encodeTo(enc *encoder) {
	enc.u64(m.session)
	enc.u64(m.key)
	enc.str(m.token)
	enc.str(m.tcpAddr)
}

func (m *dialSide2Request) decodeFrom(dec *decoder) {
	m.session = dec.u64()
	m.key = dec.u64()
	m.token = dec.str()
	m.tcpAddr = dec.str()
}
