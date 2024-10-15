package kql

import (
	"errors"
	"fmt"
	"strings"

	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/token"
)

// Parser is an interface capability that needs to be provided externally
// when implementing a KQL(kibana query language) parser.
type Parser interface {
	// Stmt parses a KQL(kibana query language) expression(AST).
	Stmt() (ast.Expr, error)
}

// Error is an error that occurs when parsing a KQL(kibana query language) expression, which carries the context.
type Error struct {
	s              string
	lastTokenKind  token.Kind
	lastTokenValue string
	pos            int
	err            error
}

// NewError creates a new kql error.
func NewError(s string, lastTokenKind token.Kind, lastTokenValue string, pos int, err error) error {
	if err == nil {
		return nil
	}

	e := &Error{}
	if errors.As(err, &e) {
		return err
	}

	return &Error{s, lastTokenKind, lastTokenValue, pos, err}
}

// Error returns the error message.
func (e *Error) Error() string {
	var (
		lineNo, column int
		buf            strings.Builder
	)

	if e.pos > len(e.s) {
		return e.err.Error()
	}

	for i := 0; i < e.pos; i++ {
		if e.s[i] == '\n' {
			lineNo++
			column = 0
		} else {
			column++
		}
	}

	buf.WriteString(fmt.Sprintf("line %d:%d %s\n", lineNo, column, e.err.Error()))

	lines := strings.Split(e.s, "\n")
	for i, line := range lines {
		if i == lineNo {
			buf.WriteString(line)
			buf.WriteByte('\n')

			for j := 0; j < column; j++ {
				buf.WriteByte(' ')
			}

			if e.lastTokenKind > 0 {
				buf.WriteString(strings.Repeat("^", len(e.lastTokenValue)))
			} else {
				buf.WriteString("^")
			}

			buf.WriteByte('\n')
		}
	}

	return buf.String()
}
