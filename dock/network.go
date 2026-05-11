package dock

import (
	"path"

	"shanhu.io/g/errcode"
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
