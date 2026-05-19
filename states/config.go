package states

import (
	"fmt"
	"net/url"
	"strings"

	"shanhu.io/g/s3util"
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
func Dial(addr *url.URL, creds any) (States, error) {
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
func DialSpec(s string, creds any) (States, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, fmt.Errorf("parse url %q: %s", s, err)
	}
	return Dial(u, creds)
}
