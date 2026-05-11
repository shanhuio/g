package goload

import (
	"go/build"
	"strings"

	"golang.org/x/tools/go/buildutil"
)

func skipPkg(p string) bool {
	if p == "C" {
		return true
	}
	if p == "lib9" {
		return true
	}
	if strings.HasPrefix(p, "lib9/") {
		return true
	}
	if p == "cmd" {
		return true
	}
	if strings.HasPrefix(p, "cmd/") {
		return true
	}
	return false
}

// ListStdPkgs returns a list of all packages in
// the golang standard library with the specific GOPATH.
func ListStdPkgs() []string {
	c := build.Default
	c.CgoEnabled = false
	c.GOPATH = ""

	pkgs := buildutil.AllPackages(&c)

	var ret []string
	for _, p := range pkgs {
		if skipPkg(p) {
			continue
		}

		if _, err := build.Import(p, "", 0); err != nil {
			continue // ignore error packges
		}

		ret = append(ret, p)
	}
	return ret
}
