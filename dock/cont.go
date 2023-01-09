// Copyright (C) 2023  Shanhu Tech Inc.
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
	"net/http"
	"net/url"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/httputil"
	"shanhu.io/pub/strtoken"
)

// Cont wraps a container.
type Cont struct {
	c  *Client
	id string
}

// NewCont creates a new container handle by ID or name.
func NewCont(c *Client, id string) *Cont {
	return &Cont{c: c, id: id}
}

// ID returns the container's ID.
func (c *Cont) ID() string { return c.id }

func (c *Cont) path(m string) string { return contPath(c.id, m) }

// Exec executes a command line. Returns the exit value or any error.
func (c *Cont) Exec(line string) (int, error) {
	args, errs := strtoken.Parse(line)
	if len(errs) > 0 {
		return 0, errs[0]
	}
	return c.ExecArgs(args)
}

// ExecArgs executes a command with the given args. Returns the exit value or
// any error.
func (c *Cont) ExecArgs(args []string) (int, error) {
	return c.ExecWithSetup(&ExecSetup{Cmd: args})
}

// ExecWithSetup executes a command with the given setup.
func (c *Cont) ExecWithSetup(s *ExecSetup) (int, error) {
	return contExec(c.c, c.id, s)
}

// CopyInTar copies a tar stream into the container.
func (c *Cont) CopyInTar(r io.Reader, to string) error {
	return c.c.put(c.path("archive"), singleQuery("path", to), r)
}

func (c *Cont) getTar(p string) (*http.Response, error) {
	return c.c.get(c.path("archive"), singleQuery("path", p))
}

// CopyOutTar copies a file or a directory out of the container
// into a tarball stream.
func (c *Cont) CopyOutTar(fromPath string, w io.Writer) error {
	resp, err := c.getTar(fromPath)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if _, err := io.Copy(w, resp.Body); err != nil {
		return err
	}
	return resp.Body.Close()
}

// CopyOut copies a file or a directory into the destination.
func (c *Cont) CopyOut(src, destDir string) error {
	resp, err := c.getTar(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := writeTarToDir(resp.Body, destDir); err != nil {
		return err
	}
	return resp.Body.Close()
}

// CopyOutFile copies a single file into the destination.
func (c *Cont) CopyOutFile(src, dest string) error {
	resp, err := c.getTar(src)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if err := writeFirstFileAs(resp.Body, dest); err != nil {
		return err
	}
	return resp.Body.Close()
}

// Start starts the container.
func (c *Cont) Start() error { return c.c.poke(c.path("start"), nil) }

// SendSIGINT sends SIGINT (Ctrl-C) to the container.
func (c *Cont) SendSIGINT() error {
	return c.c.poke(c.path("kill"), singleQuery("signal", "SIGINT"))
}

// Stop stops the container.
func (c *Cont) Stop() error {
	err := c.c.poke(c.path("stop"), singleQuery("t", "60"))
	if httputil.ErrorStatusCode(err) == http.StatusNotModified {
		return nil // already stopped
	}
	return err
}

// Remove removes the container.
func (c *Cont) Remove() error { return c.c.del(c.path(""), nil) }

// Drop stops and removes the container.
func (c *Cont) Drop() error {
	if err := c.Stop(); err != nil {
		return err
	}
	return c.Remove()
}

// ForceRemove force removes the container.
func (c *Cont) ForceRemove() error {
	return c.c.del(c.path(""), singleQuery("force", "1"))
}

// Container events to wait for.
const (
	NotRunning = "not-running" // Wait till a container finishes running.
	NextExit   = "next-exit"   // Wait till the next exit event.
	Removed    = "removed"     // Wait till a container is removed.
)

// Wait for container to finish running.
func (c *Cont) Wait(cond string) (int, error) {
	q := singleQuery("condition", cond)
	var resp struct{ StatusCode int }
	if err := c.c.call(c.path("wait"), q, nil, &resp); err != nil {
		return 0, err
	}
	return resp.StatusCode, nil
}

// Inspect inspects the container.
func (c *Cont) Inspect() (*ContInfo, error) {
	info := new(ContInfo)
	if err := c.c.jsonGet(c.path("json"), nil, info); err != nil {
		return nil, err
	}
	return info, nil
}

// Exists checks if the container exists.
func (c *Cont) Exists() (bool, error) {
	if _, err := c.Inspect(); err != nil {
		if errcode.IsNotFound(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// FollowLogs follows the container's logs and forwards it into the writer.
func (c *Cont) FollowLogs(out io.Writer) error {
	q := make(url.Values)
	q.Add("follow", "true")
	q.Add("stdout", "true")
	q.Add("stderr", "true")
	sink := newLogSink(out, out)
	if _, err := c.c.getInto(c.path("logs"), q, sink); err != nil {
		return err
	}
	return sink.waitDone()
}

func (c *Cont) rename(to string) error {
	// after the renaming, the container will become not usable.
	return c.c.poke(c.path("rename"), singleQuery("name", to))
}

// RenameCont renames an existing container.
func RenameCont(client *Client, from, to string) error {
	return NewCont(client, from).rename(to)
}
