package parser_test

import (
	"fmt"

	"github.com/laojianzi/kql-go/parser"
)

func Example() {
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
}
