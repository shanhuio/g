// Copyright (C) 2022  Shanhu Tech Inc.
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
	"go/build"
	"path/filepath"
	"sort"
	"strings"

	"shanhu.io/pub/goenv"
	"shanhu.io/pub/goload"
)

func absGOPATH() (string, error) {
	gopath, err := goenv.GOPATH()
	if err != nil {
		return "", err
	}
	abs, err := filepath.Abs(gopath)
	if err != nil {
		return "", err
	}
	return abs, nil
}

type relPkg struct {
	abs string
	rel string
	pkg *build.Package
}

func relPkgs(rootPkg string, scanRes *goload.ScanResult) ([]*relPkg, error) {
	var pkgs []string
	for pkg := range scanRes.Pkgs {
		pkgs = append(pkgs, pkg)
	}
	sort.Strings(pkgs)

	var ret []*relPkg
	prefix := rootPkg + "/"

	for _, pkg := range pkgs {
		rel := &relPkg{
			abs: pkg,
			pkg: scanRes.Pkgs[pkg].Build,
		}

		if pkg == rootPkg {
			rel.rel = "."
			ret = append(ret, rel)
			continue
		}

		if strings.HasPrefix(pkg, prefix) {
			rel.rel = "./" + strings.TrimPrefix(pkg, prefix)
			ret = append(ret, rel)
			continue
		}

		return nil, fmt.Errorf("%q is not in %q", pkg, rootPkg)
	}
	return ret, nil
}
