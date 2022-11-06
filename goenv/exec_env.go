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

package goenv

import (
	"go/build"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"shanhu.io/pub/osutil"
)

// ExecJob is an execution job to be executed in an execution environment.
type ExecJob struct {
	Dir  string
	Name string
	Args []string
}

// ExecEnv is an execution environment for executing a Go language related
// command.
type ExecEnv struct {
	gopath   string
	buildCtx *build.Context
	pipe     io.Writer
}

// NewExecEnv creates a new execution environment for a particular GOPATH.
func NewExecEnv(gopath string) *ExecEnv {
	ctx := build.Default
	ctx.GOPATH = gopath
	return &ExecEnv{
		gopath:   gopath,
		buildCtx: &ctx,
	}
}

// GOPATH returns the GOPATH for this environment.
func (env *ExecEnv) GOPATH() string { return env.gopath }

// BindPipe will forward stdout and stderr to the given writer,
// rather than os.Stdout and os.Stderr.
func (env *ExecEnv) BindPipe(w io.Writer) { env.pipe = w }

// Context returns the build context that is used by this environment.
func (env *ExecEnv) Context() *build.Context {
	return env.buildCtx
}

// IsDir checks if p exists as a directory under the GOPATH.
func (env *ExecEnv) IsDir(p string) (bool, error) {
	if env.gopath != "" {
		p = filepath.Join(env.gopath, p)
	}
	stat, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return stat.IsDir(), nil
}

// Cmd creates an execution process using a given execution job.
func (env *ExecEnv) Cmd(j *ExecJob) *exec.Cmd {
	ret := exec.Command(j.Name, j.Args...)
	if env.gopath == "" {
		if j.Dir == "" {
			ret.Dir = "/"
		} else {
			ret.Dir = j.Dir
		}
	} else {
		if j.Dir == "" {
			ret.Dir = env.gopath
		} else {
			ret.Dir = filepath.Join(env.gopath, j.Dir)
		}
	}
	for _, k := range []string{
		"HOME", "PATH", "SSH_AUTH_SOCK",
	} {
		osutil.CmdCopyEnv(ret, k)
	}
	osutil.CmdAddEnv(ret, "GO111MODULE", "off")
	if env.gopath != "" {
		osutil.CmdAddEnv(ret, "GOPATH", env.gopath)
	}
	return ret
}

// PipedCmd creates an execution process using a given execution job similar to
// Cmd but also forward through the Stdout and Stderr.
func (env *ExecEnv) PipedCmd(j *ExecJob) *exec.Cmd {
	ret := env.Cmd(j)

	if env.pipe != nil {
		ret.Stdout = env.pipe
		ret.Stderr = env.pipe
	} else {
		ret.Stdout = os.Stdout
		ret.Stderr = os.Stderr
	}
	return ret
}

// Exec executes a process in the environment.
func (env *ExecEnv) Exec(dir, name string, args ...string) error {
	cmd := env.PipedCmd(&ExecJob{
		Dir: dir, Name: name, Args: args,
	})
	return cmd.Run()
}

// Call executes a process in the environment and returns true if the
// process ends and exits with a success exit code, false if the process
// ends and exists with a non-success exit code.
func (env *ExecEnv) Call(dir, name string, args ...string) (bool, error) {
	cmd := env.Cmd(&ExecJob{
		Dir: dir, Name: name, Args: args,
	})
	if err := cmd.Run(); err != nil {
		if err, ok := err.(*exec.ExitError); ok {
			return err.Success(), nil
		}
		return false, err
	}
	return true, nil
}

// StrOut executes a process in the environment and returns the output as a
// string.
func (env *ExecEnv) StrOut(dir, name string, args ...string) (string, error) {
	cmd := env.Cmd(&ExecJob{
		Dir:  dir,
		Name: name,
		Args: args,
	})
	if env.pipe != nil {
		cmd.Stderr = env.pipe
	} else {
		cmd.Stderr = os.Stderr
	}
	bs, err := cmd.Output()
	return string(bs), err
}
