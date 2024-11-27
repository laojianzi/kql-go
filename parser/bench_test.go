package parser

import (
	"testing"
)

var benchmarkQueries = []struct {
	name  string
	query string
}{
	{
		name:  "simple_field",
		query: `status: "active"`,
	},
	{
		name:  "numeric_comparison",
		query: `age >= 18`,
	},
	{
		name:  "multiple_conditions",
		query: `status: "active" AND age >= 18`,
	},
	{
		name:  "complex_query",
		query: `(status: "active" OR status: "pending") AND age >= 18 AND name: "john*"`,
	},
	{
		name:  "escaped_chars",
		query: `message: "Hello \"World\"" AND path: "C:\\Program Files\\*"`,
	},
	{
		name:  "many_conditions",
		query: `status: "active" AND age >= 18 AND name: "john*" AND city: "New York" AND country: "USA" AND role: "admin"`,
	},
}

func BenchmarkParser(b *testing.B) {
	for _, bq := range benchmarkQueries {
		b.Run(bq.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				stmt, err := New(bq.query).Stmt()
				if err != nil {
					b.Fatal(err)
				}

				_ = stmt.String()
			}
		})
	}
}

func BenchmarkParserParallel(b *testing.B) {
	for _, bq := range benchmarkQueries {
		b.Run(bq.name, func(b *testing.B) {
			b.ReportAllocs()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					stmt, err := New(bq.query).Stmt()
					if err != nil {
						b.Fatal(err)
					}

					_ = stmt.String()
				}
			})
		})
	}
}

func BenchmarkLexer(b *testing.B) {
	for _, bq := range benchmarkQueries {
		b.Run(bq.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				l := newLexer(bq.query)

				for {
					if l.nextToken(); l.eof() {
						break
					}
				}
			}
		})
	}
}

func BenchmarkEscapeSequence(b *testing.B) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "no_escape",
			input: `hello world`,
		},
		{
			name:  "single_escape",
			input: `hello \"world\"`,
		},
		{
			name:  "multiple_escapes",
			input: `\"hello\" \"world\"`,
		},
		{
			name:  "mixed_escapes",
			input: `path: "C:\\Program Files\\*" AND message: \"Hello World\"`,
		},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()

			for i := 0; i < b.N; i++ {
				l := newLexer(tt.input)

				for {
					if l.nextToken(); l.eof() {
						break
					}
				}
			}
		})
	}
}
