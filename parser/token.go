package parser

import (
	"fmt"
	"strings"
)

func init() {
	keywords = make(map[string]Kind, keywordEnd-keywordBeg-1)
	for i := keywordBeg + 1; i < keywordEnd; i++ {
		keywords[strings.ToUpper(tokenKinds[i])] = i
	}

	operators = make(map[string]Kind, operatorEnd-operatorBeg-1)
	for i := operatorBeg + 1; i < operatorEnd; i++ {
		operators[tokenKinds[i]] = i
	}
}

type Kind int

const (
	TokenKindIllegal Kind = iota
	TokenKindEof
	TokenKindInt
	TokenKindFloat
	TokenKindString
	TokenKindIdent
	keywordBeg
	TokenKindKeywordOr
	TokenKindKeywordAnd
	TokenKindKeywordNot
	keywordEnd
	operatorBeg
	TokenKindOperatorEql
	TokenKindOperatorLss
	TokenKindOperatorGtr
	TokenKindOperatorLeq
	TokenKindOperatorGeq
	operatorEnd
	TokenKindLparen
	TokenKindRparen
	TokenKindWildcard
)

var tokenKinds = [...]string{
	TokenKindIllegal:     "Illegal",
	TokenKindEof:         "Eof",
	TokenKindInt:         "Int",
	TokenKindFloat:       "Float",
	TokenKindString:      "String",
	TokenKindIdent:       "Ident",
	TokenKindKeywordOr:   "OR",
	TokenKindKeywordAnd:  "AND",
	TokenKindKeywordNot:  "NOT",
	TokenKindOperatorEql: ":",
	TokenKindOperatorLss: "<",
	TokenKindOperatorGtr: ">",
	TokenKindOperatorLeq: "<=",
	TokenKindOperatorGeq: ">=",
	TokenKindLparen:      "(",
	TokenKindRparen:      ")",
	TokenKindWildcard:    "*",
}

func (kind Kind) String() string {
	var s string
	if kind > 0 && int(kind) < len(tokenKinds) {
		s = tokenKinds[kind]
	}

	if s != "" {
		return s
	}

	return fmt.Sprintf("Kind(%d)", kind)
}

var (
	keywords  map[string]Kind
	operators map[string]Kind
)

func (kind Kind) IsKeyword() bool {
	return kind > keywordBeg && kind < keywordEnd
}

func (kind Kind) IsOperator() bool {
	return kind > operatorBeg && kind < operatorEnd
}

func (kind Kind) IsField() bool {
	return kind == TokenKindIdent
}

func (kind Kind) IsValue() bool {
	return kind == TokenKindIdent || kind == TokenKindInt || kind == TokenKindFloat || kind == TokenKindString
}

func IsKeyword(s string) bool {
	kind, ok := keywords[strings.ToUpper(s)]

	return ok && kind.IsKeyword()
}

func IsOperator(s string) bool {
	kind, ok := operators[s]

	return ok && kind.IsOperator()
}

func IsSpecialChar(s string) bool {
	return IsOperator(s) || s == TokenKindLparen.String() || s == TokenKindRparen.String()
}

func ToKeyword(s string) Kind {
	kind, ok := keywords[strings.ToUpper(s)]
	if ok && kind.IsKeyword() {
		return kind
	}

	return TokenKindIllegal
}

func ToOperator(s string) Kind {
	kind, ok := operators[s]
	if ok && kind.IsOperator() {
		return kind
	}

	return TokenKindIllegal
}

func KeywordsExpected(got string) error {
	var expectedList []string
	for k := keywordBeg + 1; k < keywordEnd; k++ {
		expectedList = append(expectedList, k.String())
	}

	return fmt.Errorf("expected keyword %s, but got %q", strings.Join(expectedList, "|"), got)
}

func OperatorsExpected(got string) error {
	var expectedList []string
	for k := operatorBeg + 1; k < operatorEnd; k++ {
		expectedList = append(expectedList, k.String())
	}

	return fmt.Errorf("expected operator %s, but got %q", strings.Join(expectedList, ""), got)
}

type Token struct {
	Pos   int
	End   int
	Kind  Kind
	Value string
}
