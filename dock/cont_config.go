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
