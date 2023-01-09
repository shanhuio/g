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

package caco3bin

import (
	"shanhu.io/pub/caco3"
	"shanhu.io/pub/flagutil"
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
