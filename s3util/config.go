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
