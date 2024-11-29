# KQL (Kibana Query Language) Parser

![GitHub CI](https://github.com/laojianzi/kql-go/actions/workflows/ci.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/laojianzi/kql-go)](https://goreportcard.com/report/github.com/laojianzi/kql-go)
[![LICENSE](https://img.shields.io/github/license/laojianzi/kql-go.svg)](https://github.com/laojianzi/kql-go/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://pkg.go.dev/github.com/laojianzi/kql-go)
[![DeepSource](https://app.deepsource.com/gh/laojianzi/kql-go.svg/?label=code+coverage&show_trend=false&token=BgPgeWYICSssJGgLh2UosQw7)](https://app.deepsource.com/gh/laojianzi/kql-go/)

A Kibana Query Language (KQL) parser implemented in Go.

## Features

- Escaped character handling
- Wildcard patterns
- Parentheses grouping
- AND/OR/NOT operators
- Field:value pairs
- String literals with quotes

## Installation

```bash
go get github.com/laojianzi/kql-go
```

## Quick Start

```go
package main

import (
    "fmt"

    "github.com/laojianzi/kql-go/parser"
)

func main() {
    query := `(service_name: "redis" OR service_name: "mysql") AND level: ("error" OR "warn") and start_time > 1723286863 anD latency >= 1.5`
    // Parse query into AST
    stmt, err := parser.New(query).Stmt()
    if err != nil {
        panic(err)
    }

    // output AST to KQL(kibana query language) query
    fmt.Println(stmt.String())
    // output:
    // (service_name: "redis" OR service_name: "mysql") AND level: ("error" OR "warn") AND start_time > 1723286863 AND latency >= 1.5
}
```

## Performance

Recent benchmark results:

```
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i5-10500 CPU @ 3.10GHz

BenchmarkParser/simple_field-12                           459882              2500 ns/op            1280 B/op         34 allocs/op
BenchmarkParser/numeric_comparison-12                     728577              1646 ns/op             688 B/op         19 allocs/op
BenchmarkParser/multiple_conditions-12                    211783              5966 ns/op            2385 B/op         62 allocs/op
BenchmarkParser/complex_query-12                           63580             18675 ns/op            7235 B/op        168 allocs/op
BenchmarkParser/escaped_chars-12                          108622             10926 ns/op            5416 B/op        131 allocs/op
BenchmarkParser/many_conditions-12                         35870             34985 ns/op           12454 B/op        257 allocs/op
BenchmarkParserParallel/simple_field-12                  1582999             773.8 ns/op            1280 B/op         34 allocs/op
BenchmarkParserParallel/numeric_comparison-12            2465758             468.9 ns/op             688 B/op         19 allocs/op
BenchmarkParserParallel/multiple_conditions-12            743210              1661 ns/op            2386 B/op         62 allocs/op
BenchmarkParserParallel/complex_query-12                  219790              5692 ns/op            7238 B/op        168 allocs/op
BenchmarkParserParallel/escaped_chars-12                  331581              3735 ns/op            5416 B/op        131 allocs/op
BenchmarkParserParallel/many_conditions-12                125736              9812 ns/op           12459 B/op        257 allocs/op
BenchmarkLexer/simple_field-12                            572068              1947 ns/op             832 B/op         25 allocs/op
BenchmarkLexer/numeric_comparison-12                     1000000              1082 ns/op             264 B/op         11 allocs/op
BenchmarkLexer/multiple_conditions-12                     278456              4342 ns/op            1360 B/op         42 allocs/op
BenchmarkLexer/complex_query-12                            77738             16504 ns/op            4768 B/op        119 allocs/op
BenchmarkLexer/escaped_chars-12                           129708              8450 ns/op            3768 B/op         96 allocs/op
BenchmarkLexer/many_conditions-12                          39974             29785 ns/op            8944 B/op        192 allocs/op
BenchmarkEscapeSequence/no_escape-12                      581481              2017 ns/op             720 B/op         26 allocs/op
BenchmarkEscapeSequence/single_escape-12                  487568              2400 ns/op             936 B/op         32 allocs/op
BenchmarkEscapeSequence/multiple_escapes-12               432496              2645 ns/op            1152 B/op         38 allocs/op
BenchmarkEscapeSequence/mixed_escapes-12                  129600              9215 ns/op            3672 B/op        100 allocs/op
```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

### Development Requirements

- Go 1.16 or higher
- golangci-lint for code quality checks

### Running Tests

```bash
# Run unit tests
go test -v -count=1 -race ./...

# Run benchmarks
go test -bench=. -benchmem ./...

# Run fuzz tests, but require go version >= 1.18
go test -fuzz=. ./parser/...
```

## Examples

### Basic Queries
```go
// Simple field value query
query := `status: "active"`

// Numeric comparison
query := `age >= 18`

// Multiple conditions
query := `status: "active" AND age >= 18`
```

### Advanced Queries
```go
// Complex grouping with wildcards
query := `(status: "active" OR status: "pending") AND name: "john*"`

// Escaped characters
query := `message: "Hello \"World\"" AND path: "C:\\Program Files\\*"`

// Multiple conditions with various operators
query := `status: "active" AND age >= 18 AND name: "john*" AND city: "New York"`
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

This project is inspired by:
- [github.com/AfterShip/clickhouse-sql-parser](https://github.com/AfterShip/clickhouse-sql-parser)
- [github.com/cloudspannerecosystem/memefish](https://github.com/cloudspannerecosystem/memefish)

## References

- [Kibana Query Language Documentation](https://www.elastic.co/guide/en/kibana/current/kuery-query.html)
