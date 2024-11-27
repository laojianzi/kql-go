package token

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"
)

// Kind represents token kind.
type Kind int

const (
	TokenKindIllegal Kind = iota // illegal token
	TokenKindEof                 // end of file
	TokenKindInt
	TokenKindFloat
	TokenKindString
	TokenKindIdent // identifier
	keywordBeg
	TokenKindKeywordOr  // keyword or
	TokenKindKeywordAnd // keyword and
	TokenKindKeywordNot // keyword not
	keywordEnd
	operatorBeg
	TokenKindOperatorEql // operator :
	TokenKindOperatorLss // operator <
	TokenKindOperatorGtr // operator >
	TokenKindOperatorLeq // operator <=
	TokenKindOperatorGeq // operator >=
	operatorEnd
	TokenKindLparen   // (
	TokenKindRparen   // )
	TokenKindWildcard // *
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

// String converts the Kind type to a string representation.
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
	keywords = func() map[string]Kind {
		v := make(map[string]Kind, keywordEnd-keywordBeg-1)
		for i := keywordBeg + 1; i < keywordEnd; i++ {
			v[strings.ToUpper(tokenKinds[i])] = i
		}

		return v
	}()

	operators = func() map[string]Kind {
		v := make(map[string]Kind, operatorEnd-operatorBeg-1)
		for i := operatorBeg + 1; i < operatorEnd; i++ {
			v[tokenKinds[i]] = i
		}

		return v
	}()
)

// IsKeyword checks if the Kind type is a keyword.
func (kind Kind) IsKeyword() bool {
	return kind > keywordBeg && kind < keywordEnd
}

// IsOperator checks if the Kind type is an operator.
func (kind Kind) IsOperator() bool {
	return kind > operatorBeg && kind < operatorEnd
}

// IsField checks if the Kind type is a field(identifier).
func (kind Kind) IsField() bool {
	return kind == TokenKindIdent
}

// IsValue checks if the Kind type is a value(identifier, int, float or string).
func (kind Kind) IsValue() bool {
	return kind == TokenKindIdent || kind == TokenKindInt || kind == TokenKindFloat || kind == TokenKindString
}

// IsKeyword checks if the string is a keyword.
func IsKeyword(s string) bool {
	kind, ok := keywords[strings.ToUpper(s)]

	return ok && kind.IsKeyword()
}

// IsOperator checks if the string is an operator.
func IsOperator(s string) bool {
	kind, ok := operators[s]

	return ok && kind.IsOperator()
}

// IsSpecialChar checks if the string is a special character(operator, ( or )).
func IsSpecialChar(s string) bool {
	return IsOperator(s) || s == TokenKindLparen.String() || s == TokenKindRparen.String()
}

var numberRegex = regexp.MustCompile(`^[+-]?\d+(\.\d+)?$`)

// IsNumber checks if the string is a number.
func IsNumber(s string) bool {
	return numberRegex.MatchString(s)
}

// ToKeyword converts the string to a keyword Kind type.
func ToKeyword(s string) Kind {
	kind, ok := keywords[strings.ToUpper(s)]
	if ok && kind.IsKeyword() {
		return kind
	}

	return TokenKindIllegal
}

// ToOperator converts the string to an operator Kind type.
func ToOperator(s string) Kind {
	kind, ok := operators[s]
	if ok && kind.IsOperator() {
		return kind
	}

	return TokenKindIllegal
}

// KeywordsExpected returns an error indicating that the given string was not the expected keyword.
func KeywordsExpected(got string) error {
	var expectedList []string
	for k := keywordBeg + 1; k < keywordEnd; k++ {
		expectedList = append(expectedList, k.String())
	}

	return fmt.Errorf("expected keyword %s, but got %q", strings.Join(expectedList, "|"), got)
}

// OperatorsExpected returns an error indicating that the given string was not the expected operator.
func OperatorsExpected(got string) error {
	var expectedList []string
	for k := operatorBeg + 1; k < operatorEnd; k++ {
		expectedList = append(expectedList, k.String())
	}

	return fmt.Errorf("expected operator %s, but got %q", strings.Join(expectedList, "|"), got)
}

// RequireEscape checks if a character requires escaping in the given context
func RequireEscape(s string, kind Kind) bool {
	if s == "" {
		return false
	}

	if r, _ := utf8.DecodeRuneInString(s); r == '"' || r == '\\' {
		return true
	}

	if kind == TokenKindString {
		return false
	}

	// only kind TokenKindIdent
	return IsSpecialChar(s) || IsKeyword(s)
}
