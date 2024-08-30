package parser

import "strings"

type Expr interface {
	Pos() int
	End() int
	String() string
}

type WrapExpr struct {
	pos    int
	Field  string
	Layers int
	Expr   Expr
}

func (e *WrapExpr) Pos() int {
	return e.pos
}

func (e *WrapExpr) End() int {
	return e.Expr.End() + e.Layers
}

func (e *WrapExpr) String() string {
	var buf strings.Builder
	if e.Field != "" {
		buf.WriteString(e.Field)
		buf.WriteString(": ")
	}

	buf.WriteString(strings.Repeat("(", e.Layers))
	buf.WriteString(e.Expr.String())
	buf.WriteString(strings.Repeat(")", e.Layers))

	return buf.String()
}

type CombineExpr struct {
	LeftExpr  Expr
	Keyword   Kind
	RightExpr Expr
}

func (e *CombineExpr) Pos() int {
	return e.LeftExpr.Pos()
}

func (e *CombineExpr) End() int {
	return e.RightExpr.End()
}

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

type MatchExpr struct {
	pos      int
	Field    string
	Operator Kind
	Value    *Literal
	HasNot   bool
}

func (e *MatchExpr) Pos() int {
	return e.pos
}

func (e *MatchExpr) End() int {
	if e.Value.WithDoubleQuote {
		return e.Value.End() + 1
	}

	return e.Value.End()
}

func (e *MatchExpr) String() string {
	var buf strings.Builder
	if e.HasNot {
		buf.WriteString("NOT ")
	}

	if e.Field != "" {
		buf.WriteString(e.Field)

		if e.Operator != TokenKindOperatorEql {
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

type Literal struct {
	pos             int
	end             int
	Kind            Kind // int or float or string or ident
	Value           string
	WithDoubleQuote bool
}

func (e *Literal) Pos() int {
	return e.pos
}

func (e *Literal) End() int {
	return e.end
}

func (e *Literal) String() string {
	if e.WithDoubleQuote {
		return `"` + e.Value + `"`
	}

	return e.Value
}
