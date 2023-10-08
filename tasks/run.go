// Copyright (C) 2023  Shanhu Tech Inc.
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

package tasks

import (
	"fmt"
	"log"
	"path"

	"shanhu.io/g/creds"
)

// Run issues a list of tasks to a particular server.
func Run(server, prefix string, tasks []string) error {
	c, err := creds.Dial(server)
	if err != nil {
		return err
	}

	if len(tasks) == 0 {
		// list available tasks
		var names []string
		p := path.Join(prefix, "help")
		if err := c.JSONCall(p, nil, &names); err != nil {
			return err
		}

		for _, name := range names {
			fmt.Println(name)
		}
		return nil
	}

	for _, t := range tasks {
		log.Println(t)
		if err := c.Poke(path.Join(prefix, t)); err != nil {
			return err
		}
	}

	return nil
}
