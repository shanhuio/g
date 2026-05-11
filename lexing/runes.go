package lexing

// IsLetter returns true if the rune is in a-z or A-Z
func IsLetter(r rune) bool {
	if r >= 'a' && r <= 'z' {
		return true
	}
	if r >= 'A' && r <= 'Z' {
		return true
	}
	return false
}

// IsDigit returns true when the rune is in 0-9
func IsDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

// IsHexDigit returns true when the rune is in 0-9, a-f or A-F
func IsHexDigit(r rune) bool {
	if IsDigit(r) {
		return true
	}
	if r >= 'a' && r <= 'f' {
		return true
	}
	if r >= 'A' && r <= 'F' {
		return true
	}
	return false
}
