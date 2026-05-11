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
