package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"shanhu.io/g/dags"
	"shanhu.io/g/godep"
	"shanhu.io/g/goload"
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
