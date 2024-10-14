package ast

import (
	"strings"

	"github.com/laojianzi/kql-go/token"
)

// Expr represents an expression in KQL.
type Expr interface {
	// Pos returns the position of the expression.
	Pos() int

	// End returns the end position of the expression.
	End() int

	// String returns the string representation of the expression.
	String() string
}

// ParenExpr is a parenthesis expression.
//
// Example:
//
//	`(f1: "v1" AND num: 1)`
type ParenExpr struct {
	L, R int // left and right position of the parenthesis
	Expr Expr
}

// NewParenExpr creates a new parenthesis expression.
func NewParenExpr(L, R int, expr Expr) *ParenExpr {
	return &ParenExpr{
		L:    L,
		R:    R,
		Expr: expr,
	}
}

// Pos returns the position of the parenthesis expression.
func (e *ParenExpr) Pos() int {
	return e.L
}

// End returns the end position of the parenthesis expression.
func (e *ParenExpr) End() int {
	return e.R
}

// String returns the string representation of the parenthesis expression.
func (e *ParenExpr) String() string {
	var buf strings.Builder

	buf.WriteByte('(')
	buf.WriteString(e.Expr.String())
	buf.WriteByte(')')

	return buf.String()
}

// CombineExpr is a combination expression.
//
// Example:
//
//	`f1: "v1" AND num: 1`
type CombineExpr struct {
	LeftExpr  Expr
	Keyword   token.Kind
	RightExpr Expr
}

// NewCombineExpr creates a new combination expression.
func NewCombineExpr(leftExpr Expr, keyword token.Kind, rightExpr Expr) *CombineExpr {
	return &CombineExpr{
		LeftExpr:  leftExpr,
		Keyword:   keyword,
		RightExpr: rightExpr,
	}
}

// Pos returns the position of the combination expression.
func (e *CombineExpr) Pos() int {
	return e.LeftExpr.Pos()
}

// End returns the end position of the combination expression.
func (e *CombineExpr) End() int {
	return e.RightExpr.End()
}

// String returns the string representation of the combination expression.
func (e *CombineExpr) String() string {
	var buf strings.Builder
	if e.LeftExpr != nil {
		buf.WriteString(e.LeftExpr.String())
	}

	if e.RightExpr != nil {
		buf.WriteByte(' ')
		buf.WriteString(e.Keyword.String())
		buf.WriteByte(' ')
		buf.WriteString(e.RightExpr.String())
	}

	return buf.String()
}

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
