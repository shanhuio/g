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

package smake

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type context struct {
	dir     string
	modRoot string
	env     []string
	errLog  io.Writer
}

func newContext(gopath, modRoot, dir string) *context {
	var env []string
	for _, v := range []string{
		"PATH", "HOME", "SSH_AUTH_SOCK",
	} {
		if s := os.Getenv(v); s != "" {
			env = append(env, fmt.Sprintf("%s=%s", v, s))
		}
	}
	env = append(env, "GO111MODULE=on")
	if gopath != "" {
		env = append(env, fmt.Sprintf("GOPATH=%s", gopath))
	}

	return &context{
		dir:     dir,
		modRoot: modRoot,
		env:     env,
		errLog:  os.Stderr,
	}
}

func (c *context) workDir() string { return c.dir }

func (c *context) modRootDir() string { return c.modRoot }

func (c *context) atModRoot() bool { return c.dir == c.modRoot }

type execConfig struct {
	Stdout io.Writer
	Stderr io.Writer
}

func (c *context) execPkgs(
	pkgs []*relPkg, args []string, config *execConfig,
) error {
	line := strings.Join(args, " ")
	fmt.Println(line)

	if len(pkgs) > 0 {
		for _, pkg := range pkgs {
			args = append(args, pkg.rel)
		}
	}
	p, err := exec.LookPath(args[0])
	if err != nil {
		return err
	}

	var stdout, stderr io.Writer
	if config != nil {
		stdout, stderr = config.Stdout, config.Stderr
	}
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}
	cmd := exec.Cmd{
		Path:   p,
		Args:   args,
		Dir:    c.dir,
		Stdout: stdout,
		Stderr: stderr,
		Env:    c.env,
	}
	return cmd.Run()
}

func (c *context) logf(f string, args ...interface{}) {
	fmt.Fprintf(c.errLog, f, args...)
}

func (c *context) logln(args ...interface{}) {
	fmt.Fprintln(c.errLog, args...)
}
