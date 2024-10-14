package token_test

import (
	"fmt"
	"testing"

	"github.com/laojianzi/kql-go/token"
	"github.com/stretchr/testify/assert"
)

func TestKind_String(t *testing.T) {
	t.Run("valid kind", func(t *testing.T) {
		kind := token.Kind(1)
		expected := token.TokenKindEof.String()
		actual := kind.String()
		assert.Equal(t, expected, actual)
	})

	t.Run("invalid kind", func(t *testing.T) {
		kind := token.Kind(999)
		expected := fmt.Sprintf("Kind(%d)", kind)
		actual := kind.String()
		assert.Equal(t, expected, actual)
	})

	t.Run("zero kind", func(t *testing.T) {
		kind := token.Kind(0)
		expected := fmt.Sprintf("Kind(%d)", kind)
		actual := kind.String()
		assert.Equal(t, expected, actual)
	})
}

func TestKind_IsKeyword(t *testing.T) {
	t.Run("valid keyword", func(t *testing.T) {
		for _, kind := range []token.Kind{
			token.TokenKindKeywordOr,
			token.TokenKindKeywordAnd,
			token.TokenKindKeywordNot,
		} {
			assert.True(t, kind.IsKeyword())
		}
	})

	t.Run("invalid keyword", func(t *testing.T) {
		assert.False(t, token.TokenKindOperatorEql.IsKeyword())
	})
}

func TestKind_IsOperator(t *testing.T) {
	t.Run("valid operator", func(t *testing.T) {
		for _, kind := range []token.Kind{
			token.TokenKindOperatorEql,
			token.TokenKindOperatorLss,
			token.TokenKindOperatorGtr,
			token.TokenKindOperatorLeq,
			token.TokenKindOperatorGeq,
		} {
			assert.True(t, kind.IsOperator())
		}
	})

	t.Run("invalid operator", func(t *testing.T) {
		assert.False(t, token.TokenKindKeywordNot.IsOperator())
		assert.False(t, token.TokenKindRparen.IsOperator())
	})
}

func TestKind_IsField(t *testing.T) {
	t.Run("valid field", func(t *testing.T) {
		assert.True(t, token.TokenKindIdent.IsField())
	})

	t.Run("invalid field", func(t *testing.T) {
		assert.False(t, token.TokenKindKeywordOr.IsField())
	})
}

func TestKind_IsValue(t *testing.T) {
	// Test cases for different kinds
	cases := []struct {
		name     string
		kind     token.Kind
		expected bool
	}{
		{"token ident is a value", token.TokenKindIdent, true},
		{"token int is a value", token.TokenKindInt, true},
		{"token float is a value", token.TokenKindFloat, true},
		{"token string is a value", token.TokenKindString, true},
		{"token keyword is not a value", token.TokenKindKeywordAnd, false}, // Test a non-matching kind
		{"token operator is not a value", token.TokenKindOperatorEql, false},
		{"token paren is not a value", token.TokenKindLparen, false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			actual := c.kind.IsValue()
			assert.Equal(t, c.expected, actual)
		})
	}
}

func TestIsKeyword(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"or is a keyword", "OR", true},
		{"and is a keyword", "AND", true},
		{"not is a keyword", "not", true},
		{"empty is not a keyword", "", false},
		{"operator is not a keyword", ":", false},
		{"paren is not a keyword", ")", false},
	}

	for _, c := range cases {
		actual := token.IsKeyword(c.input)
		if actual != c.expected {
			assert.Equal(t, c.expected, actual)
		}
	}
}

func TestIsOperator(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{": is a operator", ":", true},
		{"< is a operator", "<", true},
		{"> is a operator", ">", true},
		{"<= is a operator", "<=", true},
		{">= is a operator", ">=", true},
		{"+ is not a operator", "+", false},
		{"- is not a operator", "-", false},
		{"keyword is not a operator", "or", false},
		{"paren is not a operator", ")", false},
	}

	for _, c := range cases {
		actual := token.IsOperator(c.input)
		assert.Equal(t, c.expected, actual)
	}
}

func TestIsSpecialChar(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"operator is a special char", ":", true},
		{"paren is a special char", ")", true},
		{"keyword is not a special char", "or", false},
	}

	for _, c := range cases {
		actual := token.IsSpecialChar(c.input)
		assert.Equal(t, c.expected, actual)
	}
}

func TestToKeyword(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected token.Kind
	}{
		{"or is a keyword", "OR", token.TokenKindKeywordOr},
		{"and is a keyword", "AND", token.TokenKindKeywordAnd},
		{"not is a keyword", "not", token.TokenKindKeywordNot},
		{"operator is not a keyword", ":", token.TokenKindIllegal},
		{"paren is not a keyword", ")", token.TokenKindIllegal},
	}

	for _, c := range cases {
		actual := token.ToKeyword(c.input)
		assert.Equal(t, c.expected, actual)
	}
}

func TestToOperator(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected token.Kind
	}{
		{": is a operator", ":", token.TokenKindOperatorEql},
		{"< is a operator", "<", token.TokenKindOperatorLss},
		{"> is a operator", ">", token.TokenKindOperatorGtr},
		{"<= is a operator", "<=", token.TokenKindOperatorLeq},
		{">= is a operator", ">=", token.TokenKindOperatorGeq},
		{"keyword is not a operator", "or", token.TokenKindIllegal},
		{"paren is not a operator", ")", token.TokenKindIllegal},
	}

	for _, c := range cases {
		actual := token.ToOperator(c.input)
		assert.Equal(t, c.expected, actual)
	}
}

func TestKeywordsExpected(t *testing.T) {
	expected := "expected keyword OR|AND|NOT, but got \"test\""
	actual := token.KeywordsExpected("test")
	assert.Equal(t, expected, actual.Error())
}

func TestOperatorsExpected(t *testing.T) {
	expected := "expected operator :|<|>|<=|>=, but got \"test\""
	actual := token.OperatorsExpected("test")
	assert.Equal(t, expected, actual.Error())
}
