package parser

import (
	"errors"
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/laojianzi/kql-go/token"
)

type defaultLexer struct {
	Value []rune
	Token Token

	pos           int
	lastTokenKind token.Kind
	dotIdent      bool
}

func (l *defaultLexer) nextToken() error {
	l.lastTokenKind = l.Token.Kind

	for {
		i := l.pos
		l.skipSpaces()

		if l.pos == i {
			break
		}
	}

	l.Token = Token{Pos: l.pos}
	defer func() {
		l.Token.End = l.pos
		if l.Token.Kind == token.TokenKindString { // skip the double quote "
			l.Token.Pos += 1
			l.Token.End -= 1
		}
	}()

	if !l.dotIdent {
		return l.consumeToken()
	}

	l.dotIdent = false

	return l.consumeFieldToken()
}

func (l *defaultLexer) consumeToken() error {
	if l.eof() {
		l.Token.Kind = token.TokenKindEof

		return nil
	}

	switch l.peek(0) {
	case ':', '<', '>': // operator
		return l.consumeOperator()
	case '(', ')':
		return l.consumeParen()
	case '+', '-': // with sign int or float
		fallthrough // jump to number case
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9': // int or float
		return l.consumeNumber()
	case '"': // double quote string
		return l.consumeString()
	case '*': // e.g * or *xx* or *"xx"*
	}

	// ident as keyword, field or value
	return l.consumeIdent()
}

func (l *defaultLexer) consumeFieldToken() error {
	if l.peekOk(0) && !unicode.IsSpace(rune(l.peek(0))) && !l.eof() {
		i := 0
		for l.peekOk(i) && !unicode.IsSpace(rune(l.peek(i))) && !l.eof() {
			i++
		}

		l.Token.Kind = token.TokenKindIdent
		l.Token.Value = string(l.Value[l.pos : l.pos+i])
		l.skipN(i)

		return nil
	}

	return l.consumeToken()
}

func (l *defaultLexer) consumeIdent() error {
	var i int
	for l.peekOk(i) && !unicode.IsSpace(rune(l.peek(i))) && !l.eof() {
		if l.peek(i) == '\\' && l.peekOk(i+1) && requireEscape(l.peek(i+1)) {
			i += 2

			continue
		}

		if requireEscape(l.peek(i)) {
			break
		}

		i++
	}

	l.Token.Kind = token.TokenKindIdent
	l.Token.Value = string(l.Value[l.pos : l.pos+i])

	if token.IsKeyword(l.Token.Value) {
		l.Token.Kind = token.ToKeyword(l.Token.Value)
	}

	l.skipN(i)

	return nil
}

func (l *defaultLexer) consumeString() error {
	i, endChar := 1, rune('"')
	for l.peekOk(i) && l.peek(i) != endChar {
		i++
	}

	if !l.peekOk(i) {
		return errors.New("expected double quote closed")
	}

	l.Token.Kind = token.TokenKindString
	l.Token.Value = l.slice(1, i)

	l.skipN(i + 1)

	return nil
}

func (l *defaultLexer) consumeNumber() error {
	var i int
	if l.peek(0) == '+' || l.peek(0) == '-' { // skip sign
		if !l.peekOk(i + 1) {
			return fmt.Errorf("expected digit, but got Eof")
		}

		if nextChar := l.peek(i + 1); !unicode.IsDigit(nextChar) {
			return fmt.Errorf("expected digit, but got %q", string(nextChar))
		}

		i++
	}

	l.Token.Kind = token.TokenKindInt

	for l.peekOk(i) {
		b := l.peek(i)
		if unicode.IsSpace(rune(b)) || b == ')' {
			break
		}

		if !unicode.IsDigit(rune(b)) && b != '.' {
			return fmt.Errorf("expected digit or decimal point, but got %q", string(b))
		}

		if b == '.' {
			if !l.peekOk(i + 1) {
				return fmt.Errorf("expected digit, but got Eof")
			}

			if nextChar := l.peek(i + 1); !unicode.IsDigit(nextChar) {
				return fmt.Errorf("expected digit, but got %q", string(nextChar))
			}

			l.Token.Kind = token.TokenKindFloat
		}

		i++
	}

	l.Token.Value = l.slice(0, i)

	l.skipN(i)

	return nil
}

func (l *defaultLexer) consumeOperator() error {
	length := 1
	if (l.peek(0) == '<' || l.peek(0) == '>') && l.peekOk(1) && l.peek(1) == '=' { // <= or >=
		length = 2
	}

	l.Token.Value = l.slice(0, length)
	l.Token.Kind = token.ToOperator(l.Token.Value)

	l.skipN(length)

	return nil
}

func (l *defaultLexer) consumeParen() error {
	l.Token.Value = l.slice(0, 1)

	switch l.peek(0) {
	case '(':
		l.Token.Kind = token.TokenKindLparen
	case ')':
		l.Token.Kind = token.TokenKindRparen
	default:
		return fmt.Errorf("expected token \"(\" or \")\", but got %q", string(l.peek(0)))
	}

	l.skipN(1)

	return nil
}

func (l *defaultLexer) skipSpaces() {
	for !l.eof() {
		r, size := utf8.DecodeRuneInString(string(l.Value[l.pos:]))
		if !unicode.IsSpace(r) {
			return
		}

		l.skipN(size)
	}
}

func (l *defaultLexer) skip() rune {
	r := l.Value[l.pos]
	l.pos++

	return r
}

func (l *defaultLexer) skipN(n int) {
	l.pos += n
}

func (l *defaultLexer) peek(i int) rune {
	return l.Value[l.pos+i]
}

func (l *defaultLexer) peekOk(i int) bool {
	return l.pos+i < len(l.Value)
}

func (l *defaultLexer) slice(start, end int) string {
	if len(l.Value) < l.pos+end {
		end = len(l.Value) - l.pos
	}

	return string(l.Value[l.pos+start : l.pos+end])
}

func (l *defaultLexer) eof() bool {
	return l.pos >= len(l.Value)
}

func requireEscape(r rune) bool {
	return r == '"' || token.IsSpecialChar(string(r)) || r == '\\'
}
