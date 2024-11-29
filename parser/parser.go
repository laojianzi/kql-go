package parser

import (
	"errors"
	"fmt"
	"strings"

	"github.com/laojianzi/kql-go"
	"github.com/laojianzi/kql-go/ast"
	"github.com/laojianzi/kql-go/token"
)

type defaultParser struct {
	lexer *defaultLexer
}

// New creates a new KQL parser.
func New(input string) kql.Parser {
	return &defaultParser{lexer: newLexer(input)}
}

// Stmt parses a statement from the input.
func (p *defaultParser) Stmt() (ast.Expr, error) {
	expr, err := p.parseStmt()
	if err != nil {
		return nil, p.toKQLError(err)
	}

	return expr, nil
}

func (p *defaultParser) parseStmt() (ast.Expr, error) {
	if strings.TrimSpace(string(p.lexer.Value)) == "" {
		return nil, errors.New("expected KQL(kibana query language) string, but got empty string")
	}

	stmt, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if p.lexer.Token.Kind != token.TokenKindEof {
		return nil, fmt.Errorf("expected <EOF>, but got %q", p.lexer.Token.Kind.String())
	}

	return stmt, nil
}

func (p *defaultParser) parseExpr() (ast.Expr, error) {
	if err := p.lexer.nextToken(); err != nil {
		return nil, err
	}

	expr, err := p.parseBinary()
	if err != nil {
		return nil, err
	}

	return p.parseCombine(expr)
}

func (p *defaultParser) parseCombine(left ast.Expr) (ast.Expr, error) {
	kind := p.lexer.Token.Kind
	if kind == token.TokenKindEof || kind == token.TokenKindRparen {
		return left, nil
	}

	if !kind.IsKeyword() && kind != token.TokenKindKeywordNot {
		return nil, token.KeywordsExpected(p.lexer.Token.Kind.String())
	}

	if err := p.lexer.nextToken(); err != nil {
		return nil, err
	}

	right, err := p.parseBinary()
	if err != nil {
		return nil, err
	}

	return p.parseCombine(&ast.CombineExpr{
		LeftExpr:  left,
		Keyword:   kind,
		RightExpr: right,
	})
}

func (p *defaultParser) parseBinary() (ast.Expr, error) {
	pos, hasNot := 0, false

	if p.lexer.Token.Kind == token.TokenKindKeywordNot {
		pos = p.lexer.Token.Pos
		hasNot = true

		if err := p.lexer.nextToken(); err != nil {
			return nil, err
		}
	}

	expr, err := p.parseLiteral()
	if err != nil {
		return nil, err
	}

	if !hasNot {
		pos = expr.Pos()
	}

	op := p.lexer.Token.Kind
	if !op.IsOperator() || !p.lexer.lastTokenKind.IsField() {
		return ast.NewBinaryExpr(pos, "", 0, expr, hasNot), nil
	}

	if err := p.lexer.nextToken(); err != nil {
		return nil, err
	}

	right, err := p.parseLiteral()
	if err != nil {
		return nil, err
	}

	// check >=, >, <=, < operator must be followed by Int|Float|Ident
	switch op {
	case token.TokenKindOperatorGeq, token.TokenKindOperatorGtr, token.TokenKindOperatorLeq, token.TokenKindOperatorLss:
		if n := strings.ReplaceAll(right.String(), token.TokenKindWildcard.String(), ""); n != "" && !token.IsNumber(n) {
			return nil, fmt.Errorf("expected number or number with wildcard, but got %q", n)
		}
	}

	return ast.NewBinaryExpr(pos, expr.String(), op, right, hasNot), nil
}

func (p *defaultParser) parseLiteral() (ast.Expr, error) {
	if p.lexer.Token.Kind == token.TokenKindLparen {
		return p.parseParen()
	}

	switch p.lexer.Token.Kind {
	case token.TokenKindInt, token.TokenKindFloat, token.TokenKindString, token.TokenKindIdent:
		return p.parseWildcard()
	}

	return nil, fmt.Errorf("unexpected token: %s", p.lexer.Token.Kind)
}

func (p *defaultParser) parseParen() (ast.Expr, error) {
	tok := p.lexer.Token

	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	if p.lexer.Token.Kind != token.TokenKindRparen {
		return nil, fmt.Errorf("expected token <Rparen>, but got %q", p.lexer.Token.Kind.String())
	}

	rparen := p.lexer.Token.End

	if err := p.lexer.nextToken(); err != nil {
		return nil, err
	}

	return ast.NewParenExpr(tok.Pos, rparen, expr), nil
}

func (p *defaultParser) parseWildcard() (ast.Expr, error) {
	kind := p.lexer.Token.Kind

	tok, err := p.expect(kind)
	if err != nil {
		return nil, err
	}

	pos, end := tok.Pos, tok.End
	if kind == token.TokenKindString { // with double quote "
		pos -= 1
		end += 1
	}

	lit := ast.NewLiteral(pos, end, kind, tok.Value, tok.EscapeIndexes)
	if kind != token.TokenKindIdent && kind != token.TokenKindString {
		return lit, nil
	}

	var indexes []int

	runes := []rune(lit.String())
	for i := range runes {
		if runes[i] == '*' && (i == 0 || runes[i-1] != '\\') { // skip escaped wildcard
			indexes = append(indexes, i)
		}
	}

	if len(indexes) == 0 { // not found wildcard
		return lit, nil
	}

	return ast.NewWildcardExpr(lit, indexes), nil
}

func (p *defaultParser) expect(kind token.Kind) (*Token, error) {
	if p.lexer.Token.Kind != kind {
		return nil, fmt.Errorf("expected token: %s, but: %s", kind, p.lexer.Token.Kind)
	}

	t := p.lexer.Token.Clone()

	if err := p.lexer.nextToken(); err != nil {
		return nil, err
	}

	return t, nil
}

func (p *defaultParser) toKQLError(err error) error {
	return kql.NewError(string(p.lexer.Value), p.lexer.lastTokenKind, p.lexer.Token.Value, p.lexer.Token.Pos, err)
}
