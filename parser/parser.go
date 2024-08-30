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
	return p.parseWrapExpr(0)
}

func (p *defaultParser) parseWrapExpr(layers int) (Expr, error) {
	oldCurrent := p.lexer.current

	token, err := p.lexer.peekToken()
	if err != nil {
		return nil, err
	}

	if token.Kind == TokenKindLparen {
		return p.parseWrapExpr(layers + 1)
	}

	p.lexer.current = oldCurrent // rollback current index

	expr, err := p.parseCombineExpr(nil)
	if err != nil {
		return nil, err
	}

	if layers == 0 {
		return expr, nil
	}

	// close wrap
	for i := 0; i < layers; i++ {
		token, err = p.lexer.peekToken()
		if err != nil {
			return nil, err
		}

		if token.Kind != TokenKindRparen {
			return nil, fmt.Errorf("expected token <Rparen>, but got %q", token.Kind.String())
		}
	}

	return p.parseCombineExpr(&WrapExpr{pos: expr.Pos() - layers, Layers: layers, Expr: expr})
}

func (p *defaultParser) parseCombineExpr(left Expr) (Expr, error) {
	switch expr := left.(type) {
	case nil:
		matchExpr, err := p.parseMatchExpr()
		if err != nil {
			return nil, err
		}

		return p.parseCombineExpr(matchExpr)
	case *MatchExpr, *WrapExpr:
		return p.parseCombineExpr(&CombineExpr{LeftExpr: expr})
	case *CombineExpr:
		if p.isEof() {
			if expr.Keyword.IsKeyword() {
				return expr, nil
			}

			return expr.LeftExpr, nil
		}

		// try peek wrap close
		if token, err := p.lexer.peekWrapper(); err == nil {
			// rollback pos
			p.lexer.current = token.Pos

			if token.Kind == TokenKindRparen {
				return expr.LeftExpr, nil
			}
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

		expr.RightExpr, err = p.parseExpr()
		if err != nil {
			return nil, err
		}

		return expr, nil
	}

	return nil, fmt.Errorf("unexpected Expr(%T)", left)
}

// need fix lint (funlen)
// //nolint: funlen
func (p *defaultParser) parseMatchExpr() (Expr, error) {
	if p.isEof() {
		return nil, errors.New("expected field or value, but got Eof")
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

	// maybe is field or value
	maybeValue := &Literal{token.Pos, token.End, token.Kind, token.Value, token.Kind == TokenKindString}
	// default operator = ":" if only value
	expr := &MatchExpr{pos: pos, HasNot: hasNot, Operator: TokenKindOperatorEql, Value: maybeValue}

	if p.isEof() {
		return expr, nil
	}

	token, err = p.lexer.peekToken()
	if err != nil {
		return nil, err
	}

	if !token.Kind.IsOperator() {
		p.lexer.current = expr.End()

		return expr, nil
	}

	expr.Field, expr.Operator = maybeValue.Value, token.Kind

	token, err = p.lexer.peekToken()
	if err != nil {
		return nil, err
	}

	// e.g. field: (...)
	if token.Kind == TokenKindLparen {
		return p.parseWrapExprInMatchExpr(token, expr)
	}

	if !token.Kind.IsValue() {
		return nil, fmt.Errorf("expected value, but got %q", token.Kind.String())
	}

	expr.Value = &Literal{token.Pos, token.End, token.Kind, token.Value, token.Kind == TokenKindString}

	return expr, nil
}

func (p *defaultParser) parseWrapExprInMatchExpr(token *Token, matchExpr *MatchExpr) (Expr, error) {
	p.lexer.current = token.Pos

	wrapExpr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	switch e := wrapExpr.(type) {
	case *CombineExpr:
		if left, ok := e.LeftExpr.(*WrapExpr); ok {
			left.pos = matchExpr.pos
			left.Field = matchExpr.Field
		}

		if right, ok := e.RightExpr.(*WrapExpr); ok {
			right.pos = matchExpr.pos
			right.Field = matchExpr.Field
		}
	case *WrapExpr:
		e.pos = matchExpr.pos
		e.Field = matchExpr.Field
	}

	return wrapExpr, nil
}

func (p *defaultParser) isEof() bool {
	return p.lexer.isEof()
}
