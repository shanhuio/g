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

	"shanhu.io/g/caco3"
	"shanhu.io/g/errcode"
	"shanhu.io/g/lexing"
)

const sumsFile = "sums.jsonx"

func cmdSync(args []string) error {
	flags := cmdFlags.New()
	config := new(caco3.Config)
	declareBuildFlags(flags, config)
	pull := flags.Bool("pull", false, "pull latest commit")
	save := flags.Bool("save", false, "save latest commit into sums file")
	setRemotes := flags.Bool("set_remotes", false, "sets remote URLs")
	flags.ParseArgs(args)

	wd, err := os.Getwd()
	if err != nil {
		return errcode.Annotate(err, "get work dir")
	}

	b, err := caco3.NewBuilder(wd, config)
	if err != nil {
		return errcode.Annotate(err, "new builder")
	}

	if _, errs := b.ReadWorkspace(); errs != nil {
		lexing.FprintErrs(os.Stderr, errs, wd)
		return errcode.InvalidArgf("read workspace got %d errors", len(errs))
	}
	var sums *caco3.RepoSums
	if !*pull {
		s, err := caco3.ReadRepoSums(sumsFile)
		if err != nil {
			return errcode.Annotate(err, "read build sums")
		}
		sums = s
	}

	opts := &caco3.SyncOptions{
		SetRemotes: *setRemotes,
	}

	newSums, err := b.SyncRepos(sums, opts)
	if err != nil {
		return err
	}
	if *save {
		if err := caco3.SaveRepoSums(sumsFile, newSums); err != nil {
			return errcode.Annotate(err, "save build sums")
		}
	}
	return nil
}
