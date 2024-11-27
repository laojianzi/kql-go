package ast_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/token"
)

func TestWildcard(t *testing.T) {
	type args struct {
		pos             int
		end             int
		kind            token.Kind
		value           string
		withDoubleQuote bool
		indexes         []int
	}

	cases := []struct {
		name        string
		args        args
		wantPos     int
		wantEnd     int
		wantString  string
		wantIndexes []int
	}{
		{
			name: "only wildcard on ident",
			args: args{
				end:   1,
				kind:  token.TokenKindIdent,
				value: "*",
			},
			wantEnd:     1,
			wantString:  "*",
			wantIndexes: []int{0},
		},
		{
			name: "only wildcard on string",
			args: args{
				pos:             1,
				end:             2,
				kind:            token.TokenKindString,
				value:           "*",
				withDoubleQuote: true,
			},
			wantPos:     1,
			wantEnd:     2,
			wantString:  `"*"`,
			wantIndexes: []int{1},
		},
		{
			name: "int value with wildcard on ident",
			args: args{
				end:   3,
				kind:  token.TokenKindIdent,
				value: "4*9",
			},
			wantEnd:     3,
			wantString:  "4*9",
			wantIndexes: []int{1},
		},
		{
			name: "int value with multi-wildcard on ident",
			args: args{
				end:   3,
				kind:  token.TokenKindIdent,
				value: "*0*",
			},
			wantEnd:     3,
			wantString:  "*0*",
			wantIndexes: []int{0, 2},
		},
		{
			name: "float value with wildcard on ident",
			args: args{
				end:   4,
				kind:  token.TokenKindIdent,
				value: "0.*9",
			},
			wantEnd:     4,
			wantString:  "0.*9",
			wantIndexes: []int{2},
		},
		{
			name: "float value with multi-wildcard on ident",
			args: args{
				end:   4,
				kind:  token.TokenKindIdent,
				value: "*.9*",
			},
			wantEnd:     4,
			wantString:  "*.9*",
			wantIndexes: []int{0, 3},
		},
		{
			name: "string value with wildcard on ident",
			args: args{
				end:   3,
				kind:  token.TokenKindIdent,
				value: "f*o",
			},
			wantEnd:     3,
			wantString:  "f*o",
			wantIndexes: []int{1},
		},
		{
			name: "string value with multi-wildcard on ident",
			args: args{
				end:   3,
				kind:  token.TokenKindIdent,
				value: "*o*",
			},
			wantEnd:     3,
			wantString:  "*o*",
			wantIndexes: []int{0, 2},
		},
		{
			name: "value with wildcard on string",
			args: args{
				end:             5,
				kind:            token.TokenKindString,
				value:           "f*o",
				withDoubleQuote: true,
			},
			wantEnd:     5,
			wantString:  `"f*o"`,
			wantIndexes: []int{2},
		},
		{
			name: "value with multi-wildcard on string",
			args: args{
				end:             5,
				kind:            token.TokenKindString,
				value:           "*o*",
				withDoubleQuote: true,
			},
			wantEnd:     5,
			wantString:  `"*o*"`,
			wantIndexes: []int{1, 3},
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
