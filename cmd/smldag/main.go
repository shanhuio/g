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

package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"shanhu.io/pub/dags"
	"shanhu.io/pub/godep"
	"shanhu.io/pub/goload"
)

func exitIf(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
}

func saveLayoutBytes(bs []byte, f string) {
	if strings.HasSuffix(f, ".js") {
		out, err := os.Create(f)
		exitIf(err)
		defer out.Close()

		_, err = io.WriteString(out, "var dag = ")
		exitIf(err)

		_, err = out.Write(bs)
		exitIf(err)

		_, err = io.WriteString(out, ";")
		exitIf(err)

		exitIf(out.Close())
		return
	}

	exitIf(os.WriteFile(f, bs, 0644))
}

func saveLayout(g *dags.Graph, f string) {
	m, err := dags.LayoutJSON(g)
	exitIf(err)
	saveLayoutBytes(m, f)
}

func repoDep(repo string) (*dags.Graph, error) {
	if repo == "" {
		return godep.StdDep()
	}

	pkgs, err := goload.ListPkgs(repo)
	if err != nil {
		return nil, err
	}
	g, err := godep.PkgDep(pkgs)
	if err != nil {
		return nil, err
	}

	repoSlash := repo + "/"
	return g.Rename(func(name string) (string, error) {
		if name == repo {
			return "~", nil
		}
		return strings.TrimPrefix(name, repoSlash), nil
	})
}

func main() {
	repo := flag.String("repo", "", "repository to generate the dependency map")
	out := flag.String("out", "godag.json", "output JSON file")
	flag.Parse()

	g, e := repoDep(*repo)
	exitIf(e)

	saveLayout(g, *out)
}
