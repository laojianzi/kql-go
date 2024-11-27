# KQL (Kibana Query Language) Parser

![GitHub CI](https://github.com/laojianzi/kql-go/actions/workflows/ci.yaml/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/laojianzi/kql-go)](https://goreportcard.com/report/github.com/laojianzi/kql-go)
[![LICENSE](https://img.shields.io/github/license/laojianzi/kql-go.svg)](https://github.com/laojianzi/kql-go/blob/master/LICENSE)
[![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://pkg.go.dev/github.com/laojianzi/kql-go)
[![DeepSource](https://app.deepsource.com/gh/laojianzi/kql-go.svg/?label=code+coverage&show_trend=false&token=BgPgeWYICSssJGgLh2UosQw7)](https://app.deepsource.com/gh/laojianzi/kql-go/)

A high-performance Kibana Query Language (KQL) parser implemented in Go.

## Features

- Full KQL syntax support
- Escaped character handling
- Wildcard patterns
- Parentheses grouping
- AND/OR/NOT operators
- Field:value pairs
- String literals with quotes
- Detailed error messages
- Thread-safe design
- High performance with minimal allocations

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
BenchmarkParser/simple         	  749059	      1543 ns/op	     576 B/op	      15 allocs/op
BenchmarkParser/with_escape   	  653020	      1845 ns/op	     624 B/op	      19 allocs/op
BenchmarkParser/complex       	  103722	     11127 ns/op	    2848 B/op	      68 allocs/op
BenchmarkLexer/simple        	218364832	      5.514 ns/op	       0 B/op	       0 allocs/op
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

# Run fuzzing tests
go test -fuzz=. ./...
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
