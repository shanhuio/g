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

package lexing

// Parser provides the common parser functions for parsing.
// It does NOT provide a working parser for any grammar.
type Parser struct {
	tokener Tokener
	errs    *ErrorList
	types   *Types
	t       *Token
}

// NewParser creates a new parser around a tokener
func NewParser(t Tokener, types *Types) *Parser {
	ret := new(Parser)

	ret.tokener = t
	ret.errs = NewErrorList()
	ret.Next() // read in
	ret.types = types

	return ret
}

// Errorf adds a new parser error to the parser's error list at a particular
// position.
func (p *Parser) Errorf(pos *Pos, f string, args ...interface{}) {
	p.errs.CodeErrorf(pos, "", f, args...)
}

// CodeErrorf adds a new parser error with a error code
func (p *Parser) CodeErrorf(pos *Pos, c, f string, args ...interface{}) {
	p.errs.CodeErrorf(pos, c, f, args...)
}

// ErrorfHere adds a new parser error at the current token position
func (p *Parser) ErrorfHere(f string, args ...interface{}) {
	p.Errorf(p.t.Pos, f, args...)
}

// CodeErrorfHere adds a new parser error at the current token position
func (p *Parser) CodeErrorfHere(c, f string, args ...interface{}) {
	p.CodeErrorf(p.t.Pos, c, f, args...)
}

// See checks if the current token is of type t.
func (p *Parser) See(t int) bool { return p.t.Type == t }

// Accept shifts the tokener by one token and returns true
// if the current token is of type t. It otherwise returns false and nothing
// happens.
func (p *Parser) Accept(t int) bool {
	if p.See(t) {
		p.Next()
		return true
	}
	return false
}

// SeeLit checks if the current token is of type t and the lit is
// exactly lit.
func (p *Parser) SeeLit(t int, lit string) bool {
	return p.See(t) && p.t.Lit == lit
}

// Token returns the current token.
func (p *Parser) Token() *Token { return p.t }

// Next shifts the tokener by one token and returns the new current token.
func (p *Parser) Next() *Token {
	p.t = p.tokener.Token()
	return p.t
}

// Shift shifts the token by one token and returns the last current token.
func (p *Parser) Shift() *Token {
	ret := p.t
	p.Next()
	return ret
}

// InError checks if the parser is in error state.
// A parser can enter error state by adding a parser error with Errorf() or
// ErrorfAt().  A parser leaves error by calling BailOut().
func (p *Parser) InError() bool { return p.errs.InJail() }

// BailOut bails out the parser from an error state.
func (p *Parser) BailOut() { p.errs.BailOut() }

// Jail puts the parser in error state without adding an error message.
func (p *Parser) Jail() { p.errs.Jail() }

// TypeStr returns the name of a type used by the type register of this
// parser.
func (p *Parser) TypeStr(t int) string {
	return p.types.Name(t)
}

// ExpectLit checks if the current token is type t and has literal lit.  If it
// is, the token is accepted, the current token is shifted, and it returns the
// accepted token.  If it is not, the call reports an error, enters the parser
// into error state, and returns nil.  If the parser is already in error
// state, the call returns nil immediately, and nothing is checked.
func (p *Parser) ExpectLit(t int, lit string) *Token {
	if p.InError() {
		return nil
	}

	if p.SeeLit(t, lit) {
		return p.Shift()
	}

	p.CodeErrorfHere("lexing.unexpected", "expect %s %s, got %s",
		p.TypeStr(t), lit, p.TypeStr(p.t.Type),
	)
	return nil
}

// Expect checks if the current token is type t. If it is, the token is
// accepted, the current token is shifted, and it returns the accepted token.
// If it is not, the call reports an error, enters the parser into error
// state, and returns nil. If the parser is already in error state, the call
// returns nil immediately, and nothing is checked.
func (p *Parser) Expect(t int) *Token {
	if p.InError() {
		return nil
	}

	if p.See(t) {
		return p.Shift()
	}

	p.CodeErrorfHere("lexing.unexpected", "expect %s, got %s",
		p.TypeStr(t), p.TypeStr(p.t.Type),
	)
	return nil
}

// SkipErrStmt skips tokens until it meets a token of type sep or the end of
// file (token EOF) and returns true, but only when the parser is in error
// state.  If the parser is not in error state, it returns false and nothing
// is skipped.
func (p *Parser) SkipErrStmt(sep int) bool {
	if !p.InError() {
		return false
	}

	for !p.See(sep) || p.See(EOF) {
		p.Next()
	}
	if p.See(sep) {
		p.Next()
	}

	p.BailOut()
	return true
}

// Errs returns the parsing error list if the lexing has no error.
// If lexing has error, it returns the lexing error list instead.
func (p *Parser) Errs() []*Error {
	ret := p.tokener.Errs()
	if ret != nil {
		return ret
	}
	return p.errs.Errs()
}
