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
	"os"
	"path/filepath"

	"shanhu.io/pub/caco3"
	"shanhu.io/pub/errcode"
	"shanhu.io/pub/lexing"
)

func cmdBuild(args []string) error {
	flags := cmdFlags.New()
	config := new(caco3.Config)
	declareBuildFlags(flags, config)
	args = flags.ParseArgs(args)

	wd, err := os.Getwd()
	if err != nil {
		return errcode.Annotate(err, "get work dir")
	}
	if config.Root != "" {
		root, err := filepath.Abs(config.Root)
		if err != nil {
			return errcode.Annotate(err, "get abs root dir")
		}
		config.Root = root
	}

	b, err := caco3.NewBuilder(wd, config)
	if err != nil {
		return errcode.Annotate(err, "new builder")
	}

	if _, errs := b.ReadWorkspace(); errs != nil {
		lexing.FprintErrs(os.Stderr, errs, wd)
		return errcode.InvalidArgf("read workspace got %d errors", len(errs))
	}

	if errs := b.Build(args); errs != nil {
		lexing.FprintErrs(os.Stderr, errs, wd)
		return errcode.InvalidArgf("build got %d errors", len(errs))
	}

	return nil
}
