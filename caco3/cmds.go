package caco3

import (
	"io"
	"os"
	"os/exec"

	"shanhu.io/g/osutil"
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
