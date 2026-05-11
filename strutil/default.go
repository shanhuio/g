package strutil

// Default returns s if s is not empty, and returns def if is is.
func Default(s, def string) string {
	if s == "" {
		return def
	}
	return s
}
