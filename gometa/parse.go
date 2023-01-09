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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"strings"

	"shanhu.io/pub/pathutil"
)

func charsetReader(charset string, r io.Reader) (io.Reader, error) {
	switch strings.ToLower(charset) {
	case "ascii":
		return r, nil
	}
	return nil, fmt.Errorf("charset %q not supported", charset)
}

func attrValue(attrs []xml.Attr, name string) string {
	for _, a := range attrs {
		if strings.EqualFold(a.Name.Local, name) {
			return a.Value
		}
	}
	return ""
}

func findMeta(r io.Reader, name string) ([]string, error) {
	var ret []string
	dec := xml.NewDecoder(r)
	dec.CharsetReader = charsetReader
	dec.Strict = false
	for {
		t, err := dec.RawToken()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if e, ok := t.(xml.StartElement); ok {
			if strings.EqualFold(e.Name.Local, "body") {
				break
			}
		}

		if e, ok := t.(xml.EndElement); ok {
			if strings.EqualFold(e.Name.Local, "head") {
				break
			}
		}

		e, ok := t.(xml.StartElement)
		if !ok || !strings.EqualFold(e.Name.Local, "meta") {
			continue
		}

		if attrValue(e.Attr, "name") != name {
			continue
		}

		ret = append(ret, attrValue(e.Attr, "content"))
	}
	return ret, nil
}

func bitBucketRepoType(r io.Reader) (string, error) {
	dec := json.NewDecoder(r)
	var dat struct {
		SCM string `json:"scm"`
	}
	if err := dec.Decode(&dat); err != nil {
		return "", err
	}
	return dat.SCM, nil
}

// ParseGoImport takes an HTML page and parses for the go-import meta tag.
func ParseGoImport(r io.Reader, pkg string) (*Repo, error) {
	metas, err := findMeta(r, "go-import")
	if err != nil {
		return nil, err
	}

	for _, meta := range metas {
		fields := strings.Fields(meta)
		if len(fields) != 3 {
			continue
		}

		repoPath := fields[0]
		if !pathutil.IsParent(repoPath, pkg) {
			continue
		}

		// we found it
		vcs := fields[1]
		url := fields[2]
		return &Repo{
			ImportRoot: repoPath,
			VCS:        vcs,
			VCSRoot:    url,
		}, nil
	}

	return nil, fmt.Errorf("go meta not found for %q", pkg)
}
