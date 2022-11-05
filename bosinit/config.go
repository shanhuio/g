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

package bosinit

import (
	"bytes"
	"io"

	yaml "gopkg.in/yaml.v2"
)

// Config is the root of the cloud-init configuration YAML.
type Config struct {
	Rancher    *Rancher     `yaml:",omitempty"`
	WriteFiles []*WriteFile `yaml:"write_files,omitempty"`

	SSHAuthorizedKeys []string `yaml:"ssh_authorized_keys,omitempty"`
}

func (c *Config) encodeInto(w io.Writer) error {
	enc := yaml.NewEncoder(w)
	if err := enc.Encode(c); err != nil {
		return err
	}
	return enc.Close()
}

// CloudConfig encodes the cloud-init config into YAML form.
// It also has "#cloud-config" shebang on the first line.
func (c *Config) CloudConfig() ([]byte, error) {
	buf := new(bytes.Buffer)
	io.WriteString(buf, "#cloud-config\n")
	if err := c.encodeInto(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// Encode encodes the cloud-init config into YAML form.
func (c *Config) Encode() ([]byte, error) {
	buf := new(bytes.Buffer)
	if err := c.encodeInto(buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// ParseConfig parses an Rancher/Burmilla OS init config.
func ParseConfig(bs []byte) (*Config, error) {
	c := new(Config)
	if err := yaml.Unmarshal(bs, c); err != nil {
		return nil, err
	}
	return c, nil
}
