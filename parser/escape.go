package parser

import (
	"bytes"
	"errors"
	"fmt"

	"github.com/laojianzi/kql-go/token"
)

// CharProcResult represents the result of character processing in token lexing
type CharProcResult struct {
	// Position represents the next position to process
	Position int
	// IsEscaped indicates if the next character should be treated as escaped
	IsEscaped bool
	// EscapeIndexes stores the positions of escape characters
	EscapeIndexes []int
}

// NewCharProcResult creates a new CharProcResult with the given values
func NewCharProcResult(position int, isEscaped bool, escapeIndexes []int) *CharProcResult {
	return &CharProcResult{
		Position:      position,
		IsEscaped:     isEscaped,
		EscapeIndexes: escapeIndexes,
	}
}

// String returns a string representation of CharProcResult
func (r CharProcResult) String() string {
	return fmt.Sprintf("CharProcResult{pos:%d, escaped:%v, indexes:%v}",
		r.Position, r.IsEscaped, r.EscapeIndexes)
}

// handleNonEscaped processes a non-escaped character and updates the token state
func (l *defaultLexer) handleNonEscaped(pos int, k token.Kind, buf *bytes.Buffer, indexes []int) *CharProcResult {
	ch := l.peek(pos)
	if ch == '\\' {
		nextPos, newIndexes := l.handleBackslash(pos, buf, k == token.TokenKindString, indexes)

		return NewCharProcResult(nextPos, true, newIndexes)
	}

	if k != token.TokenKindString && !l.processNonEscaped(pos, k) {
		return NewCharProcResult(pos, false, indexes)
	}

	buf.WriteRune(ch)

	return NewCharProcResult(pos+1, false, indexes)
}

// handleEscaped processes an escaped character and updates the token state
func (l *defaultLexer) handleEscaped(pos int, k token.Kind, buf *bytes.Buffer, indexes []int) (*CharProcResult, error) {
	valid, err := l.handleEscapeSequence(pos, k)
	if err != nil {
		return NewCharProcResult(pos, true, indexes), err
	}

	if !valid {
		return NewCharProcResult(pos, false, indexes), nil
	}

	buf.WriteRune(l.peek(pos))

	return NewCharProcResult(pos+1, false, indexes), nil
}

// handleBackslash processes a backslash character and updates the token state
func (l *defaultLexer) handleBackslash(pos int, buf *bytes.Buffer, isString bool, indexes []int) (int, []int) {
	l.Token.Value += buf.String()
	buf.Reset()

	offset := 0
	if isString {
		offset = -1 // adjust for opening quote
	}

	return pos + 1, append(indexes, pos+offset-len(indexes))
}

// handleEscapeSequence validates and processes an escape sequence
func (l *defaultLexer) handleEscapeSequence(pos int, k token.Kind) (bool, error) {
	if !l.peekOk(pos) {
		return false, nil
	}

	ch := l.peek(pos)

	// Handle string literals specially
	if k == token.TokenKindString {
		switch ch {
		case '"', '\\', '*':
			return true, nil
		default:
			return false, nil
		}
	}

	// Handle special cases first
	if ch == '*' || token.RequireEscape(string(ch), k) {
		return true, nil
	}

	// Check if it's part of a keyword or operator
	nextToken := l.collectNextToken(pos)
	if token.IsKeyword(nextToken) || token.IsOperator(nextToken) {
		return true, nil
	}

	// If it's not a keyword or operator, check if it's a valid special character
	if !token.IsSpecialChar(string(ch)) {
		return false, errors.New("unexpected escapes")
	}

	return true, nil
}
