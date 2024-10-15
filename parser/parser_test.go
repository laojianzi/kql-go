package parser_test

import (
	"testing"

	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/parser"
	"github.com/laojianzi/kql-go/token"
	"github.com/stretchr/testify/assert"
)

func Test_defaultParser_Stmt(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cases := []struct {
			input string
			want  ast.Expr
		}{
			{
				input: "foo",
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindIdent, "foo"), false),
			},
			{
				input: "1",
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 1, token.TokenKindInt, "1"), false),
			},
			{
				input: "0.1",
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindFloat, "0.1"), false),
			},
			{
				input: `"0.1"`,
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 5, token.TokenKindString, "0.1"), false),
			},
			{
				input: `f1: "v1"`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1"), false),
			},
			{
				input: `f1 > 1`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorGtr, ast.NewLiteral(5, 6, token.TokenKindInt, "1"), false),
			},
			{
				input: `f1 < 1.1`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorLss, ast.NewLiteral(5, 8, token.TokenKindFloat, "1.1"), false),
			},
			{
				input: `f1 >= 1000.0001`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorGeq, ast.NewLiteral(6, 15, token.TokenKindFloat, "1000.0001"), false),
			},
			{
				input: `f1 <= 100000011`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorLeq, ast.NewLiteral(6, 15, token.TokenKindInt, "100000011"), false),
			},
		}

		for _, c := range cases {
			t.Run(c.input, func(t *testing.T) {
				stmt, err := parser.New(c.input).Stmt()
				assert.NoError(t, err)
				assert.EqualValues(t, c.want, stmt)
				assert.Equal(t, c.input, stmt.String())
			})
		}
	})

	t.Run("with keyword", func(t *testing.T) {
		cases := []struct {
			input string
			want  ast.Expr
		}{
			{
				input: "NOT bar",
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(4, 7, token.TokenKindIdent, "bar"), true),
			},
			{
				input: "foo AND bar",
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindIdent, "foo"), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(8, "", 0, ast.NewLiteral(8, 11, token.TokenKindIdent, "bar"), false),
				),
			},
			{
				input: "foo AND NOT bar",
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindIdent, "foo"), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(8, "", 0, ast.NewLiteral(12, 15, token.TokenKindIdent, "bar"), true),
				),
			},
			{
				input: `v1 AND 2 OR 0.3 AND NOT "v4" OR NOT 5.0`,
				want: ast.NewCombineExpr(
					ast.NewCombineExpr(
						ast.NewCombineExpr(
							ast.NewCombineExpr(
								ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 2, token.TokenKindIdent, "v1"), false),
								token.TokenKindKeywordAnd,
								ast.NewBinaryExpr(7, "", 0, ast.NewLiteral(7, 8, token.TokenKindInt, "2"), false),
							),
							token.TokenKindKeywordOr,
							ast.NewBinaryExpr(12, "", 0, ast.NewLiteral(12, 15, token.TokenKindFloat, "0.3"), false),
						),
						token.TokenKindKeywordAnd,
						ast.NewBinaryExpr(20, "", 0, ast.NewLiteral(24, 28, token.TokenKindString, "v4"), true),
					),
					token.TokenKindKeywordOr,
					ast.NewBinaryExpr(32, "", 0, ast.NewLiteral(36, 39, token.TokenKindFloat, "5.0"), true),
				),
			},
			{
				input: `f1: "v1" AND f2 > 2`,
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1"), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(13, "f2", token.TokenKindOperatorGtr, ast.NewLiteral(18, 19, token.TokenKindInt, "2"), false),
				),
			},
			// {
			// 	input: `f1: "v1" AND NOT f2 > 2`,
			// 	want: ast.NewCombineExpr(
			// 		ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1"), false),
			// 		token.TokenKindKeywordAnd,
			// 		ast.NewBinaryExpr(13, "f2", token.TokenKindOperatorGtr, ast.NewLiteral(18, 19, token.TokenKindInt, "2"), true),
			// 	),
			// },
			// {
			// 	input: `f1: "v1" AND f2 > 2 OR f3 < 0.3 AND NOT f4 >= 4 OR NOT f5 <= 5.0`,
			// 	want: ast.NewCombineExpr(
			// 		ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1"), false),
			// 		token.TokenKindKeywordAnd,
			// 		ast.NewBinaryExpr(13, "f2", token.TokenKindOperatorGtr, ast.NewLiteral(18, 19, token.TokenKindInt, "2"), false),
			// 	),
			// },
		}

		for _, c := range cases {
			t.Run(c.input, func(t *testing.T) {
				stmt, err := parser.New(c.input).Stmt()
				assert.NoError(t, err)
				assert.EqualValues(t, c.want, stmt)
				assert.Equal(t, c.input, stmt.String())
			})
		}
	})

	t.Run("with paren", func(t *testing.T) {
		cases := []struct {
			input string
			want  ast.Expr
		}{
			{
				input: "foo AND (NOT bar)",
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindIdent, "foo"), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(8, "", 0, ast.NewParenExpr(8, 17, ast.NewBinaryExpr(9, "", 0, ast.NewLiteral(13, 16, token.TokenKindIdent, "bar"), true)), false),
				),
			},
		}

		for _, c := range cases {
			t.Run(c.input, func(t *testing.T) {
				stmt, err := parser.New(c.input).Stmt()
				assert.NoError(t, err)
				assert.EqualValues(t, c.want, stmt)
				assert.Equal(t, c.input, stmt.String())
			})
		}
	})
}
