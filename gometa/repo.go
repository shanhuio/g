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

package gometa

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"shanhu.io/g/aries"
	"shanhu.io/g/errcode"
)

// Repo is a Golang repository that this handler will handle.
type Repo struct {
	ImportRoot string
	VCS        string
	VCSRoot    string
}

func host(path string) string {
	host, _, found := strings.Cut(path, "/")
	if found {
		return host
	}
	return path
}

// NewGitRepo creates a new git repository for import redirection.
func NewGitRepo(path, repoAddr string) *Repo {
	return &Repo{
		ImportRoot: path,
		VCS:        "git",
		VCSRoot:    repoAddr,
	}
}

// Meta returns the HTML meta line that needs to be included in the
// header of the page.
func (r *Repo) Meta() string {
	return fmt.Sprintf(
		`<meta name="go-import" content="%s %s %s">`,
		r.ImportRoot, r.VCS, r.VCSRoot,
	)
}

// MetaContent returns the go-import meta content of the meta line.
func (r *Repo) MetaContent() string {
	return fmt.Sprintf("%s %s %s", r.ImportRoot, r.VCS, r.VCSRoot)
}

func (r *Repo) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := aries.NewContext(w, req)
	c.ErrCode(r.Serve(c))
}

// Serve serves the incomiing webapp request.
func (r *Repo) Serve(c *aries.C) error {
	path := strings.TrimSuffix(host(r.ImportRoot)+c.Req.URL.Path, "/")

	if !strings.HasPrefix(path, r.ImportRoot) {
		return errcode.NotFoundf("repo not found: %s", path)
	}

	d := &data{
		ImportRoot: r.ImportRoot,
		VCS:        r.VCS,
		VCSRoot:    r.VCSRoot,
		Suffix:     strings.TrimSuffix(path, r.ImportRoot),
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, d); err != nil {
		return err
	}
	c.Resp.Write(buf.Bytes())
	return nil
}
