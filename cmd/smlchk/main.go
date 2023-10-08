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
	"os"

	"shanhu.io/g/gocheck"
	"shanhu.io/g/goload"
	"shanhu.io/g/lexing"
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
