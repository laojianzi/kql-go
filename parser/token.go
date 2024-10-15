package parser

import "github.com/laojianzi/kql-go/token"

// Token is a token parsed from lexer.
type Token struct {
	Pos   int
	End   int
	Kind  token.Kind
	Value string
}

// Clone returns a copy of the token.
func (t *Token) Clone() *Token {
	tok := *t

	return &tok
}
