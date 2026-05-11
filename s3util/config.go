package s3util

import (
	"fmt"
	"net/url"
)

// Config contains the configuration for connecting to an S3-compatible storage
// service.
type Config struct {
	Endpoint string
	Bucket   string
	BasePath string
}

// Host returns the host name of this endpoint config.
func (c *Config) Host() string {
	if c.Bucket != "" {
		return fmt.Sprintf("%s.%s", c.Bucket, c.Endpoint)
	}
	return c.Endpoint
}

// URL returns the S3 URL of this endpoint config.
func (c *Config) URL() *url.URL {
	return &url.URL{
		Scheme: "s3",
		Host:   c.Host(),
		Path:   c.BasePath,
	}
}

// Credential contains the credentials required for accessing
// S3-compatible storage.
type Credential struct {
	Key    string
	Secret string
}
