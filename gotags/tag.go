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

package gotags

import (
	"bytes"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// Tag represents a single tag.
type Tag struct {
	Name    string
	File    string
	Address string
	Type    TagType
	Fields  map[TagField]string
}

// TagField represents a single field in a tag line.
type TagField string

// Tag fields.
const (
	Access        TagField = "access"
	Signature     TagField = "signature"
	TypeField     TagField = "type"
	ReceiverType  TagField = "ctype"
	Line          TagField = "line"
	InterfaceType TagField = "ntype"
	Language      TagField = "language"
	ExtraTags     TagField = "extraTag"
)

// TagType represents the type of a tag in a tag line.
type TagType string

// Tag types.
const (
	Package     TagType = "p"
	Import      TagType = "i"
	Constant    TagType = "c"
	Variable    TagType = "v"
	Type        TagType = "t"
	Interface   TagType = "n"
	Field       TagType = "w"
	Embedded    TagType = "e"
	Method      TagType = "m"
	Constructor TagType = "r"
	Function    TagType = "f"
)

// NewTag creates a new Tag.
func NewTag(name, file string, line int, tagType TagType) Tag {
	l := strconv.Itoa(line)
	return Tag{
		Name:    name,
		File:    file,
		Address: l,
		Type:    tagType,
		Fields:  map[TagField]string{Line: l},
	}
}

// The tags file format string representation of this tag.
func (t Tag) String() string {
	var b bytes.Buffer

	b.WriteString(t.Name)
	b.WriteByte('\t')
	b.WriteString(t.File)
	b.WriteByte('\t')
	b.WriteString(t.Address)
	b.WriteString(";\"\t")
	b.WriteString(string(t.Type))
	b.WriteByte('\t')

	fields := make([]string, 0, len(t.Fields))
	i := 0
	for k, v := range t.Fields {
		if len(v) == 0 {
			continue
		}
		fields = append(fields, fmt.Sprintf("%s:%s", k, v))
		i++
	}

	sort.Sort(sort.StringSlice(fields))
	b.WriteString(strings.Join(fields, "\t"))

	return b.String()
}
