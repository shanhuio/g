package lexing

// Remover removes a particular type of token from a token stream
type Remover struct {
	Tokener
	t int
}

// NewRemover creates a new remover that removes token of type t
func NewRemover(t Tokener, typ int) *Remover {
	if typ == EOF {
		panic("cannot remove EOF")
	}

	ret := new(Remover)
	ret.Tokener = t
	ret.t = typ

	return ret
}

// Token implements the Tokener interface but only returns
// token that is not the particular type.
func (r *Remover) Token() *Token {
	for {
		ret := r.Tokener.Token()
		if ret.Type != r.t {
			return ret
		}
	}
}

// NewCommentRemover creates a new remover that removes token
func NewCommentRemover(t Tokener) *Remover {
	return NewRemover(t, Comment)
}
