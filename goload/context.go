package goload

import (
	"go/build"
)

// MakeContext makes a build context for the given GOPATH.
func MakeContext(gopath string) *build.Context {
	ctx := build.Default
	if gopath != "" {
		ctx.GOPATH = gopath
	}
	return &ctx
}
