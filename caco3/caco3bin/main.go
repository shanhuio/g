package caco3bin

import (
	"shanhu.io/g/subcmd"
)

func cmd() *subcmd.List {
	c := subcmd.New()
	c.Add("build", "build rules", cmdBuild)
	c.Add("sync", "sync source repos", cmdSync)
	return c
}

// Main is the entrance for the caco3 binary.
func Main() { cmd().Main() }
