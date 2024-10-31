package ast

import "github.com/laojianzi/kql-go/token"

// Literal is a literal(int, float, string or identifier) value.
type Literal struct {
	pos             int
	end             int
	Kind            token.Kind // int, float, string or identifier
	Value           string
	WithDoubleQuote bool
}

// NewLiteral creates a new literal value.
func NewLiteral(pos, end int, kind token.Kind, value string) *Literal {
	return &Literal{
		pos:             pos,
		end:             end,
		Kind:            kind,
		Value:           value,
		WithDoubleQuote: kind == token.TokenKindString,
	}
}

// Pos returns the position of the literal value.
func (e *Literal) Pos() int {
	return e.pos
}

// End returns the end position of the literal value.
func (e *Literal) End() int {
	return e.end
}

// String returns the string representation of the literal value.
func (e *Literal) String() string {
	if e.WithDoubleQuote {
		return `"` + e.Value + `"`
	}

	return e.Value
}
