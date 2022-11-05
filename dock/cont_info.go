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

// ContInfo is the inspection result of a container.
type ContInfo struct {
	ID     string `json:"Id"`
	Image  string
	Config struct {
		Image    string
		Hostname string
		Labels   map[string]string
	}
	State struct {
		ExitCode int
		Error    string
		Running  bool
	}
	HostConfig struct {
		Mounts []*ContMountInfo
	}
}

// ContMountInfo is the information of a mount in the container.
type ContMountInfo struct {
	Type        string
	Target      string
	Source      string
	ReadOnly    bool   `json:",omitempty"`
	Consistency string `json:",omitempty"`
}
