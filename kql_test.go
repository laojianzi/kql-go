package kql_test

import (
	"testing"

	"github.com/laojianzi/kql-go"
	"github.com/laojianzi/kql-go/token"
	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	err := kql.NewError("foo bar", 0, 4, token.KeywordsExpected("bar"))
	assert.Error(t, err)
	t.Logf("\n\n%v", err)
	assert.EqualError(t, err, "line 0:4 expected keyword OR|AND|NOT, but got \"bar\"\nfoo bar\n    ^\n")
}
