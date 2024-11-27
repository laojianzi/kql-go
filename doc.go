/*
Package kql provides a robust parser for Kibana Query Language (KQL).

KQL is a simple yet powerful query language used in Kibana for filtering and searching data.
This package implements a complete parser that converts KQL queries into an Abstract Syntax Tree (AST).

Basic Usage:

	import "github.com/laojianzi/kql-go"

	query := "response:200 AND (method:GET OR method:POST)"
	ast, err := kql.Parse(query)
	if err != nil {
	    log.Fatal(err)
	}

Features:
  - Full KQL syntax support
  - Escaped character handling
  - Wildcard support
  - Parentheses grouping
  - AND/OR/NOT operators
  - Field:value pairs
  - String literals with quotes

Performance:
  - Efficient lexical analysis
  - Minimal memory allocations
  - Object pooling for tokens
  - Optimized string handling

Thread Safety:
The parser is designed to be thread-safe. Each Parse call creates a new parser instance,
making it safe to use across multiple goroutines.

Error Handling:
The parser provides detailed error messages with position information,
making it easy to identify and fix syntax errors in queries.

For more information about KQL syntax, visit:
https://www.elastic.co/guide/en/kibana/current/kuery-query.html
*/
package kql
