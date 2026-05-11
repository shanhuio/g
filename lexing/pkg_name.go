package lexing

// IsPkgName checks if a literal is a valid package name.
func IsPkgName(s string) bool {
	// must not an empty string
	if s == "" {
		return false
	}

	for i, r := range s {
		if r >= '0' && r <= '9' {
			if i == 0 {
				return false
			}
			continue
		}
		if r >= 'a' && r <= 'z' {
			continue
		}
		return false
	}
	return true
}
