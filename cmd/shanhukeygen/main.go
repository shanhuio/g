package main

import (
	"flag"
	"log"
	"os"

	"shanhu.io/g/osutil"
	"shanhu.io/g/rsautil"
	"shanhu.io/g/termutil"
	"shanhu.io/std/errcode"
)

type config struct {
	nbit         int
	noPassphrase bool
}

func keygen(output string, config *config) error {
	var passphrase []byte
	if !config.noPassphrase {
		pass, err := termutil.ReadPassword("Key passphrase: ")
		if err != nil {
			return err
		}
		passphrase = pass
	}

	pri, pub, err := rsautil.GenerateKey(passphrase, config.nbit)
	if err != nil {
		return errcode.Annotate(err, "generate key")
	}

	if output == "" {
		return errcode.InvalidArgf("empty key name")
	}

	pemPath := output + ".pem"
	if yes, err := osutil.Exist(pemPath); err != nil {
		return err
	} else if yes {
		return errcode.InvalidArgf("key file %q already exists", pemPath)
	}

	if err := os.WriteFile(pemPath, pri, 0600); err != nil {
		return errcode.Annotate(err, "write out key file")
	}

	return os.WriteFile(output+".pub", pub, 0600)
}

func main() {
	out := flag.String("out", "", "key path to output")
	nopass := flag.Bool("nopass", false, "no passphrase")
	nbit := flag.Int("nbit", 4096, "number of RSA bits")
	flag.Parse()

	conf := &config{
		nbit:         *nbit,
		noPassphrase: *nopass,
	}

	if err := keygen(*out, conf); err != nil {
		log.Fatal(err)
	}
}
