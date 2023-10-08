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

package main

import (
	"encoding/json"
	"flag"
	"io"
	"log"
	"os"

	"shanhu.io/g/dags"
	"shanhu.io/g/errcode"
)

func readInput(in string) ([]byte, error) {
	if in == "" {
		return io.ReadAll(os.Stdin)
	}
	return os.ReadFile(in)
}

func writeOutput(out string, bs []byte) error {
	if out == "" {
		_, err := os.Stdout.Write(bs)
		return err
	}
	return os.WriteFile(out, bs, 0644)
}

func layout(in, out string) error {
	bs, err := readInput(in)
	if err != nil {
		return errcode.Annotate(err, "read input")
	}

	g := new(dags.Graph)
	g = g.Reverse()

	if err := json.Unmarshal(bs, &g.Nodes); err != nil {
		return errcode.Annotate(err, "parse graph")
	}

	_, v, err := dags.Layout(g)
	if err != nil {
		return errcode.Annotate(err, "layout graph")
	}

	m := dags.Output(v)
	outBytes, err := json.Marshal(m)
	if err != nil {
		return errcode.Annotate(err, "encode output")
	}

	if err := writeOutput(out, outBytes); err != nil {
		return errcode.Annotate(err, "write output")
	}

	return nil
}

func main() {
	in := flag.String("in", "", "input file")
	out := flag.String("out", "", "output file")
	flag.Parse()

	if err := layout(*in, *out); err != nil {
		log.Fatal(err)
	}
}
