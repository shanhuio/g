package smake

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"shanhu.io/std/errcode"
)

func workDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(wd)
	if err != nil {
		return "", err
	}
	return filepath.EvalSymlinks(abs)
}

func usingGoMod() bool {
	v, ok := os.LookupEnv("GO111MODULE")
	if !ok {
		return true
	}
	return strings.ToLower(v) != "off"
}

func run(dir string) error {
	mod := usingGoMod()
	if !mod {
		return errcode.Internalf("must use go module")
	}

	if dir == "" {
		wd, err := workDir()
		if err != nil {
			return errcode.Annotate(err, "get work dir")
		}
		dir = wd
	} else {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return errcode.Annotate(err, "get absolute work dir")
		}
		dir = abs
	}

	modRoot, err := findGoModuleRoot(dir)
	if err != nil {
		return errcode.Annotate(err, "find module root")
	}

	// This is to make sure that we run under the absolute directory path.
	// Otherwise, some go tools will fail to recognize the directory structure.
	if err := os.Chdir(dir); err != nil {
		return err
	}

	gopath, err := absGOPATH()
	if err != nil {
		return err
	}

	c := newContext(gopath, modRoot, dir)
	return smake(c)
}

// Main is the entry point for smake.
func Main() {
	workDir := flag.String("dir", "", "work directory")
	flag.Parse()

	if err := run(*workDir); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
