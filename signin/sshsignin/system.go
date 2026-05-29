package sshsignin

import (
	"net"
	"os"
	osuser "os/user"

	"golang.org/x/crypto/ssh/agent"
	"shanhu.io/std/errcode"
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
