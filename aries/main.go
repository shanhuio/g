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

package aries

import (
	"flag"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"shanhu.io/pub/unixhttp"
)

// ListenAndServe serves on the address. If the address ends
// with .sock, it ListenAndServe's on the unix domain socket.
func ListenAndServe(addr string, s Service) error {
	log.Printf("serve on %q", addr)
	if strings.HasSuffix(addr, ".sock") {
		return unixhttp.ListenAndServe(addr, Serve(s))
	}
	return http.ListenAndServe(addr, Serve(s))
}

// DeclareAddrFlag declares the -addr flag.
func DeclareAddrFlag(def string) *string {
	if def == "" {
		if addr := os.Getenv("ADDR"); addr != "" {
			def = addr
		}
		if port := os.Getenv("PORT"); port != "" {
			def = ":" + port
		}
		if len(os.Args) > 0 {
			h := fnv.New32()
			io.WriteString(h, os.Args[0])
			const offset = 8000
			def = fmt.Sprintf("localhost:%d", offset+h.Sum32()%1000)
		}
	}

	return flag.String("addr", def, "address to listen on")
}
