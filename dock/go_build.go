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
