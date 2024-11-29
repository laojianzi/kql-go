package parser_test

import (
	"strings"
	"testing"

	"github.com/laojianzi/kql-go/parser"
	"github.com/stretchr/testify/assert"
)

func FuzzParser(f *testing.F) {
	// Add initial corpus
	seeds := []string{
		"field:value",
		"field: value",
		"field : value",
		`field: "value"`,
		"field: *",
		"field: value*",
		"field: *value",
		"field: *value*",
		"field > 10",
		"field >= 10",
		"field < 10",
		"field <= 10",
		"field: true",
		"field: false",
		"field: null",
		"field1: value1 AND field2: value2",
		"field1: value1 OR field2: value2",
		"NOT field: value",
		"(field: value)",
		"(field1: value1) AND (field2: value2)",
		`field: "value with spaces"`,
		`field: "value with \"escaped\" quotes"`,
		`field: "value with \n newline"`,
		"field1: value1 AND field2: value2 OR field3: value3",
		"field1: (value1 OR value2) AND field2: value3",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, query string) {
		if strings.TrimSpace(query) == "" {
			return
		}

		// Current fuzzing implementation has limitations in input/output validation.
		// This test only covers basic safety checks:
		// 1. No panics during parsing
		// 2. String() output can be re-parsed
		// 3. String() output remains stable
		//
		// Contributions welcome for better validation approaches :)
		stmt, err := parser.New(query).Stmt()
		if err != nil || stmt == nil {
			return
		}

		stmt2, err := parser.New(stmt.String()).Stmt()
		assert.NoError(t, err)
		assert.Equal(t, stmt.String(), stmt2.String())
	})
}
