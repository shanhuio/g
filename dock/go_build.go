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

package dock

import (
	"fmt"
	"path"
)

// ShanhuBuild returns the golang build for shanhu releases.
func ShanhuBuild(name, sshTar string) *GoBuild {
	return &GoBuild{
		Bin:     path.Base(name),
		Git:     "git@bitbucket.org:shanhuio/p2.git",
		RepoPkg: "shanhu.io/p2",
		Pkg:     path.Join("shanhu.io/p2/", name),
		SSHTar:  sshTar,
	}
}

// GoBuild is a golang build.
type GoBuild struct {
	Git     string
	RepoPkg string
	Pkg     string
	Bin     string
	SSHTar  string
}

func goSrcPath(s string) string {
	return path.Join("/go/src", s)
}

// Run runs a Go language build job.
func (b *GoBuild) Run(client *Client) error {
	c, err := CreateCont(client, "shanhu/builder", nil)
	if err != nil {
		return err
	}
	defer c.ForceRemove()

	if err := c.Start(); err != nil {
		return err
	}

	if b.SSHTar != "" {
		if err := RunTask(c, "mkdir -p -m700 /root/.ssh"); err != nil {
			return err
		}
		if err := CopyInTarFile(c, b.SSHTar, "/root/.ssh"); err != nil {
			return err
		}
	}

	pkgPath := path.Join("/go/src", b.RepoPkg)

	if err := RunTasks(c, []string{
		fmt.Sprintf("mkdir -p %s", pkgPath),
		fmt.Sprintf("git clone --depth 1 %s %s", b.Git, pkgPath),
		fmt.Sprintf("go install -v %s", b.Pkg),
	}); err != nil {
		return err
	}

	// Copy the binary out.
	if err := c.CopyOut(fmt.Sprintf("/go/bin/%s", b.Bin), "."); err != nil {
		return err
	}

	return c.Drop()
}
