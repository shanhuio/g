package lexing

// Recorder is a token filter that records all the token
// a tokener generates
type Recorder struct {
	Tokener
	tokens []*Token

	closed bool
}

// NewRecorder creates a new recorder that filters the tokener
func NewRecorder(t Tokener) *Recorder {
	ret := new(Recorder)
	ret.Tokener = t
	return ret
}

// Token implements the Tokener interface by
// relaying the call to the internal Tokener.
func (r *Recorder) Token() *Token {
	ret := r.Tokener.Token()
	r.tokens = append(r.tokens, ret)
	return ret
}

// Tokens returns the slice of recorded tokens.
func (r *Recorder) Tokens() []*Token { return r.tokens }

// TokenAll returns all the tokens fetched out of a tokener.
func TokenAll(t Tokener) []*Token {
	rec := NewRecorder(t)
	for {
		tok := rec.Token()
		if tok.Type == EOF {
			break
		}
	}
	return rec.Tokens()
}
