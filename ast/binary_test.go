package ast_test

import (
	"testing"

	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/token"
	"github.com/stretchr/testify/assert"
)

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
