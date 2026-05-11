package caco3

import (
	"os"
	"path"
	"path/filepath"
	"strings"

	"shanhu.io/g/dock"
	"shanhu.io/g/errcode"
)

type env struct {
	dock *dock.Client

	rootDir     string
	workDir     string
	workSrcPath string

	srcDir string
	outDir string

	workspace *Workspace // Lazily loaded.

	nodeType func(name string) string
	ruleType func(name string) string
}

func (e *env) prepareOut(ps ...string) (string, error) {
	p := e.out(ps...)
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", err
	}
	return p, nil
}

func dirFilePath(dir string, ps ...string) string {
	if len(ps) == 0 {
		return dir
	}
	p := path.Join(ps...)
	return filepath.Join(dir, filepath.FromSlash(p))
}

func (e *env) root(ps ...string) string {
	return dirFilePath(e.rootDir, ps...)
}

func (e *env) out(ps ...string) string {
	return dirFilePath(e.outDir, ps...)
}

func (e *env) src(ps ...string) string {
	return dirFilePath(e.srcDir, ps...)
}

func (e *env) nameToRepoTag(name string) (string, error) {
	parts := strings.Split(name, "/")
	if len(parts) == 0 {
		return "", errcode.InvalidArgf("empty name")
	}
	if len(parts) != 4 {
		return "", errcode.InvalidArgf("invalid name %q", name)
	}

	domain := parts[0]
	project := parts[1]
	dockers := parts[2]
	base := parts[3]

	if dockers != "dockers" && !strings.HasSuffix(dockers, "-dockers") {
		return "", errcode.InvalidArgf("not a docker image name: %q", name)
	}

	if domain == "shanhu.io" {
		domain = "cr.shanhu.io"
	}

	return path.Join(domain, project, base), nil
}
