package dock

import (
	"io"

	"shanhu.io/g/errcode"
)

// ExecSetup contains the setup arguments to run a command inside a running
// container.
type ExecSetup struct {
	Cmd        []string
	Env        []string `json:",omitempty"`
	User       string   `json:",omitempty"`
	WorkingDir string   `json:",omitempty"`

	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func contExec(c *Client, id string, s *ExecSetup) (int, error) {
	createReq := struct {
		Cmd        []string
		Env        []string `json:",omitempty"`
		User       string   `json:",omitempty"`
		WorkingDir string   `json:",omitempty"`

		AttachStdout bool
		AttachStderr bool
	}{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          s.Cmd,
		Env:          s.Env,
		User:         s.User,
		WorkingDir:   s.WorkingDir,
	}

	var createResp struct {
		ID string `json:"Id"`
	}

	if err := c.call(
		contPath(id, "exec"), nil, &createReq, &createResp,
	); err != nil {
		return 0, err
	}

	execID := createResp.ID
	sink := newLogSink(s.Stdout, s.Stderr)
	if err := c.jsonPost(
		execPath(execID, "start"), nil, struct{}{}, sink,
	); err != nil {
		return 0, err
	}

	if err := sink.waitDone(); err != nil {
		return 0, err
	}

	var inspectResp struct {
		ExitCode int
		Running  bool
	}
	if err := c.jsonGet(
		execPath(execID, "json"), nil, &inspectResp,
	); err != nil {
		return 0, err
	}

	if inspectResp.Running {
		return 0, errcode.Internalf("process still running")
	}
	return inspectResp.ExitCode, nil
}
