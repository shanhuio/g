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

package caco3

import (
	"log"
	"strings"
)

type gitRemote struct {
	name  string
	git   string
	fetch bool
	push  bool
}

func listRemotes(dir string) (map[string]*gitRemote, error) {
	output, err := runCmdOutput(dir, "git", "remote", "-v")
	if err != nil {
		return nil, err
	}

	remotes := make(map[string]*gitRemote)

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 3 {
			name := fields[0]
			git := fields[1]
			method := fields[2]
			remote, ok := remotes[name]
			if ok {
				if git != remote.git {
					log.Printf(
						"inconsistent remote url for %q: %q",
						name, line,
					)
					continue
				}
			} else {
				remote = &gitRemote{
					name: name,
					git:  git,
				}
				remotes[name] = remote
			}
			if method == "(fetch)" {
				remote.fetch = true
			} else if method == "(push)" {
				remote.push = true
			} else {
				log.Printf("unknown git remote method: %q", line)
			}
		} else {
			log.Printf("weird git remote line: %q", line)
		}
	}

	return remotes, nil
}
