package parser_test

import (
	"fmt"

	"github.com/laojianzi/kql-go/parser"
)

func Example_simple() {
	query := `service_name: "redis" OR service_name: "mysql" AND NOT level: "error" and start_time > 1723286863 anD latency >= 1.5`
	// Parse query into AST
	stmt, err := parser.New(query).Stmt()
	if err != nil {
		panic(err)
	}

	// output AST to KQL(kibana query language) query
	fmt.Println(stmt.String())
	// output:
	// service_name: "redis" OR service_name: "mysql" AND NOT level: "error" AND start_time > 1723286863 AND latency >= 1.5
}

func Example_withParen() {
	query := `(service_name: "redis" OR service_name: "mysql" AND level: ("error" OR "warn") and (NOT (start_time > 1723286863 anD latency >= 1.5))) AND end_time < 1723386863`
	// Parse query into AST
	stmt, err := parser.New(query).Stmt()
	if err != nil {
		panic(err)
	}

	// output AST to KQL(kibana query language) query
	fmt.Println(stmt.String())
	// output:
	// (service_name: "redis" OR service_name: "mysql" AND level: ("error" OR "warn") AND (NOT (start_time > 1723286863 AND latency >= 1.5))) AND end_time < 1723386863
}
