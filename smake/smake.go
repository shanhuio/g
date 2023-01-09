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

package smake

import (
	"fmt"
	"path"
	"path/filepath"

	lintpkg "golang.org/x/lint"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/gocheck"
	"shanhu.io/pub/goload"
	"shanhu.io/pub/gomod"
	"shanhu.io/pub/gotags"
	"shanhu.io/pub/lexing"
)

func smlchkPkg(c *context, pkg *relPkg) []*lexing.Error {
	const textHeight = 320 // 20 lines for license notice.
	const textWidth = 80

	dir := filepath.Join(c.workDir(), filepath.FromSlash(pkg.rel))
	return gocheck.ModCheckAll(dir, pkg.abs, textHeight, textWidth)
}

func smlchk(c *context, pkgs []*relPkg) error {
	c.logln("smlchk")

	for _, pkg := range pkgs {
		if errs := smlchkPkg(c, pkg); len(errs) != 0 {
			for _, err := range errs {
				c.logln(err)
			}
			return fmt.Errorf("smlchk %q failed", pkg.rel)
		}
	}
	return nil
}

func lint(c *context, pkgs []*relPkg) error {
	c.logln("lint")

	const minConfidence = 0.8
	for _, pkg := range pkgs {
		files, err := fileSourceMap(pkg)
		if err != nil {
			return err
		}

		l := new(lintpkg.Linter)
		ps, err := l.LintFiles(files)
		if err != nil {
			return err
		}

		errCount := 0
		for _, p := range ps {
			if p.Confidence < minConfidence {
				continue
			}
			c.logf("%v: %s\n", p.Position, p.Text)
			errCount++
		}

		if errCount > 0 {
			return fmt.Errorf("lint %q failed", pkg.rel)
		}
	}
	return nil
}

func tags(c *context, pkgs []*relPkg) error {
	if !c.atModRoot() {
		return nil
	}
	c.logln("tags")

	var files []string
	for _, pkg := range pkgs {
		list := listAbsFiles(pkg.pkg)
		files = append(files, list...)
	}
	return gotags.Write(files, "tags")
}

func listPkgs(c *context) ([]*relPkg, error) {
	root := c.modRootDir()
	workDir := c.workDir()

	modFile := filepath.Join(root, "go.mod")
	mod, err := gomod.Parse(modFile)
	if err != nil {
		return nil, errcode.Annotate(err, "parse go.mod")
	}

	relPath, err := filepath.Rel(root, workDir)
	if err != nil {
		return nil, errcode.Annotate(err, "get relative path")
	}
	relPkg := filepath.ToSlash(relPath)

	workPkg := path.Join(mod.Name, relPkg)

	scanRes, err := goload.ScanModPkgs(workPkg, workDir, nil)
	if err != nil {
		return nil, errcode.Annotate(err, "scan packages")
	}
	return relPkgs(workPkg, scanRes)
}

func smake(c *context) error {
	pkgs, err := listPkgs(c)
	if err != nil {
		return errcode.Annotate(err, "list packages")
	}

	if len(pkgs) == 0 {
		c.logln("no packages found")
		return nil
	}

	installCmd := []string{"go", "install"}

	if err := c.execPkgs(pkgs, []string{
		"gofmt", "-s", "-w", "-l",
	}, nil); err != nil {
		return err
	}
	if err := c.execPkgs(pkgs, installCmd, nil); err != nil {
		return err
	}

	if err := smlchk(c, pkgs); err != nil {
		return err
	}
	if err := lint(c, pkgs); err != nil {
		return err
	}

	if err := c.execPkgs(pkgs, []string{
		"go", "vet",
	}, nil); err != nil {
		return err
	}

	return tags(c, pkgs)
}
