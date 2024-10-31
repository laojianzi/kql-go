package ast_test

import (
	"testing"

	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/token"
	"github.com/stretchr/testify/assert"
)

func TestCombineExpr(t *testing.T) {
	type args struct {
		leftExpr  ast.Expr
		keyword   token.Kind
		rightExpr ast.Expr
	}

	cases := []struct {
		name       string
		args       args
		wantPos    int
		wantEnd    int
		wantString string
	}{
		{
			name: `f1: "v1" OR NOT f1: "v2"`,
			args: args{
				leftExpr:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1"), false),
				keyword:   token.TokenKindKeywordOr,
				rightExpr: ast.NewBinaryExpr(12, "f1", token.TokenKindOperatorEql, ast.NewLiteral(20, 24, token.TokenKindString, "v2"), true),
			},
			wantEnd:    24,
			wantString: `f1: "v1" OR NOT f1: "v2"`,
		},
		{
			name: `NOT f1: ("v1" OR "v2") AND f3: "v3"`,
			args: args{
				leftExpr: ast.NewBinaryExpr(
					0,
					"f1",
					token.TokenKindOperatorEql,
					ast.NewParenExpr(
						8,
						22,
						ast.NewCombineExpr(
							ast.NewBinaryExpr(9, "", 0, ast.NewLiteral(9, 13, token.TokenKindString, "v1"), false),
							token.TokenKindKeywordOr,
							ast.NewBinaryExpr(17, "", 0, ast.NewLiteral(17, 21, token.TokenKindString, "v2"), false),
						),
					),
					true,
				),
				keyword:   token.TokenKindKeywordAnd,
				rightExpr: ast.NewBinaryExpr(27, "f3", token.TokenKindOperatorEql, ast.NewLiteral(31, 35, token.TokenKindString, "v3"), false),
			},
			wantEnd:    35,
			wantString: `NOT f1: ("v1" OR "v2") AND f3: "v3"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr := ast.NewCombineExpr(c.args.leftExpr, c.args.keyword, c.args.rightExpr)
			assert.Equal(t, c.wantPos, expr.Pos())
			assert.Equal(t, c.wantEnd, expr.End())
			assert.Equal(t, c.wantString, expr.String())
		})
	}
}
