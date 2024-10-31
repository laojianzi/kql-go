package ast

import (
	"strings"

	"github.com/laojianzi/kql-go/token"
)

// BinaryExpr is a binary expression.
//
// Example:
//
//	`NOT f1: "v1"`
type BinaryExpr struct {
	pos      int
	Field    string
	Operator token.Kind
	Value    Expr
	HasNot   bool
}

// NewBinaryExpr creates a new binary expression.
func NewBinaryExpr(pos int, field string, operator token.Kind, value Expr, hasNot bool) *BinaryExpr {
	return &BinaryExpr{
		pos:      pos,
		Field:    field,
		Operator: operator,
		Value:    value,
		HasNot:   hasNot,
	}
}

// Pos returns the position of the binary expression.
func (e *BinaryExpr) Pos() int {
	return e.pos
}

// End returns the end position of the binary expression.
func (e *BinaryExpr) End() int {
	if e.Value.End() < len(e.Value.String()) { // e.g. string values with double quotes
		return e.Value.End() + 1
	}

	return e.Value.End()
}

// String returns the string representation of the binary expression.
func (e *BinaryExpr) String() string {
	var buf strings.Builder
	if e.HasNot {
		buf.WriteString("NOT ")
	}

	if e.Field != "" {
		buf.WriteString(e.Field)

		if e.Operator != token.TokenKindOperatorEql {
			buf.WriteByte(' ')
		}

		buf.WriteString(e.Operator.String())
		buf.WriteByte(' ')
	}

	if e.Value != nil {
		buf.WriteString(e.Value.String())
	}

	return buf.String()
}
