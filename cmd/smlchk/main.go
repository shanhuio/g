package main

import (
	"flag"
	"fmt"
	"os"

	"shanhu.io/g/gocheck"
	"shanhu.io/g/goload"
	"shanhu.io/std/lexing"
)

func errExit(err error) {
	if err == nil {
		return
	}
	fmt.Fprintln(os.Stderr, err)
	os.Exit(-1)
}

func errsExit(errs []*lexing.Error) {
	if len(errs) == 0 {
		return
	}
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(-1)
}

func main() {
	path := flag.String("path", "shanhu.io/smlvm", "go repo path to check")
	textHeight := flag.Int("height", 300, "maximum height for a single file")
	textWidth := flag.Int("width", 80, "maximum width for a single file")
	verbose := flag.Bool("v", false, "prints package names")
	flag.Parse()

	pkgs, err := goload.ListPkgs(*path)
	errExit(err)

	for _, pkg := range pkgs {
		if *verbose {
			fmt.Println(pkg)
		}
		errs := gocheck.CheckAll(pkg, *textHeight, *textWidth)
		errsExit(errs)
	}
}
