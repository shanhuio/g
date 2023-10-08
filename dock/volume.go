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

	"shanhu.io/g/errcode"
)

// VolumeConfig is the configuration of a new volume.
type VolumeConfig struct {
	Labels map[string]string
	Driver string
}

// CreateVolume creates a new volume
func CreateVolume(
	c *Client, name string, config *VolumeConfig,
) (string, error) {
	if config == nil {
		config = &VolumeConfig{}
	}
	var req = struct {
		Name   string
		Labels map[string]string `json:",omitempty"`
		Driver string            `json:",omitempty"`
	}{
		Name:   name,
		Labels: config.Labels,
		Driver: config.Driver,
	}

	var resp struct{ Name string }
	if err := c.call("volumes/create", nil, req, &resp); err != nil {
		return "", err
	}
	return resp.Name, nil
}

// VolumeInfo contains the information of a volume
type VolumeInfo struct {
	Name       string
	Driver     string
	Mountpoint string
	Labels     map[string]string
}

// InspectVolume inspects a volume.
func InspectVolume(c *Client, name string) (*VolumeInfo, error) {
	info := new(VolumeInfo)
	if err := c.jsonGet(path.Join("volumes", name), nil, info); err != nil {
		return nil, err
	}
	return info, nil
}

// RemoveVolume deletes a volume.
func RemoveVolume(c *Client, name string) error {
	return c.del(path.Join("volumes", name), nil)
}

// CreateVolumeIfNotExist creates the volume if the volume of the given name
// does not exist.
func CreateVolumeIfNotExist(
	c *Client, name string, config *VolumeConfig,
) (string, error) {
	info, err := InspectVolume(c, name)
	if err != nil {
		if errcode.IsNotFound(err) {
			return CreateVolume(c, name, config)
		}
		return "", errcode.Annotate(err, "inspect volume")
	}
	return info.Name, nil
}

// ListVolumesWithLabel lists all volumes with the specific label.
func ListVolumesWithLabel(c *Client, label string) ([]*VolumeInfo, error) {
	filters, err := labelFilters(label)
	if err != nil {
		return nil, err
	}
	q := singleQuery("filters", filters)
	var resp struct {
		Volumes []*VolumeInfo
	}
	if err := c.jsonGet("volumes", q, &resp); err != nil {
		return nil, err
	}
	return resp.Volumes, nil
}
