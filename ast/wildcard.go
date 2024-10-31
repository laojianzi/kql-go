package ast

// WildcardExpr is a wildcard expression.
//
// Example:
//
//	`*`
//	`5*0`
//	`f*o`
//	`"*foo*"`
type WildcardExpr struct {
	*Literal // identifier or string

	Indexes []int // index of wildcard
}

// NewWildcard creates a new wildcard expression.
func NewWildcardExpr(lit *Literal, indexes []int) *WildcardExpr {
	return &WildcardExpr{
		Literal: lit,
		Indexes: indexes,
	}
}

// Pos returns the position of the wildcard expression.
func (e *WildcardExpr) Pos() int {
	return e.pos
}

// End returns the end position of the wildcard expression.
func (e *WildcardExpr) End() int {
	return e.end
}

// String returns the string representation of the wildcard expression.
func (e *WildcardExpr) String() string {
	return e.Literal.String()
}
