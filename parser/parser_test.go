package parser

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseMatchExpr(t *testing.T) {
	t.Run("only value", func(t *testing.T) {
		tests := []struct {
			input string
			kind  Kind
		}{
			{"*", TokenKindIdent},
			{"abc", TokenKindIdent},
			{"123", TokenKindInt},
			{"1.23", TokenKindFloat},
			{`"123"`, TokenKindString},
			{`"1.23"`, TokenKindString},
		}

		for _, test := range tests {
			for _, hasNot := range []bool{false, true} {
				if hasNot {
					test.input = "NOT " + test.input
				}

				t.Run(fmt.Sprintf("input: %s", test.input), func(t *testing.T) {
					expr, err := New(test.input).parseMatchExpr()
					require.NoError(t, err)

					pos, end := 0, len(test.input)
					if hasNot {
						pos += 4
					}

					if test.kind == TokenKindString {
						pos += 1
						end -= 1
					}

					require.EqualValues(t, expr, &MatchExpr{
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:             pos,
							end:             end,
							Kind:            test.kind,
							Value:           strings.Trim(strings.TrimPrefix(test.input, "NOT "), `"`),
							WithDoubleQuote: test.kind == TokenKindString,
						},
						HasNot: hasNot,
					})
				})
			}
		}
	})

	t.Run("match expr", func(t *testing.T) {
		tests := []struct {
			format string
			kind   Kind
		}{
			{"field%s*", TokenKindIdent},
			{"field%sabc", TokenKindIdent},
			{"field%s123", TokenKindInt},
			{"field%s1.23", TokenKindFloat},
			{`field%s"123"`, TokenKindString},
			{`field%s"1.23"`, TokenKindString},
		}

		for _, test := range tests {
			for _, operator := range operators {
				for _, space := range []string{"%s", " %s", "%s ", " %s "} {
					operatorWithSpace := fmt.Sprintf(space, operator.String())
					input := fmt.Sprintf(test.format, operatorWithSpace)

					for _, hasNot := range []bool{false, true} {
						if hasNot {
							input = "NOT " + input
						}

						t.Run(fmt.Sprintf("input: %s", input), func(t *testing.T) {
							expr, err := New(input).parseMatchExpr()
							require.NoError(t, err)

							pos, end := 5+len(operatorWithSpace), len(input)
							if hasNot {
								pos += 4
							}

							if test.kind == TokenKindString {
								pos += 1
								end -= 1
							}

							require.EqualValues(t, expr, &MatchExpr{
								Field:    "field",
								Operator: operator,
								Value:    &Literal{pos: pos, end: end, Kind: test.kind, Value: strings.Trim(strings.TrimPrefix(strings.TrimPrefix(input, "NOT "), fmt.Sprintf("field%s", operatorWithSpace)), `"`), WithDoubleQuote: test.kind == TokenKindString},
								HasNot:   hasNot,
							})
						})
					}
				}
			}
		}
	})

	t.Run("expected error", func(t *testing.T) {
		test := []struct {
			input string
			err   error
		}{
			{"", errors.New("expected field or value, but got Eof")},
			{":", fmt.Errorf("expected field or value, but got %q", ":")},
			{"<", fmt.Errorf("expected field or value, but got %q", "<")},
			{">", fmt.Errorf("expected field or value, but got %q", ">")},
			{"<=", fmt.Errorf("expected field or value, but got %q", "<=")},
			{">=", fmt.Errorf("expected field or value, but got %q", ">=")},
			{"OR", fmt.Errorf("expected field or value, but got %q", "OR")},
			{"AND", fmt.Errorf("expected field or value, but got %q", "AND")},
			{"field:", fmt.Errorf("expected value, but got %q", "Eof")},
			{"field<", fmt.Errorf("expected value, but got %q", "Eof")},
			{"field>", fmt.Errorf("expected value, but got %q", "Eof")},
			{"field<=", fmt.Errorf("expected value, but got %q", "Eof")},
			{"field>=", fmt.Errorf("expected value, but got %q", "Eof")},
			{"field::", fmt.Errorf("expected value, but got %q", ":")},
			{"field<:", fmt.Errorf("expected value, but got %q", ":")},
			{"field>:", fmt.Errorf("expected value, but got %q", ":")},
			{"field<=:", fmt.Errorf("expected value, but got %q", ":")},
			{"field>=:", fmt.Errorf("expected value, but got %q", ":")},
			{"field:<", fmt.Errorf("expected value, but got %q", "<")},
			{"field<<", fmt.Errorf("expected value, but got %q", "<")},
			{"field><", fmt.Errorf("expected value, but got %q", "<")},
			{"field<=<", fmt.Errorf("expected value, but got %q", "<")},
			{"field>=<", fmt.Errorf("expected value, but got %q", "<")},
			{"field:>", fmt.Errorf("expected value, but got %q", ">")},
			{"field<>", fmt.Errorf("expected value, but got %q", ">")},
			{"field>>", fmt.Errorf("expected value, but got %q", ">")},
			{"field<=>", fmt.Errorf("expected value, but got %q", ">")},
			{"field>=>", fmt.Errorf("expected value, but got %q", ">")},
			{"field:OR", fmt.Errorf("expected value, but got %q", "OR")},
			{"field<OR", fmt.Errorf("expected value, but got %q", "OR")},
			{"field>OR", fmt.Errorf("expected value, but got %q", "OR")},
			{"field<=OR", fmt.Errorf("expected value, but got %q", "OR")},
			{"field>=OR", fmt.Errorf("expected value, but got %q", "OR")},
			{"field:AND", fmt.Errorf("expected value, but got %q", "AND")},
			{"field<AND", fmt.Errorf("expected value, but got %q", "AND")},
			{"field>AND", fmt.Errorf("expected value, but got %q", "AND")},
			{"field<=AND", fmt.Errorf("expected value, but got %q", "AND")},
			{"field>=AND", fmt.Errorf("expected value, but got %q", "AND")},
			{"field:NOT", fmt.Errorf("expected value, but got %q", "NOT")},
			{"field<NOT", fmt.Errorf("expected value, but got %q", "NOT")},
			{"field>NOT", fmt.Errorf("expected value, but got %q", "NOT")},
			{"field<=NOT", fmt.Errorf("expected value, but got %q", "NOT")},
			{"field>=NOT", fmt.Errorf("expected value, but got %q", "NOT")},
		}

		for _, test := range test {
			t.Run(fmt.Sprintf("input: %s", test.input), func(t *testing.T) {
				_, err := New(test.input).parseMatchExpr()
				require.EqualError(t, err, test.err.Error())
			})
		}
	})
}

