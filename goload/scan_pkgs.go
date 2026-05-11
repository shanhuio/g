package goload

import (
	"path"
	"path/filepath"
	"sort"
)

// ScanModPkgs scans all packages in a module.
func ScanModPkgs(mod, dir string, opts *ScanOptions) (*ScanResult, error) {
	s := newScanner(mod, opts)
	s.gomod = true
	s.res = newScanResult(mod)
	d := &scanDir{dir: dir, path: mod}
	if err := s.walk(d); err != nil {
		return nil, err
	}
	return s.res, nil
}

// ScanPkgs scans all packages under a package path.
func ScanPkgs(p string, opts *ScanOptions) (*ScanResult, error) {
	s := newScanner(p, opts)

	// First check if the folder can be found.
	s.res = newScanResult(p)
	dir := &scanDir{
		dir:  filepath.Join(s.srcRoot(), filepath.ToSlash(p)),
		path: p,
		base: path.Base(p),
	}

	for _, scanning := range []bool{true, false} {
		s.vendorScanning = scanning
		if err := s.walk(dir); err != nil {
			return nil, err
		}
	}

	return s.res, nil
}

// ListPkgs list all packages under a package path.
func ListPkgs(p string) ([]string, error) {
	res, err := ScanPkgs(p, nil)
	if err != nil {
		return nil, err
	}

	var lst []string
	for pkg := range res.Pkgs {
		lst = append(lst, pkg)
	}
	sort.Strings(lst)
	return lst, nil
}
