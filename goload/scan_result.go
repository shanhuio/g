package goload

import (
	"go/build"
)

// Package is a package in the scan result.
type Package struct {
	Build     *build.Package
	ImportMap map[string]string // import remapping due to vendoring

	ModRoot    string // module root import path, not including version
	ModVerRoot string // module root import path, including version
	ModVerPath string // alias import path when module is enabled
}

// ScanResult has the scanning result
type ScanResult struct {
	Repo        string
	Pkgs        map[string]*Package
	HasVendor   bool
	HasInternal bool
	Warnings    []error
}

func newScanResult(repo string) *ScanResult {
	return &ScanResult{
		Repo: repo,
		Pkgs: make(map[string]*Package),
	}
}
