# KQL Parser
![GitHub CI](https://github.com/laojianzi/kql-go/actions/workflows/ci.yaml/badge.svg) [![Go Report Card](https://goreportcard.com/badge/github.com/laojianzi/kql-go)](https://goreportcard.com/report/github.com/laojianzi/kql-go) [![LICENSE](https://img.shields.io/github/license/laojianzi/kql-go.svg)](https://github.com/laojianzi/kql-go/blob/master/LICENSE) [![GoDoc](https://img.shields.io/badge/Godoc-reference-blue.svg)](https://pkg.go.dev/github.com/laojianzi/kql-go) [![DeepSource](https://app.deepsource.com/gh/laojianzi/kql-go.svg/?label=code+coverage&show_trend=false&token=BgPgeWYICSssJGgLh2UosQw7)](https://app.deepsource.com/gh/laojianzi/kql-go/)

The goal of this project is to build a KQL parser in Go with the following key features:

- Parse KQL query into AST
- output AST to KQL query

This project is inspired by [clickhouse-sql-parser](https://github.com/AfterShip/clickhouse-sql-parser) which is a SQL parser for ClickHouse in Go.

## How to use

Playground: https://go.dev/play/p/m36hkz43PQL

```Go
package main

import (
    "github.com/laojianzi/kql-go/parser"
)

query := `service_name: "redis" OR service_name: "mysql" AND level: "error" and start_time > 1723286863 anD latency >= 1.5`
// Parse query into AST
stmt, err := parser.New(query).Stmt()
if err != nil {
    panic(err)
}

// output AST to KQL query
fmt.Println(stmt.String())
// output:
// service_name: "redis" OR service_name: "mysql" AND level: "error" AND start_time > 1723286863 AND latency >= 1.5
```

## Contact us

Feel free to open an issue or discussion if you have any issues or questions.
