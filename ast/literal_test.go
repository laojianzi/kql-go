package ast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/token"
)

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
			expr := ast.NewLiteral(c.args.pos, c.args.end, c.args.kind, c.args.value, nil)
			assert.Equal(t, c.wantPos, expr.Pos())
			assert.Equal(t, c.wantEnd, expr.End())
			assert.Equal(t, c.wantString, expr.String())
		})
	}
}
