package parser

import (
	"fmt"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_peekString(t *testing.T) {
	t.Run("Double quote string", func(t *testing.T) {
		tests := []string{
			`" "`,
			`"\n"`,
			`"\t"`,
			`"\r"`,
			`"\r\n"`,
			`"\v"`,
			`"\f"`,
			`"0x85"`,
			`"0xA0"`,
			`"\\"`,
			`"hello world"`,
			`"123"`,
			`"1.23"`,
			`"*"`,
			`"*xx*"`,
			`":"`,
			`"<"`,
			`">"`,
			`"<="`,
			`">="`,
			`"OR"`,
			`"AND"`,
			`"NOT"`,
		}
		for _, s := range tests {
			t.Run(fmt.Sprintf("input: %s", s), func(t *testing.T) {
				lexer := &defaultLexer{input: s}
				tok, err := lexer.peekToken()
				require.NoError(t, err)
				require.Equal(t, TokenKindString, tok.Kind)
				require.Equal(t, strings.Trim(s, `"`), tok.Value)
				require.True(t, lexer.isEof())
			})
		}
	})

	t.Run("Invalid string", func(t *testing.T) {
		tests := []string{
			`"hello world     `,
		}

		for _, s := range tests {
			t.Run(fmt.Sprintf("input: %s", s), func(t *testing.T) {
				lexer := &defaultLexer{input: s}
				tok, err := lexer.peekToken()
				require.Nil(t, tok)
				require.EqualError(t, err, "expected double quote closed")
			})
		}
	})
}

func Test_peekNumber(t *testing.T) {
	t.Run("Int number", func(t *testing.T) {
		tests := []string{
			"1",
			"01",
			"10",
			"-1",
			"-01",
			"-10",
			fmt.Sprintf("%d", math.MinInt64),
			fmt.Sprintf("%d", math.MaxInt64),
		}

		for _, s := range tests {
			t.Run(fmt.Sprintf("input: %s", s), func(t *testing.T) {
				lexer := defaultLexer{input: s}
				tok, err := lexer.peekToken()
				require.NoError(t, err)
				require.Equal(t, TokenKindInt, tok.Kind)
				require.Equal(t, s, tok.Value)
				require.True(t, lexer.isEof())
			})
		}
	})

	t.Run("Float number", func(t *testing.T) {
		tests := []string{
			"1.0",
			"0.1",
			"10.01",
			"-1.0",
			"-0.1",
			"-10.01",
			fmt.Sprintf("%f", math.MaxFloat64),
			fmt.Sprintf("%f", math.MaxFloat64*-1-1),
		}

		for _, s := range tests {
			t.Run(fmt.Sprintf("input: %s", s), func(t *testing.T) {
				lexer := defaultLexer{input: s}
				tok, err := lexer.peekToken()
				require.NoError(t, err)
				require.Equal(t, TokenKindFloat, tok.Kind)
				require.Equal(t, s, tok.Value)
				require.True(t, lexer.isEof())
			})
		}
	})

	t.Run("Invalid number", func(t *testing.T) {
		tests := []struct {
			input  string
			errMsg string
		}{
			{"1.v", `expected digit, but got "v"`},
			{"01.v", `expected digit, but got "v"`},
			{"-1.", `expected digit, but got Eof`},
			{"1.0v", `expected digit or decimal point, but got "v"`},
			{"0.1v", `expected digit or decimal point, but got "v"`},
			{"10.01.", `expected digit, but got Eof`},
			{"-.1.0", `expected digit, but got "."`},
			{"-10.01v", `expected digit or decimal point, but got "v"`},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("input :%s and errMsg: %s", test.input, test.errMsg), func(t *testing.T) {
				lexer := defaultLexer{input: test.input}
				tok, err := lexer.peekToken()
				require.Nil(t, tok)
				require.EqualError(t, err, test.errMsg)
			})
		}
	})
}

func Test_peekOperator(t *testing.T) {
	tests := []struct {
		input string
		kind  Kind
	}{
		{":", TokenKindOperatorEql},
		{"<", TokenKindOperatorLss},
		{">", TokenKindOperatorGtr},
		{"<=", TokenKindOperatorLeq},
		{">=", TokenKindOperatorGeq},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("input: %s and kind: %d", test.input, test.kind), func(t *testing.T) {
			lexer := defaultLexer{input: test.input}
			tok, err := lexer.peekToken()
			require.NoError(t, err)
			require.Equal(t, test.kind, tok.Kind)
			require.Equal(t, test.input, tok.Value)
			require.True(t, lexer.isEof())
		})
	}
}

func Test_peekIdent(t *testing.T) {
	t.Run("Ident", func(t *testing.T) {
		tests := []string{
			"abc",
			"a.b.c",
			`a\"c`,
			`a\:c`,
			`a\<c`,
			`a\>c`,
			`a\\c`,
			`a\(c`,
			`a\)c`,
		}

		for _, s := range tests {
			t.Run(fmt.Sprintf("input: %s", s), func(t *testing.T) {
				lexer := defaultLexer{input: s}
				tok, err := lexer.peekToken()
				require.NoError(t, err)
				require.Equal(t, TokenKindIdent, tok.Kind)
				require.Equal(t, s, tok.Value)
				require.True(t, lexer.isEof())
			})
		}
	})

	t.Run("Keyword", func(t *testing.T) {
		tests := []struct {
			input string
			kind  Kind
		}{
			{"or", TokenKindKeywordOr},
			{"Or", TokenKindKeywordOr},
			{"OR", TokenKindKeywordOr},
			{"and", TokenKindKeywordAnd},
			{"And", TokenKindKeywordAnd},
			{"aNd", TokenKindKeywordAnd},
			{"anD", TokenKindKeywordAnd},
			{"ANd", TokenKindKeywordAnd},
			{"AnD", TokenKindKeywordAnd},
			{"aND", TokenKindKeywordAnd},
			{"AND", TokenKindKeywordAnd},
			{"not", TokenKindKeywordNot},
			{"Not", TokenKindKeywordNot},
			{"nOt", TokenKindKeywordNot},
			{"noT", TokenKindKeywordNot},
			{"NOt", TokenKindKeywordNot},
			{"NoT", TokenKindKeywordNot},
			{"nOT", TokenKindKeywordNot},
			{"NOT", TokenKindKeywordNot},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("input: %s and kind: %s", test.input, test.kind.String()), func(t *testing.T) {
				lexer := defaultLexer{input: test.input}
				tok, err := lexer.peekToken()
				require.NoError(t, err)
				require.Equal(t, test.kind, tok.Kind)
				require.Equal(t, test.input, tok.Value)
				require.True(t, lexer.isEof())
			})
		}
	})
}
