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

package dock

// VersionInfo saves the response of the `/version` docker endpoint.
type VersionInfo struct {
	Version       string `json:"Version"`
	APIVersion    string `json:"ApiVersion"`
	OS            string `json:"Os"`
	KernelVersion string `json:"KernelVersion"`
	GoVersion     string `json:"GoVersion"`
	Arch          string `json:"Arch"`
}

// Version returns the version info of the docker service.
func Version(c *Client) (*VersionInfo, error) {
	info := new(VersionInfo)
	if err := c.jsonGet("version", nil, info); err != nil {
		return nil, err
	}
	return info, nil
}
