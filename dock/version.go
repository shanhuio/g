package dock

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
