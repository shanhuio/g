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

package caco3

import (
	"shanhu.io/g/jsonutil"
)

type dockerSum struct {
	Repo   string
	Tag    string
	ID     string
	Origin string `json:",omitempty"`
}

func newDockerSum(repo, tag, id string) *dockerSum {
	return &dockerSum{
		Repo: repo,
		Tag:  tag,
		ID:   id,
	}
}

func dockerSumOut(name string) string { return name + ".dockersum" }

func dockerTarOut(name string) string { return name + ".tar.gz" }

func loadDockerSum(f string) (*dockerSum, error) {
	sum := new(dockerSum)
	if err := jsonutil.ReadFile(f, sum); err != nil {
		return nil, err
	}
	return sum, nil
}
