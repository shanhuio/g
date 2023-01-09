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
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"

	"shanhu.io/pub/pathutil"
)

func getBitBucketRepoType(c *http.Client, pkgParts []string) (
	string, error,
) {
	if len(pkgParts) < 3 {
		return "", fmt.Errorf("invalid bitbucket repo")
	}

	p := fmt.Sprintf("%s/%s", pkgParts[1], pkgParts[2])
	url, err := url.Parse("https://api.bitbucket.org/2.0/repositories/" + p)
	if err != nil {
		return "", err
	}

	resp, err := c.Get(url.String())
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	t, err := bitBucketRepoType(resp.Body)
	if err != nil {
		return "", err
	}
	if err := resp.Body.Close(); err != nil {
		return "", err
	}
	return t, nil
}

func get(c *http.Client, pkg string) (*Repo, error) {
	url, err := url.Parse("https://" + pkg)
	if err != nil {
		return nil, err
	}
	url.RawQuery = "go-get=1"

	resp, err := c.Get(url.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	repo, err := ParseGoImport(resp.Body, pkg)
	if err != nil {
		return nil, err
	}
	if err := resp.Body.Close(); err != nil {
		return nil, err
	}

	return repo, nil
}

func getRepo(c *http.Client, pkg string) (*Repo, error) {
	ret, err := get(c, pkg)
	if err != nil {
		return nil, err
	}

	check, err := get(c, ret.ImportRoot)
	if err != nil {
		return nil, err
	}
	if check.ImportRoot != ret.ImportRoot {
		return nil, fmt.Errorf(
			"repo path mismatch for %q, sub has %q, parent %q",
			pkg, ret.ImportRoot, check.ImportRoot,
		)
	}
	if check.VCSRoot != ret.VCSRoot {
		return nil, fmt.Errorf(
			"vcs mismatch for %q, sub has %q, parent %q",
			pkg, ret.VCSRoot, check.VCSRoot,
		)
	}

	return check, nil
}

func commonRepo(pkg, vcs string, parts []string) (*Repo, error) {
	if len(parts) < 3 {
		return nil, fmt.Errorf("cannot find repo for pkg: %q", pkg)
	}

	repoPath := path.Join(parts[:3]...)

	return &Repo{
		ImportRoot: repoPath,
		VCS:        vcs,
		VCSRoot:    "https://" + repoPath,
	}, nil
}

// GetRepo gets the repo meta data for a particular package.
func GetRepo(c *http.Client, pkg string) (*Repo, error) {
	parts, err := pathutil.Split(pkg)
	if err != nil {
		return nil, err
	}

	domain := parts[0]
	switch domain {
	case "github.com":
		return commonRepo(pkg, "git", parts)
	case "bitbucket.org":
		repoType, err := getBitBucketRepoType(c, parts)
		if err != nil {
			return nil, fmt.Errorf("pkg %q: %s", pkg, err)
		}
		return commonRepo(pkg, repoType, parts)
	}

	for i, part := range parts {
		if strings.HasSuffix(part, ".git") {
			repoPath := path.Join(parts[:i+1]...)
			return &Repo{
				ImportRoot: repoPath,
				VCS:        "git",
				VCSRoot:    "https://" + repoPath,
			}, nil
		}
	}

	return get(c, pkg)
}
