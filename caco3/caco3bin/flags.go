package caco3bin

import (
	"shanhu.io/g/caco3"
	"shanhu.io/g/flagutil"
)

var cmdFlags = flagutil.NewFactory("caco3")

func declareBuildFlags(flags *flagutil.FlagSet, c *caco3.Config) {
	flags.StringVar(&c.Root, "root", "", "root directory")
	flags.BoolVar(&c.AlwaysRebuild, "rebuild", false, "always rebuild")
	flags.BoolVar(
		&c.UseDockerBuildCache, "docker_build_cache", true,
		"use docker build cache or not",
	)
}
