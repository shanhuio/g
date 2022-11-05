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

package bosinit

import (
	"strconv"
)

// WriteFile specifies a file to be written onto the file system.
type WriteFile struct {
	Path        string
	Permissions string
	Owner       string
	Content     string
}

// FilePerm gererates file permission string for use in WriteFile.
func FilePerm(m int) string {
	return "0" + strconv.FormatInt(int64(m), 8)
}

// RCLocal creates cloud-init entry to add /etc/rc.local file on the target.
func RCLocal(content string) *WriteFile {
	return &WriteFile{
		Path:        "/etc/rc.local",
		Permissions: FilePerm(0744),
		Owner:       "root",
		Content:     content,
	}
}
