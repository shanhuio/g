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
)

// Mount types.
const (
	MountBind      = "bind"
	MountVolume    = "volume"
	MountTemp      = "tmpfs"
	MountNamedPipe = "npipe"
)

// ContMount specifies a containter bind option.
type ContMount struct {
	Host     string
	Cont     string
	ReadOnly bool
	Type     string // Default: "bind"
}

// ContDevice specifies devices to map into the container.
type ContDevice struct {
	Host        string
	Cont        string
	CgroupPerms string
}

func (b *ContMount) String() string {
	if b.ReadOnly {
		return fmt.Sprintf("%s:%s:ro", b.Host, b.Cont)
	}
	return fmt.Sprintf("%s:%s", b.Host, b.Cont)
}

// PortBind specifies a coutainer port bind option.
type PortBind struct {
	ContPort int
	HostIP   string
	HostPort int
}

// ContConfig contains the configuration for creating a container.
type ContConfig struct {
	Name          string
	Hostname      string
	Mounts        []*ContMount
	Devices       []*ContDevice
	TCPBinds      []*PortBind
	UDPBinds      []*PortBind
	Env           map[string]string
	Network       string
	Privileged    bool
	AlwaysRestart bool // Use restart policy "always".
	AutoRestart   bool // Use restart policy "unless-stopped".
	Cmd           []string
	WorkDir       string
	Labels        map[string]string

	JSONLogConfig *JSONLogConfig
}

// JSONLogConfig contains the logging option for a docker.
type JSONLogConfig struct {
	MaxSize string
	MaxFile int
}

// LimitedJSONLog returns a log config that has max-size=10 and
// max-file=3 as json-file logging options.
func LimitedJSONLog() *JSONLogConfig {
	return &JSONLogConfig{
		MaxSize: "10m",
		MaxFile: 3,
	}
}
