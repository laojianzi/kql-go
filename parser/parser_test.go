package parser_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/laojianzi/kql-go"
	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/parser"
	"github.com/laojianzi/kql-go/token"
)

func Test_defaultParser_Stmt(t *testing.T) {
	t.Run("simple", func(t *testing.T) {
		cases := []struct {
			input string
			want  ast.Expr
		}{
			{
				input: "foo",
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindIdent, "foo", nil), false),
			},
			{
				input: "1",
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 1, token.TokenKindInt, "1", nil), false),
			},
			{
				input: "0.1",
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindFloat, "0.1", nil), false),
			},
			{
				input: `"0.1"`,
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 5, token.TokenKindString, "0.1", nil), false),
			},
			{
				input: `f1: "v1"`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1", nil), false),
			},
			{
				input: `f1 > 1`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorGtr, ast.NewLiteral(5, 6, token.TokenKindInt, "1", nil), false),
			},
			{
				input: `f1 < 1.1`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorLss, ast.NewLiteral(5, 8, token.TokenKindFloat, "1.1", nil), false),
			},
			{
				input: `f1 >= 1000.0001`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorGeq, ast.NewLiteral(6, 15, token.TokenKindFloat, "1000.0001", nil), false),
			},
			{
				input: `f1 <= 100000011`,
				want:  ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorLeq, ast.NewLiteral(6, 15, token.TokenKindInt, "100000011", nil), false),
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
				want:  ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(4, 7, token.TokenKindIdent, "bar", nil), true),
			},
			{
				input: "foo AND bar",
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindIdent, "foo", nil), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(8, "", 0, ast.NewLiteral(8, 11, token.TokenKindIdent, "bar", nil), false),
				),
			},
			{
				input: "foo AND NOT bar",
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindIdent, "foo", nil), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(8, "", 0, ast.NewLiteral(12, 15, token.TokenKindIdent, "bar", nil), true),
				),
			},
			{
				input: `v1 AND 2 OR 0.3 AND NOT "v4" OR NOT 5.0`,
				want: ast.NewCombineExpr(
					ast.NewCombineExpr(
						ast.NewCombineExpr(
							ast.NewCombineExpr(
								ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 2, token.TokenKindIdent, "v1", nil), false),
								token.TokenKindKeywordAnd,
								ast.NewBinaryExpr(7, "", 0, ast.NewLiteral(7, 8, token.TokenKindInt, "2", nil), false),
							),
							token.TokenKindKeywordOr,
							ast.NewBinaryExpr(12, "", 0, ast.NewLiteral(12, 15, token.TokenKindFloat, "0.3", nil), false),
						),
						token.TokenKindKeywordAnd,
						ast.NewBinaryExpr(20, "", 0, ast.NewLiteral(24, 28, token.TokenKindString, "v4", nil), true),
					),
					token.TokenKindKeywordOr,
					ast.NewBinaryExpr(32, "", 0, ast.NewLiteral(36, 39, token.TokenKindFloat, "5.0", nil), true),
				),
			},
			{
				input: "NOT f: v",
				want:  ast.NewBinaryExpr(0, "f", token.TokenKindOperatorEql, ast.NewLiteral(7, 8, token.TokenKindIdent, "v", nil), true),
			},
			{
				input: `f1: "v1" AND f2 > 2`,
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1", nil), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(13, "f2", token.TokenKindOperatorGtr, ast.NewLiteral(18, 19, token.TokenKindInt, "2", nil), false),
				),
			},
			{
				input: `f1: "v1" AND NOT f2 > 2`,
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1", nil), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(13, "f2", token.TokenKindOperatorGtr, ast.NewLiteral(22, 23, token.TokenKindInt, "2", nil), true),
				),
			},
			{
				input: `f1: "v1" AND f2 > 2 OR f3 < 0.3 AND NOT f4 >= 4 OR NOT f5 <= 5.0`,
				want: ast.NewCombineExpr(
					ast.NewCombineExpr(
						ast.NewCombineExpr(
							ast.NewCombineExpr(
								ast.NewBinaryExpr(0, "f1", token.TokenKindOperatorEql, ast.NewLiteral(4, 8, token.TokenKindString, "v1", nil), false),
								token.TokenKindKeywordAnd,
								ast.NewBinaryExpr(13, "f2", token.TokenKindOperatorGtr, ast.NewLiteral(18, 19, token.TokenKindInt, "2", nil), false),
							),
							token.TokenKindKeywordOr,
							ast.NewBinaryExpr(23, "f3", token.TokenKindOperatorLss, ast.NewLiteral(28, 31, token.TokenKindFloat, "0.3", nil), false),
						),
						token.TokenKindKeywordAnd,
						ast.NewBinaryExpr(36, "f4", token.TokenKindOperatorGeq, ast.NewLiteral(46, 47, token.TokenKindInt, "4", nil), true),
					),
					token.TokenKindKeywordOr,
					ast.NewBinaryExpr(51, "f5", token.TokenKindOperatorLeq, ast.NewLiteral(61, 64, token.TokenKindFloat, "5.0", nil), true),
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

	t.Run("with paren", func(t *testing.T) {
		cases := []struct {
			input string
			want  ast.Expr
		}{
			{
				input: "foo AND (NOT bar)",
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "", 0, ast.NewLiteral(0, 3, token.TokenKindIdent, "foo", nil), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(8, "", 0, ast.NewParenExpr(8, 17, ast.NewBinaryExpr(9, "", 0, ast.NewLiteral(13, 16, token.TokenKindIdent, "bar", nil), true)), false),
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

	t.Run("with wildcard", func(t *testing.T) {
		cases := []struct {
			input string
			want  ast.Expr
		}{
			{
				input: "foo: *",
				want: ast.NewBinaryExpr(0, "foo", token.TokenKindOperatorEql, ast.NewWildcardExpr(
					ast.NewLiteral(5, 6, token.TokenKindIdent, "*", nil),
					[]int{0},
				), false),
			},
			{
				input: `foo: * AND bar: *v2`,
				want: ast.NewCombineExpr(
					ast.NewBinaryExpr(0, "foo", token.TokenKindOperatorEql, ast.NewWildcardExpr(
						ast.NewLiteral(5, 6, token.TokenKindIdent, "*", nil),
						[]int{0},
					), false),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(11, "bar", token.TokenKindOperatorEql, ast.NewWildcardExpr(
						ast.NewLiteral(16, 19, token.TokenKindIdent, "*v2", nil),
						[]int{0},
					), false),
				),
			},
			{
				input: "foo: v*1",
				want: ast.NewBinaryExpr(0, "foo", token.TokenKindOperatorEql, ast.NewWildcardExpr(
					ast.NewLiteral(5, 8, token.TokenKindIdent, "v*1", nil),
					[]int{1},
				), false),
			},
			{
				input: "foo: *0 AND bar: 1* AND 2*0",
				want: ast.NewCombineExpr(
					ast.NewCombineExpr(
						ast.NewBinaryExpr(0, "foo", token.TokenKindOperatorEql, ast.NewWildcardExpr(
							ast.NewLiteral(5, 7, token.TokenKindIdent, "*0", nil),
							[]int{0},
						), false),
						token.TokenKindKeywordAnd,
						ast.NewBinaryExpr(12, "bar", token.TokenKindOperatorEql, ast.NewWildcardExpr(
							ast.NewLiteral(17, 19, token.TokenKindIdent, "1*", nil),
							[]int{1},
						), false),
					),
					token.TokenKindKeywordAnd,
					ast.NewBinaryExpr(24, "", 0, ast.NewWildcardExpr(
						ast.NewLiteral(24, 27, token.TokenKindIdent, "2*0", nil),
						[]int{1},
					), false),
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

	t.Run("expect kql error", func(t *testing.T) {
		cases := []struct {
			input string
			want  error
		}{
			{
				input: "foo bar",
				want: kql.NewError(
					"foo bar",
					token.TokenKindIdent,
					"bar",
					4,
					token.KeywordsExpected(token.TokenKindIdent.String()),
				),
			},
			{
				input: "foo: ",
				want: kql.NewError(
					"foo:",
					token.TokenKindEof,
					"",
					4,
					errors.New("unexpected token: Eof"),
				),
			},
		}

		for _, c := range cases {
			t.Run(c.input, func(t *testing.T) {
				_, err := parser.New(c.input).Stmt()
				assert.Error(t, err)
				assert.Equal(t, c.want.Error(), err.Error())
			})
		}
	})
}

func TestParser_EscapedKeywords(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "escaped AND keyword",
			input:    `app: \AND`,
			expected: `app: \AND`,
			wantErr:  false,
		},
		{
			name:     "escaped OR keyword",
			input:    `app: \OR`,
			expected: `app: \OR`,
			wantErr:  false,
		},
		{
			name:     "escaped NOT keyword",
			input:    `app: \NOT`,
			expected: `app: \NOT`,
			wantErr:  false,
		},
		{
			name:     "mix of escaped and normal keywords",
			input:    `app: foo AND msg: \OR`,
			expected: `app: foo AND msg: \OR`,
			wantErr:  false,
		},
		{
			name:     "multiple escaped keywords",
			input:    `app: \AND AND msg: \OR`,
			expected: `app: \AND AND msg: \OR`,
			wantErr:  false,
		},
		{
			name:     "escaped keyword in parentheses",
			input:    `(app: \AND)`,
			expected: `(app: \AND)`,
			wantErr:  false,
		},
		{
			name:     "escaped backslash before keyword",
			input:    `app: \\ AND msg: foo`,
			expected: `app: \\ AND msg: foo`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.input)
			expr, err := p.Stmt()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, expr.String())
			}

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestParser_EscapedOperators(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "escaped equality operator",
			input:    `app: \:`,
			expected: `app: \:`,
			wantErr:  false,
		},
		{
			name:     "escaped less than operator",
			input:    `app: \>`,
			expected: `app: \>`,
			wantErr:  false,
		},
		{
			name:     "escaped greater than operator",
			input:    `app: \<`,
			expected: `app: \<`,
			wantErr:  false,
		},
		{
			name:     "escaped less than or equal operator",
			input:    `app: \<=`,
			expected: `app: \<=`,
			wantErr:  false,
		},
		{
			name:     "escaped greater than or equal operator",
			input:    `app: \>=`,
			expected: `app: \>=`,
			wantErr:  false,
		},
		{
			name:     "mix of escaped operators and normal keywords",
			input:    `app: foo AND msg: \>`,
			expected: `app: foo AND msg: \>`,
			wantErr:  false,
		},
		{
			name:     "escaped operator in parentheses",
			input:    `(app: \>)`,
			expected: `(app: \>)`,
			wantErr:  false,
		},
		{
			name:     "escaped backslash before operator",
			input:    `app: \: AND msg: foo`,
			expected: `app: \: AND msg: foo`,
			wantErr:  false,
		},
		{
			name:     "escaped backslash before operator 2",
			input:    `app: \`,
			expected: ``,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.input)
			expr, err := p.Stmt()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, expr.String())
			}

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestParser_EscapedWildcard(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "escaped wildcard",
			input:    `app: \*`,
			expected: `app: \*`,
			wantErr:  false,
		},
		{
			name:     "escaped wildcard in parentheses",
			input:    `(app: \*)`,
			expected: `(app: \*)`,
			wantErr:  false,
		},
		{
			name:     "escaped wildcard with AND",
			input:    `app: \* AND msg: 5*0`,
			expected: `app: \* AND msg: 5*0`,
			wantErr:  false,
		},
		{
			name:     "multiple escaped wildcards",
			input:    `\*escaped\*`,
			expected: `\*escaped\*`,
			wantErr:  false,
		},
		{
			name:     "escaped wildcard in string",
			input:    `message: "Hello \* World"`,
			expected: `message: "Hello \* World"`,
			wantErr:  false,
		},
		{
			name:     "escaped wildcard at string boundaries",
			input:    `message: "\*Hello World\*"`,
			expected: `message: "\*Hello World\*"`,
			wantErr:  false,
		},
		{
			name:     "escaped wildcard with other escapes",
			input:    `message: "\\* \"Hello\" \* World"`,
			expected: `message: "\\* \"Hello\" \* World"`,
			wantErr:  false,
		},
		{
			name:     "invalid escape sequence",
			input:    `message: \a`,
			expected: ``,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.input)
			expr, err := p.Stmt()

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected, expr.String())
		})
	}
}

func TestParser_EscapedDoubleQuote(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "escaped a double quote in ident",
			input:    `app: \"`,
			expected: `app: \"`,
			wantErr:  false,
		},
		{
			name:     "escaped double quote warped in ident",
			input:    `app: \"v1\"`,
			expected: `app: \"v1\"`,
			wantErr:  false,
		},
		{
			name:     "escaped double quote in ident",
			input:    `foo: b\"ar\"`,
			expected: `foo: b\"ar\"`,
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.input)
			expr, err := p.Stmt()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, expr.String())
			}

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestParser_EscapedString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "escaped double quote in string",
			input:    `message: "hello \"world\""`,
			expected: `message: "hello \"world\""`,
			wantErr:  false,
		},
		{
			name:     "escaped backslash in string",
			input:    `message: "C:\\Program Files"`,
			expected: `message: "C:\\Program Files"`,
			wantErr:  false,
		},
		{
			name:     "escaped asterisk in string",
			input:    `message: "hello \* world"`,
			expected: `message: "hello \* world"`,
			wantErr:  false,
		},
		{
			name:     "multiple escaped characters",
			input:    `message: "path: \"C:\\Program Files\\*\""`,
			expected: `message: "path: \"C:\\Program Files\\*\""`,
			wantErr:  false,
		},
		{
			name:     "escaped characters in complex query",
			input:    `field1: "value with \"quotes\"" AND field2: "\*wildcard\*" OR field3: "back\\slash"`,
			expected: `field1: "value with \"quotes\"" AND field2: "\*wildcard\*" OR field3: "back\\slash"`,
			wantErr:  false,
		},
		{
			name:     "escaped characters in complex query with escaped double quotes",
			input:    `foo:\*bar AND field:\"value\"`,
			expected: `foo: \*bar AND field: \"value\"`,
			wantErr:  false,
		},
		{
			name:    "invalid escape sequence",
			input:   `message: "hello \n world"`,
			wantErr: true,
		},
		{
			name:    "unclosed string",
			input:   `message: "unclosed`,
			wantErr: true,
		},
		{
			name:    "invalid escape at end",
			input:   `message: "test\`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.input)
			expr, err := p.Stmt()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, expr.String())
			}

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}

func TestParser_EdgeCases(t *testing.T) {
	tests := []struct {
		name  string
		query string
	}{
		{
			name:  "empty query",
			query: "",
		},
		{
			name:  "only whitespace",
			query: "    \t\n",
		},
		{
			name:  "only operator",
			query: "AND",
		},
		{
			name:  "incomplete field value",
			query: "field:",
		},
		{
			name:  "missing value after operator",
			query: "age >",
		},
		{
			name:  "invalid numeric value",
			query: "age > abc",
		},
		{
			name:  "unmatched parenthesis",
			query: "(field: value",
		},
		{
			name:  "extra parenthesis",
			query: "field: value)",
		},
		{
			name:  "consecutive operators",
			query: "field: value AND OR value2",
		},
		{
			name:  "invalid field name",
			query: "field space: value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.query)
			_, err := p.Stmt()
			assert.Error(t, err)
		})
	}
}

func TestParser_ComplexQueries(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantStr string
		wantErr bool
	}{
		{
			name:    "nested parentheses",
			query:   `((field1: value1 OR field2: value2) AND (field3: value3 OR field4: value4))`,
			wantStr: `((field1: value1 OR field2: value2) AND (field3: value3 OR field4: value4))`,
		},
		{
			name:    "mixed operators",
			query:   `field1: value1 AND NOT (field2: value2 OR field3: value3)`,
			wantStr: `field1: value1 AND NOT (field2: value2 OR field3: value3)`,
		},
		{
			name:    "multiple wildcards",
			query:   `field1: *val* AND field2: val?ue*`,
			wantStr: `field1: *val* AND field2: val?ue*`,
		},
		{
			name:    "mixed comparisons",
			query:   `age >= 18 AND score > 90 OR rank <= 3`,
			wantStr: `age >= 18 AND score > 90 OR rank <= 3`,
		},
		{
			name:    "complex escapes",
			query:   `message: "Hello \"World\"" AND path: "C:\\Program Files\\*" OR command: "\"quoted\""`,
			wantStr: `message: "Hello \"World\"" AND path: "C:\\Program Files\\*" OR command: "\"quoted\""`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.New(tt.query)
			stmt, err := p.Stmt()

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantStr, stmt.String())
			}

			assert.Equal(t, tt.wantErr, err != nil)
		})
	}
}
