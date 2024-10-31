package ast

// Expr represents an expression in KQL.
type Expr interface {
	// Pos returns the position of the expression.
	Pos() int

	// End returns the end position of the expression.
	End() int

	// String returns the string representation of the expression.
	String() string
}
