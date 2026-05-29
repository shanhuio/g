package smake

import (
	"fmt"
	"path"
	"path/filepath"

	lintpkg "golang.org/x/lint"
	"shanhu.io/g/errcode"
	"shanhu.io/g/gocheck"
	"shanhu.io/g/goload"
	"shanhu.io/g/gomod"
	"shanhu.io/g/gotags"
	"shanhu.io/std/lexing"
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

	installCmd := []string{"go", "install", "-buildvcs=false", "-trimpath"}

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
