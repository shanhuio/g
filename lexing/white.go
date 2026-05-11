package lexing

// WhiteFunc is a function type that checks if a rune is white space.
type WhiteFunc func(r rune) bool

// IsWhite is the default IsWhite function for a lexer. Returns true on spaces,
// \t and \r.  Returns false on \n.
func IsWhite(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r'
}

// IsWhiteOrEndl is another IsWhite function that also returns true for \n.
func IsWhiteOrEndl(r rune) bool {
	return IsWhite(r) || r == '\n'
}
