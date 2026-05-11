package lexing

// Logger is an error logging interface
type Logger interface {
	Errorf(p *Pos, fmt string, args ...interface{})
}

// LogError adds a error to the logger if the error is not nil and returns
// true.  If the error is nil, it returns false.
func LogError(log Logger, e error) bool {
	if e == nil {
		return false
	}

	log.Errorf(nil, "%s", e.Error())
	return true
}
