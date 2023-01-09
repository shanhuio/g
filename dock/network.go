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

import (
	"path"

	"shanhu.io/pub/errcode"
)

// CreateNetwork creates a network with the given name.
func CreateNetwork(c *Client, name string) error {
	req := struct{ Name string }{Name: name}
	return c.jsonPost("networks/create", nil, &req, nil)
}

// RemoveNetwork removes a network of the given name.
func RemoveNetwork(c *Client, name string) error {
	return c.del(path.Join("networks", name), nil)
}

// NetworkInfo contains information of a docker network.
type NetworkInfo struct {
	Name   string
	ID     string `json:"Id"`
	Driver string
	IPAM   *NetworkIPAM `json:",omitempty"`
}

// NetworkIPAM is the configuration section of IP address management.
type NetworkIPAM struct {
	Config []*NetworkIPAMConfig `json:",omitempty"`
}

// NetworkIPAMConfig is the IP address management config entry.
type NetworkIPAMConfig struct {
	Subnet  string
	Gateway string
}

// InspectNetwork inspects a network.
func InspectNetwork(c *Client, name string) (*NetworkInfo, error) {
	info := new(NetworkInfo)
	if err := c.jsonGet(path.Join("networks", name), nil, info); err != nil {
		return nil, err
	}
	return info, nil
}

// HasNetwork checks if a network exists.
func HasNetwork(c *Client, name string) (bool, error) {
	if _, err := InspectNetwork(c, name); err != nil {
		if errcode.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
