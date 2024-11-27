package parser

import (
	"testing"
)

func FuzzParser(f *testing.F) {
	// Add seed corpus with various query patterns
	seeds := []string{
		// Basic queries
		"foo:bar",
		"foo:bar AND bar:baz",
		"foo:bar OR bar:baz",
		"NOT foo:bar",

		// Complex queries
		"foo:bar AND (bar:baz OR qux:*)",
		"(status:active OR status:pending) AND age>=18",
		"foo:b* AND bar:*az AND NOT baz:qux",

		// Escape sequences
		"foo:\\*bar AND field:\\\"value\\\"",
		"field1:\\AND AND field2:\\OR OR field3:\\NOT",
		"path:\"C:\\\\Program Files\\\\*\"",

		// Numeric values
		"age>=18",
		"score>90.5",
		"count<=100",
		"value<-10.5",

		// Special characters
		"field:*value*",
		"name:?ohn*",
		"path:\"*/temp/*\"",

		// Mixed cases
		"Status:Active AND Age>=18 OR Name:\"John*\"",
		"(field1:value1 OR field2:value2) AND NOT field3:value3",
	}

	for _, seed := range seeds {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, query string) {
		p := New(query)
		_, err := p.Stmt()
		if err != nil {
			// Some errors are expected for invalid input
			return
		}
	})
}
