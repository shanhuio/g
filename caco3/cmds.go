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

package caco3

import (
	"io"
	"os"
	"os/exec"

	"shanhu.io/pub/osutil"
)

type execJob struct {
	dir  string
	bin  string
	args []string
	out  io.Writer
}

func (j *execJob) command() *exec.Cmd {
	cmd := exec.Command(j.bin, j.args...)
	cmd.Dir = j.dir
	if j.out == nil {
		cmd.Stdout = os.Stdout
	} else {
		cmd.Stdout = j.out
	}
	cmd.Stderr = os.Stderr
	osutil.CmdCopyEnv(cmd, "HOME")
	osutil.CmdCopyEnv(cmd, "PATH")
	osutil.CmdCopyEnv(cmd, "SSH_AUTH_SOCK")
	return cmd
}

func runCmd(dir, bin string, args ...string) error {
	j := &execJob{
		dir:  dir,
		bin:  bin,
		args: args,
	}
	return j.command().Run()
}

func runCmdOutput(dir, bin string, args ...string) ([]byte, error) {
	j := &execJob{
		dir:  dir,
		bin:  bin,
		args: args,
	}
	cmd := j.command()
	cmd.Stdout = nil
	return cmd.Output()
}

func callCmd(dir, bin string, args ...string) (bool, error) {
	if err := runCmd(dir, bin, args...); err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return err.Success(), nil
		}
		return false, err
	}
	return true, nil
}
