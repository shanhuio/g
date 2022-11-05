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

const (
	ruleFileSet     = "file_set"
	ruleBundle      = "bundle"
	ruleDockerPull  = "docker_pull"
	ruleDockerBuild = "docker_build"
	ruleDockerRun   = "docker_run"
	ruleDownload    = "download"
)

// FileSet selects a set of files.
type FileSet struct {
	Name string

	// The list of files to include in the fileset.
	Files []string `json:",omitempty"`

	// Selects a set of source input files.
	Select []string `json:",omitempty"`

	// Ignores a set of source input files after selection.
	Ignore []string `json:",omitempty"`

	// Merge in other file sets
	Include []string `json:",omitempty"`
}

// Bundle is a set of build rules in a bundle. A bundle has no build actions;
// it just group rules together.
type Bundle struct {
	// Name of the fule
	Name string

	// Other rule names.
	Deps []string
}

// DockerPull is a rule to pull down a docker container image.
type DockerPull struct {
	Name      string
	Pull      string `json:",omitempty"`
	Digest    string `json:",omitempty"`
	OutputTar bool   `json:",omitempty"`
}

// DockerBuild is a rule to build a docker container image.
type DockerBuild struct {
	Name         string
	Dockerfile   string   `json:",omitempty"`
	From         []string `json:",omitempty"`
	Input        []string `json:",omitempty"`
	ArchiveInput []string `json:",omitempty"`
	PrefixDir    string   `json:",omitempty"`
	Args         []string `json:",omitempty"`
	OutputTar    bool     `json:",omitempty"`
}

// DockerRun is a rule to run a command inside a docker container image.
type DockerRun struct {
	Name    string
	Image   string
	User    string   `json:",omitempty"`
	Envs    []string `json:",omitempty"`
	WorkDir string   `json:",omitempty"`

	MountWorkspace string `json:",omitempty"`

	Command []string `json:",omitempty"`

	// Map from input to file inside the container.
	Input map[string]string

	// Map files from zip or tarball archives to directories
	// inside the container.
	ArchiveInput map[string]string

	// Map from output path to file inside the container.
	Output map[string]string `json:",omitempty"`

	// Extra dependencies.
	Deps []string `json:",omitempty"`
}

// Download is a rule to download an artifact from the Internet.
type Download struct {
	Name     string
	URL      string
	Checksum string
	Output   string
}
