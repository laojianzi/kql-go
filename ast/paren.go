package ast

import "strings"

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