func Test_parseExpr(t *testing.T) {
	t.Run("success case", func(t *testing.T) {
		tests := []struct {
			input string
			expr  Expr
		}{
			{
				"field: value",
				&MatchExpr{
					Field:    "field",
					Operator: TokenKindOperatorEql,
					Value: &Literal{
						pos:   7,
						end:   12,
						Kind:  TokenKindIdent,
						Value: "value",
					},
				},
			},
			{
				"v1 AND v2 AND v3",
				&CombineExpr{
					LeftExpr: &MatchExpr{
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							end:   2,
							Kind:  TokenKindIdent,
							Value: "v1",
						},
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &CombineExpr{
						LeftExpr: &MatchExpr{
							pos:      7,
							Operator: TokenKindOperatorEql,
							Value: &Literal{
								pos:   7,
								end:   9,
								Kind:  TokenKindIdent,
								Value: "v2",
							},
						},
						Keyword: TokenKindKeywordAnd,
						RightExpr: &MatchExpr{
							pos:      14,
							Operator: TokenKindOperatorEql,
							Value: &Literal{
								pos:   14,
								end:   16,
								Kind:  TokenKindIdent,
								Value: "v3",
							},
						},
					},
				},
			},
			{
				"f1: v1 AND f2 > 2",
				&CombineExpr{
					LeftExpr: &MatchExpr{
						Field:    "f1",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   4,
							end:   6,
							Kind:  TokenKindIdent,
							Value: "v1",
						},
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &MatchExpr{
						pos:      11,
						Field:    "f2",
						Operator: TokenKindOperatorGtr,
						Value: &Literal{
							pos:   16,
							end:   17,
							Kind:  TokenKindInt,
							Value: "2",
						},
					},
				},
			},
			{
				"f1: v1 AND f2 > 2 OR f3: \"v3\"",
				&CombineExpr{
					LeftExpr: &MatchExpr{
						Field:    "f1",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   4,
							end:   6,
							Kind:  TokenKindIdent,
							Value: "v1",
						},
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &CombineExpr{
						LeftExpr: &MatchExpr{
							pos:      11,
							Field:    "f2",
							Operator: TokenKindOperatorGtr,
							Value: &Literal{
								pos:   16,
								end:   17,
								Kind:  TokenKindInt,
								Value: "2",
							},
						},
						Keyword: TokenKindKeywordOr,
						RightExpr: &MatchExpr{
							pos:      21,
							Field:    "f3",
							Operator: TokenKindOperatorEql,
							Value: &Literal{
								pos:             26,
								end:             28,
								Kind:            TokenKindString,
								Value:           "v3",
								WithDoubleQuote: true,
							},
						},
					},
				},
			},
			{
				"NOT f1: v1 AND f2 > 2 OR f3: \"v3\" AND NOT f4 < 4",
				&CombineExpr{
					LeftExpr: &MatchExpr{
						Field:    "f1",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   8,
							end:   10,
							Kind:  TokenKindIdent,
							Value: "v1",
						},
						HasNot: true,
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &CombineExpr{
						LeftExpr: &MatchExpr{
							pos:      15,
							Field:    "f2",
							Operator: TokenKindOperatorGtr,
							Value: &Literal{
								pos:   20,
								end:   21,
								Kind:  TokenKindInt,
								Value: "2",
							},
						},
						Keyword: TokenKindKeywordOr,
						RightExpr: &CombineExpr{
							LeftExpr: &MatchExpr{
								pos:      25,
								Field:    "f3",
								Operator: TokenKindOperatorEql,
								Value: &Literal{
									pos:             30,
									end:             32,
									Kind:            TokenKindString,
									Value:           "v3",
									WithDoubleQuote: true,
								},
							},
							Keyword: TokenKindKeywordAnd,
							RightExpr: &MatchExpr{
								pos:      38,
								Field:    "f4",
								Operator: TokenKindOperatorLss,
								Value: &Literal{
									pos:   47,
									end:   48,
									Kind:  TokenKindInt,
									Value: "4",
								},
								HasNot: true,
							},
						},
					},
				},
			},
			{
				"NOT f1: v1 AND f2 > 2 OR f3: \"v3\" AND NOT f4 < 4 OR NOT f5 <= 5.1",
				&CombineExpr{
					LeftExpr: &MatchExpr{
						Field:    "f1",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   8,
							end:   10,
							Kind:  TokenKindIdent,
							Value: "v1",
						},
						HasNot: true,
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &CombineExpr{
						LeftExpr: &MatchExpr{
							pos:      15,
							Field:    "f2",
							Operator: TokenKindOperatorGtr,
							Value: &Literal{
								pos:   20,
								end:   21,
								Kind:  TokenKindInt,
								Value: "2",
							},
						},
						Keyword: TokenKindKeywordOr,
						RightExpr: &CombineExpr{
							LeftExpr: &MatchExpr{
								pos:      25,
								Field:    "f3",
								Operator: TokenKindOperatorEql,
								Value: &Literal{
									pos:             30,
									end:             32,
									Kind:            TokenKindString,
									Value:           "v3",
									WithDoubleQuote: true,
								},
							},
							Keyword: TokenKindKeywordAnd,
							RightExpr: &CombineExpr{
								LeftExpr: &MatchExpr{
									pos:      38,
									Field:    "f4",
									Operator: TokenKindOperatorLss,
									Value: &Literal{
										pos:   47,
										end:   48,
										Kind:  TokenKindInt,
										Value: "4",
									},
									HasNot: true,
								},
								Keyword: TokenKindKeywordOr,
								RightExpr: &MatchExpr{
									pos:      52,
									Field:    "f5",
									Operator: TokenKindOperatorLeq,
									Value: &Literal{
										pos:   62,
										end:   65,
										Kind:  TokenKindFloat,
										Value: "5.1",
									},
									HasNot: true,
								},
							},
						},
					},
				},
			},
			{
				"NOT f1: v1 AND f2 > 2 OR f3: \"v3\" AND NOT f4 < 4 OR NOT f5 <= 5.1 AND f6 >= 0.0000000001",
				&CombineExpr{
					LeftExpr: &MatchExpr{
						Field:    "f1",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   8,
							end:   10,
							Kind:  TokenKindIdent,
							Value: "v1",
						},
						HasNot: true,
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &CombineExpr{
						LeftExpr: &MatchExpr{
							pos:      15,
							Field:    "f2",
							Operator: TokenKindOperatorGtr,
							Value: &Literal{
								pos:   20,
								end:   21,
								Kind:  TokenKindInt,
								Value: "2",
							},
						},
						Keyword: TokenKindKeywordOr,
						RightExpr: &CombineExpr{
							LeftExpr: &MatchExpr{
								pos:      25,
								Field:    "f3",
								Operator: TokenKindOperatorEql,
								Value: &Literal{
									pos:             30,
									end:             32,
									Kind:            TokenKindString,
									Value:           "v3",
									WithDoubleQuote: true,
								},
							},
							Keyword: TokenKindKeywordAnd,
							RightExpr: &CombineExpr{
								LeftExpr: &MatchExpr{
									pos:      38,
									Field:    "f4",
									Operator: TokenKindOperatorLss,
									Value: &Literal{
										pos:   47,
										end:   48,
										Kind:  TokenKindInt,
										Value: "4",
									},
									HasNot: true,
								},
								Keyword: TokenKindKeywordOr,
								RightExpr: &CombineExpr{
									LeftExpr: &MatchExpr{
										pos:      52,
										Field:    "f5",
										Operator: TokenKindOperatorLeq,
										Value: &Literal{
											pos:   62,
											end:   65,
											Kind:  TokenKindFloat,
											Value: "5.1",
										},
										HasNot: true,
									},
									Keyword: TokenKindKeywordAnd,
									RightExpr: &MatchExpr{
										pos:      70,
										Field:    "f6",
										Operator: TokenKindOperatorGeq,
										Value: &Literal{
											pos:   76,
											end:   88,
											Kind:  TokenKindFloat,
											Value: "0.0000000001",
										},
									},
								},
							},
						},
					},
				},
			},
			{
				"f1: v1 AND (f2 > 2 OR f3 < 3)",
				&CombineExpr{
					LeftExpr: &MatchExpr{
						Field:    "f1",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   4,
							end:   6,
							Kind:  TokenKindIdent,
							Value: "v1",
						},
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &WrapExpr{
						pos:    11,
						Layers: 1,
						Expr: &CombineExpr{
							LeftExpr: &MatchExpr{
								pos:      12,
								Field:    "f2",
								Operator: TokenKindOperatorGtr,
								Value: &Literal{
									pos:   17,
									end:   18,
									Kind:  TokenKindInt,
									Value: "2",
								},
							},
							Keyword: TokenKindKeywordOr,
							RightExpr: &MatchExpr{
								pos:      22,
								Field:    "f3",
								Operator: TokenKindOperatorLss,
								Value: &Literal{
									pos:   27,
									end:   28,
									Kind:  TokenKindInt,
									Value: "3",
								},
							},
						},
					},
				},
			},
			{
				"((f1 > 1 OR f2 < 2)) AND f3: v3",
				&CombineExpr{
					LeftExpr: &WrapExpr{
						Layers: 2,
						Expr: &CombineExpr{
							LeftExpr: &MatchExpr{
								pos:      2,
								Field:    "f1",
								Operator: TokenKindOperatorGtr,
								Value: &Literal{
									pos:   7,
									end:   8,
									Kind:  TokenKindInt,
									Value: "1",
								},
							},
							Keyword: TokenKindKeywordOr,
							RightExpr: &MatchExpr{
								pos:      12,
								Field:    "f2",
								Operator: TokenKindOperatorLss,
								Value: &Literal{
									pos:   17,
									end:   18,
									Kind:  TokenKindInt,
									Value: "2",
								},
							},
						},
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &MatchExpr{
						pos:      25,
						Field:    "f3",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   29,
							end:   31,
							Kind:  TokenKindIdent,
							Value: "v3",
						},
					},
				},
			},
			{
				"f1: v1 AND (f2 > 2 OR f3 < 3) AND f4: 4",
				&CombineExpr{
					LeftExpr: &MatchExpr{
						Field:    "f1",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   4,
							end:   6,
							Kind:  TokenKindIdent,
							Value: "v1",
						},
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &CombineExpr{
						LeftExpr: &WrapExpr{
							pos:    11,
							Layers: 1,
							Expr: &CombineExpr{
								LeftExpr: &MatchExpr{
									pos:      12,
									Field:    "f2",
									Operator: TokenKindOperatorGtr,
									Value: &Literal{
										pos:   17,
										end:   18,
										Kind:  TokenKindInt,
										Value: "2",
									},
								},
								Keyword: TokenKindKeywordOr,
								RightExpr: &MatchExpr{
									pos:      22,
									Field:    "f3",
									Operator: TokenKindOperatorLss,
									Value: &Literal{
										pos:   27,
										end:   28,
										Kind:  TokenKindInt,
										Value: "3",
									},
								},
							},
						},
						Keyword: TokenKindKeywordAnd,
						RightExpr: &MatchExpr{
							pos:      34,
							Field:    "f4",
							Operator: TokenKindOperatorEql,
							Value: &Literal{
								pos:   38,
								end:   39,
								Kind:  TokenKindInt,
								Value: "4",
							},
						},
					},
				},
			},
			{
				"f1: (v1 OR v11 AND v111) AND f2: v2",
				&CombineExpr{
					LeftExpr: &WrapExpr{
						Field:  "f1",
						Layers: 1,
						Expr: &CombineExpr{
							LeftExpr: &MatchExpr{
								pos:      5,
								Operator: TokenKindOperatorEql,
								Value: &Literal{
									pos:   5,
									end:   7,
									Kind:  TokenKindIdent,
									Value: "v1",
								},
							},
							Keyword: TokenKindKeywordOr,
							RightExpr: &CombineExpr{
								LeftExpr: &MatchExpr{
									pos:      11,
									Operator: TokenKindOperatorEql,
									Value: &Literal{
										pos:   11,
										end:   14,
										Kind:  TokenKindIdent,
										Value: "v11",
									},
									HasNot: false,
								},
								Keyword: TokenKindKeywordAnd,
								RightExpr: &MatchExpr{
									pos:      19,
									Operator: TokenKindOperatorEql,
									Value: &Literal{
										pos:   19,
										end:   23,
										Kind:  TokenKindIdent,
										Value: "v111",
									},
								},
							},
						},
					},
					Keyword: TokenKindKeywordAnd,
					RightExpr: &MatchExpr{
						pos:      29,
						Field:    "f2",
						Operator: TokenKindOperatorEql,
						Value: &Literal{
							pos:   33,
							end:   35,
							Kind:  TokenKindIdent,
							Value: "v2",
						},
					},
				},
			},
			{
				"(f2: v2 AND f1: (((v1 OR v11 AND v111))))",
				&WrapExpr{
					Layers: 1,
					Expr: &CombineExpr{
						LeftExpr: &MatchExpr{
							pos:      1,
							Field:    "f2",
							Operator: TokenKindOperatorEql,
							Value: &Literal{
								pos:   5,
								end:   7,
								Kind:  TokenKindIdent,
								Value: "v2",
							},
						},
						Keyword: TokenKindKeywordAnd,
						RightExpr: &WrapExpr{
							pos:    12,
							Field:  "f1",
							Layers: 3,
							Expr: &CombineExpr{
								LeftExpr: &MatchExpr{
									pos:      19,
									Operator: TokenKindOperatorEql,
									Value: &Literal{
										pos:   19,
										end:   21,
										Kind:  TokenKindIdent,
										Value: "v1",
									},
								},
								Keyword: TokenKindKeywordOr,
								RightExpr: &CombineExpr{
									LeftExpr: &MatchExpr{
										pos:      25,
										Operator: TokenKindOperatorEql,
										Value: &Literal{
											pos:   25,
											end:   28,
											Kind:  TokenKindIdent,
											Value: "v11",
										},
										HasNot: false,
									},
									Keyword: TokenKindKeywordAnd,
									RightExpr: &MatchExpr{
										pos:      33,
										Operator: TokenKindOperatorEql,
										Value: &Literal{
											pos:   33,
											end:   37,
											Kind:  TokenKindIdent,
											Value: "v111",
										},
									},
								},
							},
						},
					},
				},
			},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("input: %s", test.input), func(t *testing.T) {
				expr, err := New(test.input).parseExpr()
				require.NoError(t, err)
				require.EqualValues(t, test.expr, expr)
			})
		}
	})

	t.Run("fail case", func(t *testing.T) {
		tests := []struct {
			input  string
			errMsg string
		}{
			{"AND f1: v1", "expected field or value, but got \"AND\""},
			{"f1: v1 AND (", "expected field or value, but got Eof"},
			{"f1: v1 AND ()", "expected field or value, but got \")\""},
			{"f1: v1 AND (f2 > 2", "expected token <Rparen>, but got \"Eof\""},
			{"f1: (", "expected field or value, but got Eof"},
			{"f1: ()", "expected field or value, but got \")\""},
			{"f1: (f2 > 2", "expected token <Rparen>, but got \"Eof\""},
		}

		for _, test := range tests {
			t.Run(fmt.Sprintf("input: %s", test.input), func(t *testing.T) {
				_, err := New(test.input).parseExpr()
				require.EqualError(t, err, test.errMsg)
			})
		}
	})
}
