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

import (
	"fmt"
	"net/url"
	"strconv"
)

type contRestartPolicy struct {
	Name string `json:",omitempty"`
}

type contHostPortBinding struct {
	HostIP   string `json:"HostIp,omitempty"`
	HostPort string `json:",omitempty"`
}

type contLogConfig struct {
	Type   string
	Config map[string]string `json:",omitempty"`
}

func newContMount(m *ContMount) *ContMountInfo {
	typ := m.Type
	if typ == "" {
		typ = MountBind
	}
	return &ContMountInfo{
		Target:      m.Cont,
		Source:      m.Host,
		Type:        typ,
		ReadOnly:    m.ReadOnly,
		Consistency: "default",
	}
}

type contDevice struct {
	PathOnHost        string `json:",omitempty"`
	PathInContainer   string `json:",omitempty"`
	CgroupPermissions string `json:",omitempty"`
}

func newContDevice(d *ContDevice) *contDevice {
	ret := &contDevice{
		PathOnHost:        d.Host,
		PathInContainer:   d.Cont,
		CgroupPermissions: d.CgroupPerms,
	}
	if ret.PathInContainer == "" {
		ret.PathInContainer = ret.PathOnHost
	}
	return ret
}

type contHostConfig struct {
	PortBindings map[string][]*contHostPortBinding `json:",omitempty"`

	RestartPolicy *contRestartPolicy `json:",omitempty"`
	NetworkMode   string             `json:",omitempty"`
	LogConfig     *contLogConfig     `json:",omitempty"`
	Mounts        []*ContMountInfo   `json:",omitempty"`
	Devices       []*contDevice      `json:",omitempty"`

	Privileged bool `json:",omitempty"`
}

// CreateCont creates a new container with the given config.
func CreateCont(c *Client, image string, config *ContConfig) (*Cont, error) {
	req := struct {
		Image        string
		Env          []string            `json:",omitempty"`
		Hostname     string              `json:",omitempty"`
		Cmd          []string            `json:",omitempty"`
		HostConfig   *contHostConfig     `json:",omitempty"`
		ExposedPorts map[string]struct{} `json:",omitempty"`
		Labels       map[string]string   `json:",omitempty"`
		WorkingDir   string              `json:",omitempty"`
	}{
		Image: image,
	}

	q := make(url.Values)

	if config != nil {
		req.WorkingDir = config.WorkDir
		req.Hostname = config.Hostname

		hc := &contHostConfig{Privileged: config.Privileged}
		req.HostConfig = hc
		for _, m := range config.Mounts {
			hc.Mounts = append(hc.Mounts, newContMount(m))
		}
		for _, d := range config.Devices {
			hc.Devices = append(hc.Devices, newContDevice(d))
		}
		if len(config.TCPBinds) > 0 || len(config.UDPBinds) > 0 {
			hc.PortBindings = make(map[string][]*contHostPortBinding)
			req.ExposedPorts = make(map[string]struct{})
			for _, binds := range []struct {
				proto string
				binds []*PortBind
			}{
				{proto: "tcp", binds: config.TCPBinds},
				{proto: "udp", binds: config.UDPBinds},
			} {
				proto := binds.proto
				for _, b := range binds.binds {
					hb := &contHostPortBinding{
						HostIP:   b.HostIP,
						HostPort: strconv.Itoa(b.HostPort),
					}
					k := fmt.Sprintf("%d/%s", b.ContPort, proto)
					hc.PortBindings[k] = append(hc.PortBindings[k], hb)
					req.ExposedPorts[k] = struct{}{}
				}
			}
		}
		hc.NetworkMode = config.Network
		if j := config.JSONLogConfig; j != nil {
			lc := &contLogConfig{
				Type:   "json-file",
				Config: make(map[string]string),
			}
			if j.MaxSize != "" {
				lc.Config["max-size"] = j.MaxSize
			}
			if j.MaxFile > 0 {
				lc.Config["max-file"] = fmt.Sprint(j.MaxFile)
			}
			hc.LogConfig = lc
		}

		if config.AlwaysRestart {
			hc.RestartPolicy = &contRestartPolicy{Name: "always"}
		} else if config.AutoRestart {
			hc.RestartPolicy = &contRestartPolicy{Name: "unless-stopped"}
		}
		if len(config.Env) > 0 {
			req.Env = unmapEnv(config.Env)
		}
		if len(config.Cmd) > 0 {
			req.Cmd = append(req.Cmd, config.Cmd...)
		}

		if config.Name != "" {
			q.Add("name", config.Name)
		}

		if len(config.Labels) > 0 {
			req.Labels = make(map[string]string)
			for k, v := range config.Labels {
				req.Labels[k] = v
			}
		}
	}

	var resp struct {
		ID string `json:"Id"`
	}
	if err := c.call("containers/create", q, &req, &resp); err != nil {
		return nil, err
	}
	return &Cont{c: c, id: resp.ID}, nil
}
