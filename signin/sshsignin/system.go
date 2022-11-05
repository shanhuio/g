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

package sshsignin

import (
	"net"
	"os"
	osuser "os/user"

	"golang.org/x/crypto/ssh/agent"
	"shanhu.io/pub/errcode"
)

// SysUser returns the default system user. It returns the value of
// SHANHU_USER if set, or system's current user name.
func SysUser() (string, error) {
	if u, ok := os.LookupEnv("SHANHU_USER"); ok && u != "" {
		return u, nil
	}
	cur, err := osuser.Current()
	if err != nil {
		return "", errcode.Annotate(err, "get current user")
	}
	return cur.Username, nil
}

// SysAgent returns the system's SSH agent by connecting to
// SSH_AUTH_SOCK.
func SysAgent() (agent.ExtendedAgent, error) {
	sock := os.Getenv("SSH_AUTH_SOCK")
	if sock == "" {
		return nil, errcode.Internalf("ssh agent socket not specified")
	}
	conn, err := net.Dial("unix", sock)
	if err != nil {
		return nil, errcode.Annotate(err, "dial agent")
	}
	return agent.NewClient(conn), nil
}
