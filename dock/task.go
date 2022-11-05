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

package dock

import (
	"fmt"
	"log"

	"shanhu.io/pub/errcode"
)

// RunTask runs a command line.
func RunTask(c *Cont, line string) error {
	log.Println("#", line)

	exit, err := c.Exec(line)
	if err != nil {
		return errcode.Annotatef(err, "exec %q", line)
	}
	if exit != 0 {
		return fmt.Errorf("exit value: %d", exit)
	}
	return nil
}

// RunTasks runs a series of command lines. All commands must succeed and
// return 0 exit value.
func RunTasks(c *Cont, lines []string) error {
	for _, line := range lines {
		if err := RunTask(c, line); err != nil {
			return err
		}
	}
	return nil
}
