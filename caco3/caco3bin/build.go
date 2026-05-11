package caco3bin

import (
	"os"
	"path/filepath"

	"shanhu.io/g/caco3"
	"shanhu.io/g/errcode"
	"shanhu.io/g/lexing"
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
