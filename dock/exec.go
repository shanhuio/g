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
	"io"

	"shanhu.io/pub/errcode"
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
