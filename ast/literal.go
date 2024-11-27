package ast

import "github.com/laojianzi/kql-go/token"

// Literal is a literal(int, float, string or identifier) value.
type Literal struct {
	pos           int
	end           int
	escapeIndexes []int

	Kind            token.Kind // int, float, string or identifier
	Value           string
	WithDoubleQuote bool
}

// NewLiteral creates a new literal value.
func NewLiteral(pos, end int, kind token.Kind, value string, escapeIndexes []int) *Literal {
	return &Literal{
		pos:             pos,
		end:             end,
		Kind:            kind,
		Value:           value,
		WithDoubleQuote: kind == token.TokenKindString,
		escapeIndexes:   escapeIndexes,
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
	value := e.Value

	if len(e.escapeIndexes) > 0 {
		var (
			runes     = []rune(value)
			newValue  []rune
			lastIndex int
		)

		for _, escapeIndex := range e.escapeIndexes {
			newValue = append(newValue, runes[lastIndex:escapeIndex]...)
			newValue = append(newValue, '\\')
			lastIndex = escapeIndex
		}

		value = string(append(newValue, runes[lastIndex:]...))
	}

	if e.WithDoubleQuote {
		return `"` + value + `"`
	}

	return value
}
