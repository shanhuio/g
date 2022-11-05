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
	"shanhu.io/pub/jsonx"
	"shanhu.io/pub/lexing"
)

// Workspace is the structure of the build.jsonx file. It specifies how
// to build a project.
type Workspace struct {
	RepoMap *RepoMap
}

// GitRemote defines a set of remote URLs for a given name. It provides a more
// consistent remote setup for the repositories in the workspace.
type GitRemote struct {
	Name string
	URL  map[string]string
}

// RepoMap contains the list of repos to clone down.
type RepoMap struct {
	GitHosting   map[string]string `json:",omitempty"`
	Src          map[string]string
	ExtraRemotes []*GitRemote `json:",omitempty"`
}

func readWorkspace(f string) (*Workspace, []*lexing.Error) {
	tm := func(t string) interface{} {
		switch t {
		case "repo_map":
			return new(RepoMap)
		}
		return nil
	}
	entries, errs := jsonx.ReadSeriesFile(f, tm)
	if errs != nil {
		return nil, errs
	}

	ws := new(Workspace)
	for _, entry := range entries {
		switch v := entry.V.(type) {
		case *RepoMap:
			ws.RepoMap = v
		}
	}
	return ws, nil
}

// RepoSums records the checkums and git commits of a build.
type RepoSums struct {
	RepoCommits map[string]string
}

// ReadRepoSums reads in the workspaces's repo checksum file.
func ReadRepoSums(f string) (*RepoSums, error) {
	b := new(RepoSums)
	if err := jsonx.ReadFile(f, b); err != nil {
		return nil, err
	}
	return b, nil
}

// SaveRepoSums saves sums to f.
func SaveRepoSums(f string, sums *RepoSums) error {
	return jsonx.WriteFile(f, sums)
}
