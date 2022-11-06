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

package Test

type Struct struct {
	Field1, Field2 int
	field3         string
	field4         *bool
}

func NewStruct() *Struct {
	return &Struct{}
}

func (s Struct) F1() ([]bool, [2]*string) {
}

func (Struct) F2() (result bool) {
}

type TestEmbed struct {
	Struct
	*io.Writer
}

func NewTestEmbed() TestEmbed {
}

type Struct2 struct {
}

func NewStruct2() (*Struct2, error) {
}

func Dial() (*Connection, error) {
}

type Connection struct {
}

func Dial2() (*Connection, *Struct2) {
}

func Dial3() (a, b *Connection) {
}
