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

package states

import (
	"fmt"
	"net/url"
	"strings"

	"shanhu.io/pub/s3util"
)

func parseS3Endpoint(host string) (bucket, endpoint string) {
	host = strings.TrimSuffix(host, ".") // Trims ending dot in domain name.
	firstDot := strings.Index(host, ".")
	if firstDot < 0 {
		return "", host
	}
	return host[:firstDot], host[firstDot+1:]
}

// Dial connects to a States storage using the given URL address.
func Dial(addr *url.URL, creds interface{}) (States, error) {
	switch addr.Scheme {
	case "file", "":
		return newDirBack(addr.Path), nil
	case "s3":
		bucket, ep := parseS3Endpoint(addr.Host)
		config := &s3util.Config{
			Endpoint: ep,
			Bucket:   bucket,
			BasePath: addr.Path,
		}
		s3Creds, ok := creds.(*s3util.Credential)
		if !ok {
			return nil, fmt.Errorf("credential not for s3")
		}
		s3, err := s3util.NewClient(config, s3Creds)
		if err != nil {
			return nil, err
		}
		return newS3Back(s3), nil
	case "mem":
		return newMemBack(), nil
	}
	return nil, fmt.Errorf("unknown scheme: %q", addr.Scheme)
}

// DialSpec connects to a States storage using the given URL string.
func DialSpec(s string, creds interface{}) (States, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parse url %q: %s", s, err)
	}
	return Dial(u, creds)
}
