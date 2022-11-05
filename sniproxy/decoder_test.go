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
	"testing"

	"bytes"
	"math/rand"
	"reflect"
)

func TestDecoder(t *testing.T) {
	bigTrunk := make([]byte, 1024*32)
	if _, err := rand.Read(bigTrunk); err != nil {
		t.Fatal("generate big trunk data", err)
	}

	type testCase struct {
		en   func(enc *encoder)
		de   func(dec *decoder) interface{}
		want interface{}
	}

	for _, c := range []*testCase{{
		en:   func(enc *encoder) { enc.u64(32) },
		de:   func(dec *decoder) interface{} { return dec.u64() },
		want: uint64(32),
	}, {
		en:   func(enc *encoder) { enc.str("hello") },
		de:   func(dec *decoder) interface{} { return dec.str() },
		want: "hello",
	}, {
		en:   func(enc *encoder) { enc.u8(3) },
		de:   func(dec *decoder) interface{} { return dec.u8() },
		want: uint8(3),
	}, {
		en:   func(enc *encoder) { enc.bytes([]byte{0xfa, 0xce}) },
		de:   func(dec *decoder) interface{} { return dec.bytes(nil) },
		want: []byte{0xfa, 0xce},
	}, {
		en:   func(enc *encoder) { enc.bytes(bigTrunk) },
		de:   func(dec *decoder) interface{} { return dec.bytes(nil) },
		want: bigTrunk,
	}} {
		buf := new(bytes.Buffer)
		enc := newEncoder(buf)
		c.en(enc)
		if err := enc.Err(); err != nil {
			t.Errorf("encode got error: %s", err)
		}

		dec := newDecoder(buf)
		got := c.de(dec)
		dec.end()

		if err := dec.Err(); err != nil {
			t.Errorf("decode got error: %s", err)
		}
		if !reflect.DeepEqual(got, c.want) {
			t.Errorf("want %v, got %v", c.want, got)
		}
	}

}
