package parser

import (
	"errors"
	"fmt"
	"strings"
)

type Parser interface {
	Stmt() (Expr, error)
}

type defaultParser struct {
	lexer *defaultLexer
}

func New(s string) *defaultParser { // skipcq: RVV-B0011
	return &defaultParser{lexer: &defaultLexer{input: strings.TrimSpace(s)}}
}

func (p *defaultParser) Stmt() (Expr, error) {
	return p.parseExpr()
}

func (p *defaultParser) parseExpr() (Expr, error) {
	return p.parseCombineExpr(nil)
}

func (p *defaultParser) parseCombineExpr(left Expr) (Expr, error) {
	switch expr := left.(type) {
	case nil:
		matchExpr, err := p.parseMatchExpr()
		if err != nil {
			return nil, err
		}

		return p.parseCombineExpr(matchExpr)
	case *MatchExpr:
		return p.parseCombineExpr(&CombineExpr{LeftExpr: expr})
	case *CombineExpr:
		if p.isEof() {
			return expr.LeftExpr, nil
		}

		if err := p.lexer.peekWhitespace(); err != nil {
			return nil, err
		}

		token, err := p.lexer.peekToken()
		if err != nil {
			return nil, err
		}

		if !token.Kind.IsKeyword() {
			return nil, KeywordsExpected(token.Kind.String())
		}

		expr.Keyword = token.Kind

		if err := p.lexer.peekWhitespace(); err != nil {
			return nil, err
		}

		expr.RightExpr, err = p.parseCombineExpr(nil)
		if err != nil {
			return nil, err
		}

		return expr, nil
	}

	return nil, fmt.Errorf("unexpected Expr(%T)", left)
}

func (p *defaultParser) parseMatchExpr() (Expr, error) {
	if p.isEof() {
		return nil, errors.New("expected value or match expr, but got Eof")
	}

	pos := p.lexer.current

	token, err := p.lexer.peekToken()
	if err != nil {
		return nil, err
	}

	var hasNot bool
	if token.Kind == TokenKindKeywordNot {
		hasNot = true

		if err := p.lexer.peekWhitespace(); err != nil {
			return nil, err
		}

		// get token for next step
		if token, err = p.lexer.peekToken(); err != nil {
			return nil, err
		}
	}

	if !token.Kind.IsField() && !token.Kind.IsValue() {
		return nil, fmt.Errorf("expected field or value, but got %q", token.Kind.String())
	}

	// maby is field or value
	mabyValue := &Literal{token.Pos, token.End, token.Kind, token.Value, token.Kind == TokenKindString}
	// default operator = ":" if only value
	expr := &MatchExpr{pos: pos, HasNot: hasNot, Operator: TokenKindOperatorEql, Value: mabyValue}

	if p.isEof() {
		return expr, nil
	}

	token, err = p.lexer.peekToken()
	if err != nil {
		return nil, err
	}

	if !token.Kind.IsOperator() {
		p.lexer.current = mabyValue.end

		return expr, nil
	}

	expr.Field, expr.Operator = mabyValue.Value, token.Kind

	token, err = p.lexer.peekToken()
	if err != nil {
		return nil, err
	}

	if !token.Kind.IsValue() {
		return nil, fmt.Errorf("expected value, but got %q", token.Kind.String())
	}

	expr.Value = &Literal{token.Pos, token.End, token.Kind, token.Value, token.Kind == TokenKindString}

	return expr, nil
}

func (p *defaultParser) isEof() bool {
	return p.lexer.isEof()
}
