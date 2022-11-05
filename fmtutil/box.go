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

package fmtutil

import (
	"strings"
)

func isIndent(r rune) bool {
	return r == ' ' || r == '\t'
}

func indentWithSpace(s string) string {
	ret := ""
	for i, c := range s {
		if c == ' ' {
			ret += " "
		} else if c == '\t' {
			ret += "    "
		} else {
			return ret + s[i:]
		}
	}

	return ret
}

func allIndents(s string) bool {
	for _, r := range s {
		if !isIndent(r) {
			return false
		}
	}
	return true
}

func indentCount(s string) int {
	for i, r := range s {
		if !isIndent(r) {
			return i
		}
	}
	return len(s)
}

// Box removes the common indent of a multi-line block.
func Box(s string) string { return box(s, false) }

// BoxSpaceIndent converts tab indents into 4 spaces and removes the common
// indent of a multi-line block.
func BoxSpaceIndent(s string) string { return box(s, true) }

func box(s string, spaceIndent bool) string {
	lines := strings.Split(s, "\n")
	nline := len(lines)

	// convert indent into space, make blank lines empty.
	for i := range lines {
		if i == 0 {
			continue
		}
		line := lines[i]
		if spaceIndent {
			line = indentWithSpace(line)
		}
		if allIndents(line) {
			line = ""
		}
		lines[i] = line
	}

	if nline > 1 {
		minIndentCount := indentCount(lines[1])
		for i, line := range lines {
			if i <= 1 {
				continue
			}
			if line == "" {
				continue
			}

			n := indentCount(line)
			if n < minIndentCount {
				minIndentCount = n
			}
		}

		for i := range lines {
			if i == 0 {
				continue
			}
			if lines[i] == "" {
				continue
			}
			// trim prefix of minIndent
			lines[i] = lines[i][minIndentCount:]
		}
	}

	return strings.Join(lines, "\n")
}
