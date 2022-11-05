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

type statusRequest struct {
	session uint64
}

func (m *statusRequest) encodeTo(enc *encoder) {
	enc.u64(m.session)
}

func (m *statusRequest) decodeFrom(dec *decoder) {
	m.session = dec.u64()
}

type statusResponse struct {
	uptime       uint64
	totalRead    uint64
	totalWritten uint64
}

func (m *statusResponse) encodeTo(enc *encoder) {
	enc.u64(m.uptime)
	enc.u64(m.totalRead)
	enc.u64(m.totalWritten)
}

func (m *statusResponse) decodeFrom(dec *decoder) {
	m.uptime = dec.u64()
	m.totalRead = dec.u64()
	m.totalWritten = dec.u64()
}
