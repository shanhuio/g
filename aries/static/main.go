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

package static

import (
	"flag"
	"log"

	"shanhu.io/g/aries"
)

func makeService(root string) aries.Service {
	return aries.NewStaticFiles(root)
}

// Main is the main entrance for smlstatic binary
func Main() {
	root := flag.String("root", "lib/site", "static directory to serve")
	addr := aries.DeclareAddrFlag("localhost:8000")
	flag.Parse()

	s := makeService(*root)
	if err := aries.ListenAndServe(*addr, s); err != nil {
		log.Fatal(err)
	}
}
