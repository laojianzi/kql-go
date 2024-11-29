package parser

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/laojianzi/kql-go/token"
)

// defaultLexer is a lexer implementation
type defaultLexer struct {
	Value []rune
	Token Token

	pos           int
	lastTokenKind token.Kind
	dotIdent      bool
}

// newLexer creates a new lexer
func newLexer(input string) *defaultLexer {
	return &defaultLexer{Value: []rune(strings.TrimSpace(input))}
}

// nextToken returns the next token from the input stream
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
	if l.eof() {
		l.Token.Kind = token.TokenKindEof

		return nil
	}

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

// consumeToken consumes the next token from the input stream
func (l *defaultLexer) consumeToken() error {
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
	}

	// ident as keyword, field or value
	return l.consumeIdent()
}

// consumeFieldToken consumes a field token
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

// shouldBreak checks if token collection should stop
func (l *defaultLexer) shouldBreak(i int, isString, withEscape bool, endChar rune) bool {
	ch := l.peek(i)
	if isString && !withEscape && ch == endChar {
		return true
	}

	if !isString && !withEscape {
		if unicode.IsSpace(ch) || ch == ')' || ch == ':' {
			return true
		}
	}

	// not \:
	if !isString && withEscape && ch == ':' && (!l.peekOk(i-1) || l.peek(i-1) != '\\') {
		return true
	}

	return false
}

// collectNextToken collects the next complete token starting from the given position
func (l *defaultLexer) collectNextToken(start int) string {
	buf := &bytes.Buffer{}
	buf.WriteRune(l.peek(start))

	for j := start; l.peekOk(j + 1); j++ {
		currentRune, nextRune := l.peek(j), l.peek(j+1)
		if currentRune != '\\' && (unicode.IsSpace(nextRune) || nextRune == ')' || nextRune == ':') {
			break
		}

		buf.WriteRune(nextRune)
	}

	return buf.String()
}

// processNonEscaped checks if a non-escaped character should break token collection
func (l *defaultLexer) processNonEscaped(i int, kind token.Kind) bool {
	return !token.RequireEscape(string(l.peek(i)), kind)
}

// consumeEscapedToken consumes a token that may contain escape sequences
// Returns the number of characters consumed, the positions of escape characters, and any error
func (l *defaultLexer) consumeEscapedToken(kind token.Kind, endChar rune) (i int, indexes []int, err error) {
	escape, buf := false, &bytes.Buffer{}

	isString := kind == token.TokenKindString
	if isString {
		i = 1 // skip opening quote
	}

	for l.peekOk(i) && !l.eof() {
		if l.shouldBreak(i, isString, escape, endChar) {
			break
		}

		var result *CharProcResult
		if escape {
			result, err = l.handleEscaped(i, kind, buf, indexes)
		} else {
			result = l.handleNonEscaped(i, kind, buf, indexes)
		}

		if err != nil {
			return 0, nil, err
		}

		if result.Position == i && !result.IsEscaped {
			break
		}

		i, escape, indexes = result.Position, result.IsEscaped, result.EscapeIndexes
	}

	if escape {
		return 0, nil, errors.New("unexpected escapes")
	}

	if buf.Len() > 0 {
		l.Token.Value += buf.String()
	}

	return i, indexes, nil
}

// consumeIdent consumes an identifier token
func (l *defaultLexer) consumeIdent() error {
	i, escapeIndexes, err := l.consumeEscapedToken(token.TokenKindIdent, 0)
	if err != nil {
		return err
	}

	l.Token.Kind = token.TokenKindIdent
	l.Token.EscapeIndexes = escapeIndexes

	if !strings.Contains(l.slice(0, i), "\\") {
		if token.IsKeyword(l.Token.Value) {
			l.Token.Kind = token.ToKeyword(l.Token.Value)
		} else if token.IsOperator(l.Token.Value) {
			l.Token.Kind = token.ToOperator(l.Token.Value)
		}
	}

	l.skipN(i)

	return nil
}

// consumeString consumes a string token
func (l *defaultLexer) consumeString() error {
	i, escapeIndexes, err := l.consumeEscapedToken(token.TokenKindString, '"')
	if err != nil {
		return err
	}

	if !l.peekOk(i) {
		return errors.New("expected double quote closed")
	}

	l.Token.Kind = token.TokenKindString
	l.Token.EscapeIndexes = escapeIndexes

	l.skipN(i + 1)

	return nil
}

// consumeNumber consumes a number token
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
			if b == '*' {
				return l.consumeIdent()
			}

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

// consumeOperator consumes an operator token
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

// consumeParen consumes a parenthesis token
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

// skipSpaces skips whitespace characters
func (l *defaultLexer) skipSpaces() {
	for !l.eof() {
		r, size := utf8.DecodeRuneInString(string(l.Value[l.pos:]))
		if !unicode.IsSpace(r) {
			return
		}

		l.skipN(size)
	}
}

// skipN skips n characters
func (l *defaultLexer) skipN(n int) {
	l.pos += n
}

// peek returns the character at the given position
func (l *defaultLexer) peek(i int) rune {
	return l.Value[l.pos+i]
}

// peekOk checks if the character at the given position is valid
func (l *defaultLexer) peekOk(i int) bool {
	return l.pos+i < len(l.Value)
}

// slice returns a substring of the input string
func (l *defaultLexer) slice(start, end int) string {
	if len(l.Value) < l.pos+end {
		end = len(l.Value) - l.pos
	}

	return string(l.Value[l.pos+start : l.pos+end])
}

// eof checks if the end of the input stream has been reached
func (l *defaultLexer) eof() bool {
	return l.pos >= len(l.Value)
}
