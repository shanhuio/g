package osutil

import (
	"fmt"
	"os"
	"os/exec"
)

// CmdAddEnv adds an environment variable to cmd and returns true. If v is an
// empty string, nothing is added, and it returns false.
func CmdAddEnv(cmd *exec.Cmd, k, v string) bool {
	if v == "" {
		return false
	}
	cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
	return true
}

// CmdCopyEnv copies the value of environment variable k to cmd. If the value
// is empty, returns false; otherwise returns true.
func CmdCopyEnv(cmd *exec.Cmd, k string) bool {
	return CmdAddEnv(cmd, k, os.Getenv(k))
}
