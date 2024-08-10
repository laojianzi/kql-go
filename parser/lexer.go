package parser

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"
)

type defaultLexer struct {
	input   string
	current int
}

func (l *defaultLexer) peekToken() (*Token, error) {
	l.skipWhitespace()

	if l.isEof() {
		return &Token{Kind: TokenKindEof}, nil
	}

	switch l.peekN(0) {
	case ':', '<', '>': // operator
		return l.peekOperator()
	case '+', '-': // with sign int or float
		fallthrough // jump to number case
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // int or float
		return l.peekNumber()
	case '"': // double quote string
		return l.peekString()
	case '*': // e.g * or *xx* or *"xx"*
	}

	// ident as keyword, field or value
	return l.peekIdent()
}

func (l *defaultLexer) peekIdent() (*Token, error) {
	var i int
	for l.peekOk(i) && !unicode.IsSpace(rune(l.peekN(i))) && !l.isEof() {
		if l.peekN(i) == '\\' && l.peekOk(i+1) && isRequireEscape(l.peekN(i+1)) {
			i += 2

			continue
		}

		if isRequireEscape(l.peekN(i)) {
			break
		}

		i++
	}

	tok := &Token{
		Kind:  TokenKindIdent,
		Value: l.slice(0, i),
		Pos:   l.current,
		End:   l.current + i,
	}

	if IsKeyword(tok.Value) {
		tok.Kind = ToKeyword(tok.Value)
	}

	l.skipN(i)

	return tok, nil
}

func (l *defaultLexer) peekString() (*Token, error) {
	i, endChar := 1, byte('"')
	for l.peekOk(i) && l.peekN(i) != endChar {
		i++
	}

	if !l.peekOk(i) {
		return nil, errors.New("expected double quote closed")
	}

	defer l.skipN(i + 1) // cannot be called before pos/end are set

	return &Token{
		Kind:  TokenKindString,
		Value: l.slice(1, i),
		Pos:   l.current + 1,
		End:   l.current + i,
	}, nil
}

func (l *defaultLexer) peekNumber() (*Token, error) {
	var i int
	if l.peekN(0) == '+' || l.peekN(0) == '-' { // skip sign
		if !l.peekOk(i + 1) {
			return nil, fmt.Errorf("expected digit, but got Eof")
		}

		if nextChar := rune(l.peekN(i + 1)); !unicode.IsDigit(nextChar) {
			return nil, fmt.Errorf("expected digit, but got %q", string(nextChar))
		}

		i++
	}

	tok := &Token{Pos: l.current, Kind: TokenKindInt}

	for l.peekOk(i) {
		b := l.peekN(i)
		if unicode.IsSpace(rune(b)) {
			break
		}

		if !unicode.IsDigit(rune(b)) && b != '.' {
			return nil, fmt.Errorf("expected digit or decimal point, but got %q", string(b))
		}

		if b == '.' {
			if !l.peekOk(i + 1) {
				return nil, fmt.Errorf("expected digit, but got Eof")
			}

			if nextChar := rune(l.peekN(i + 1)); !unicode.IsDigit(nextChar) {
				return nil, fmt.Errorf("expected digit, but got %q", string(nextChar))
			}

			tok.Kind = TokenKindFloat
		}

		i++
	}

	tok.Value = l.slice(0, i)
	tok.End = l.current + i
	l.skipN(i)

	return tok, nil
}

func (l *defaultLexer) peekOperator() (*Token, error) {
	length, tok := 1, &Token{Pos: l.current}
	if (l.peekN(0) == '<' || l.peekN(0) == '>') && l.peekOk(1) && l.peekN(1) == '=' { // <= or >=
		length = 2
	}

	tok.End = l.current + length
	tok.Value = l.slice(0, length)
	tok.Kind = ToOperator(tok.Value)

	l.skipN(length)

	return tok, nil
}

func (l *defaultLexer) peekWhitespace() error {
	oldCurrent := l.current
	l.skipWhitespace()

	if l.isEof() || oldCurrent < l.current {
		return nil
	}

	return fmt.Errorf("expected whitespace or Eof, but got %q", string(l.input[l.current]))
}

func (l *defaultLexer) skipWhitespace() {
	for !l.isEof() {
		r, size := utf8.DecodeRuneInString(l.input[l.current:])
		if !unicode.IsSpace(r) {
			break
		}

		l.current += size
	}
}

func (l *defaultLexer) skipN(n int) {
	l.current += n
}

func (l *defaultLexer) slice(i, j int) string {
	return l.input[l.current+i : l.current+j]
}

func (l *defaultLexer) peekN(n int) byte {
	return l.input[l.current+n]
}

func (l *defaultLexer) peekOk(n int) bool {
	return l.current+n < len(l.input)
}

func (l *defaultLexer) isEof() bool {
	return l.current >= len(l.input)
}

func isRequireEscape(b byte) bool {
	return b == '"' || IsSpecialChar(string(b)) || b == '\\'
}
