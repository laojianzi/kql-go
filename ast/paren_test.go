package ast_test

import (
	"testing"

	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/token"
	"github.com/stretchr/testify/assert"
)

func TestParenExpr(t *testing.T) {
	type args struct {
		L, R int
		Expr ast.Expr
	}

	cases := []struct {
		name       string
		args       args
		wantPos    int
		wantEnd    int
		wantString string
	}{
		{
			name: `(f1: "v1")`,
			args: args{
				R:    10,
				Expr: ast.NewBinaryExpr(1, "f1", token.TokenKindOperatorEql, ast.NewLiteral(5, 9, token.TokenKindString, "v1"), false),
			},
			wantEnd:    10,
			wantString: `(f1: "v1")`,
		},
		{
			name: `("v1" OR "v2")`,
			args: args{
				R: 14,
				Expr: ast.NewCombineExpr(
					ast.NewLiteral(1, 5, token.TokenKindString, "v1"),
					token.TokenKindKeywordOr,
					ast.NewLiteral(9, 13, token.TokenKindString, "v2"),
				),
			},
			wantEnd:    14,
			wantString: `("v1" OR "v2")`,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr := ast.NewParenExpr(c.args.L, c.args.R, c.args.Expr)
			assert.Equal(t, c.wantPos, expr.Pos())
			assert.Equal(t, c.wantEnd, expr.End())
			assert.Equal(t, c.wantString, expr.String())
		})
	}
}
