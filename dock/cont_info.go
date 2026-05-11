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
