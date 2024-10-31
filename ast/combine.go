package ast

import (
	"strings"

	"github.com/laojianzi/kql-go/token"
)

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
