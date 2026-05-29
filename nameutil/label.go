package nameutil

import (
	"shanhu.io/std/errcode"
)

// CheckLabel checks whether a string can be safely used as a
// sub-domain name.
func CheckLabel(s string) error {
	if len(s) == 0 {
		return errcode.InvalidArgf("empty name")
	}
	if len(s) > 50 {
		return errcode.InvalidArgf("name too long: %q", s)
	}

	if s[0] == '-' {
		return errcode.InvalidArgf("%q starts with hypen", s)
	}
	if s[len(s)-1] == '-' {
		return errcode.InvalidArgf("%q ends with hypen", s)
	}
	lastHyphen := false
	for _, r := range s {
		if r == '-' {
			if lastHyphen {
				return errcode.InvalidArgf("%q has continous hyphen", s)
			}
			lastHyphen = true
			continue
		} else {
			lastHyphen = false
		}
		if r >= '0' && r <= '9' {
			continue
		}
		if r >= 'a' && r <= 'z' {
			continue
		}
		return errcode.InvalidArgf("%q contain invalid char: %q", s, r)
	}
	return nil
}
