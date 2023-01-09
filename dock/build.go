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
	"encoding/json"
	"io"
	"net/url"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/tarutil"
)

// BuildConfig is the configuration for building an image.
type BuildConfig struct {
	Tarball  io.Reader
	Files    *tarutil.Stream
	Args     map[string]string
	UseCache bool
}

// BuildImage builds a docker image using the given tarball stream,
// and assigns the given tag.
func BuildImage(c *Client, tag string, tarball io.Reader) error {
	return BuildImageConfig(c, tag, &BuildConfig{Tarball: tarball})
}

// BuildImageStream builds a docker image using the given tarball stream,
// and assigns the given tag.
func BuildImageStream(c *Client, tag string, files *tarutil.Stream) error {
	return BuildImageConfig(c, tag, &BuildConfig{Files: files})
}

// BuildImageConfig builds the image with the given config.
func BuildImageConfig(c *Client, tag string, config *BuildConfig) error {
	q := make(url.Values)
	q.Add("t", tag)
	if !config.UseCache {
		q.Add("nocache", "true")
	}
	if len(config.Args) > 0 {
		argsBytes, err := json.Marshal(config.Args)
		if err != nil {
			return errcode.Annotate(err, "marshal args")
		}
		q.Add("buildargs", string(argsBytes))
	}

	r := config.Tarball
	var wr *writeToReader
	if r == nil && config.Files != nil {
		wr = newWriteToReader(config.Files)
		r = wr
		defer wr.Close()
	}
	if r == nil {
		return errcode.InvalidArgf("no input files")
	}

	sink := newStreamSink()
	if err := c.post("build", q, r, sink); err != nil {
		return err
	}
	if err := sink.waitDone(); err != nil {
		return err
	}
	if wr != nil {
		if err := wr.Join(); err != nil {
			return errcode.Annotate(err, "send in files")
		}
	}
	return nil
}

// NewTarStream creates a stream with a docker file.
func NewTarStream(dockerfile string) *tarutil.Stream {
	ts := tarutil.NewStream()
	AddDockerfileToStream(ts, dockerfile)
	return ts
}

// AddDockerfileToStream adds a DockerFile of content with mode 0600.
func AddDockerfileToStream(s *tarutil.Stream, content string) {
	s.AddString("Dockerfile", tarutil.ModeMeta(0600), content)
}
