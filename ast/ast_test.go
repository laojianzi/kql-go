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

func TestBinaryExpr(t *testing.T) {
	type args struct {
		pos      int
		field    string
		operator token.Kind
		value    ast.Expr
		hasNot   bool
	}

	cases := []struct {
		name       string
		args       args
		wantPos    int
		wantEnd    int
		wantString string
	}{
		{
			name: `"v1"`,
			args: args{
				pos:    0,
				value:  ast.NewLiteral(0, 4, token.TokenKindString, "v1"),
				hasNot: false,
			},
			wantEnd:    4,
			wantString: `"v1"`,
		},
		{
			name: `NOT "v1"`,
			args: args{
				pos:    0,
				value:  ast.NewLiteral(4, 8, token.TokenKindString, "v1"),
				hasNot: true,
			},
			wantEnd:    8,
			wantString: `NOT "v1"`,
		},
		{
			name: `f1: "v1"`,
			args: args{
				field:    "f1",
				operator: token.TokenKindOperatorEql,
				value:    ast.NewLiteral(4, 8, token.TokenKindString, "v1"),
				hasNot:   false,
			},
			wantEnd:    8,
			wantString: `f1: "v1"`,
		},
		{
			name: `NOT f1: "v1"`,
			args: args{
				pos:      0,
				field:    "f1",
				operator: token.TokenKindOperatorEql,
				value:    ast.NewLiteral(8, 12, token.TokenKindString, "v1"),
				hasNot:   true,
			},
			wantEnd:    12,
			wantString: `NOT f1: "v1"`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr := ast.NewBinaryExpr(c.args.pos, c.args.field, c.args.operator, c.args.value, c.args.hasNot)
			assert.Equal(t, c.wantPos, expr.Pos())
			assert.Equal(t, c.wantEnd, expr.End())
			assert.Equal(t, c.wantString, expr.String())
		})
	}
}

func TestLiteral(t *testing.T) {
	type args struct {
		pos             int
		end             int
		kind            token.Kind
		value           string
		withDoubleQuote bool
	}

	cases := []struct {
		name       string
		args       args
		wantPos    int
		wantEnd    int
		wantString string
	}{
		{
			name: "int literal",
			args: args{
				end:   3,
				kind:  token.TokenKindInt,
				value: "101",
			},
			wantEnd:    3,
			wantString: `101`,
		},
		{
			name: "float literal",
			args: args{
				end:   3,
				kind:  token.TokenKindFloat,
				value: "10.1",
			},
			wantEnd:    3,
			wantString: `10.1`,
		},
		{
			name: "string literal",
			args: args{
				pos:             1,
				end:             3,
				kind:            token.TokenKindString,
				value:           "v1",
				withDoubleQuote: true,
			},
			wantPos:    1,
			wantEnd:    3,
			wantString: `"v1"`,
		},
		{
			name: `identifier literal`,
			args: args{
				pos:             0,
				end:             2,
				kind:            token.TokenKindIdent,
				value:           "v1",
				withDoubleQuote: true,
			},
			wantEnd:    2,
			wantString: `v1`,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			expr := ast.NewLiteral(c.args.pos, c.args.end, c.args.kind, c.args.value)
			assert.Equal(t, c.wantPos, expr.Pos())
			assert.Equal(t, c.wantEnd, expr.End())
			assert.Equal(t, c.wantString, expr.String())
		})
	}
}
