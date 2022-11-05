// Copyright (C) 2021  Shanhu Tech Inc.
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

package dags

// Repo is the overview dependency structure of a repository.
type Repo struct {
	Name     string
	RepoTopo *M
	PkgTopos map[string]*M
}

// NewRepo creates an empty overview for a repo.
func NewRepo(name string) *Repo {
	return &Repo{
		Name:     name,
		PkgTopos: make(map[string]*M),
	}
}
