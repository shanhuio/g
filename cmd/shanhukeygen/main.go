// Copyright (C) 2022  Shanhu Tech Inc.
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

package main

import (
	"flag"
	"log"
	"os"

	"shanhu.io/pub/errcode"
	"shanhu.io/pub/osutil"
	"shanhu.io/pub/rsautil"
	"shanhu.io/pub/termutil"
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
