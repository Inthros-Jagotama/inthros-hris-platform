package calculator

import (
	"fmt"
	"strconv"
	"strings"
)

// Parser mengubah stream token menjadi AST. Grammar (precedence naik):
//
//	expression := term (('+' | '-') term)*
//	term       := unary (('*' | '/') unary)*
//	unary      := ('-' | '+')* postfix
//	postfix    := primary ('%')*
//	primary    := NUMBER | IDENT | '(' expression ')'
type Parser struct {
	tokens []Token
	pos    int
}

// Parse meng-parse sebuah string formula menjadi AST.
func Parse(input string) (Node, error) {
	tokens, err := Tokenize(input)
	if err != nil {
		return nil, err
	}
	p := &Parser{tokens: tokens}
	node, err := p.parseExpression()
	if err != nil {
		return nil, err
	}
	if p.current().Type != TokenEOF {
		return nil, fmt.Errorf("ekspresi tidak valid: token tak terduga %q pada posisi %d", p.current().Value, p.current().Pos)
	}
	return node, nil
}

func (p *Parser) current() Token  { return p.tokens[p.pos] }
func (p *Parser) advance() Token  { t := p.tokens[p.pos]; p.pos++; return t }

func (p *Parser) expect(t TokenType) (Token, error) {
	if p.current().Type != t {
		return Token{}, fmt.Errorf("ekspresi tidak valid: diharapkan %s, ditemukan %q pada posisi %d", t, p.current().Value, p.current().Pos)
	}
	return p.advance(), nil
}

func (p *Parser) parseExpression() (Node, error) {
	left, err := p.parseTerm()
	if err != nil {
		return nil, err
	}
	for p.current().Type == TokenPlus || p.current().Type == TokenMinus {
		op := p.advance().Type
		right, err := p.parseTerm()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseTerm() (Node, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for p.current().Type == TokenStar || p.current().Type == TokenSlash {
		op := p.advance().Type
		right, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		left = &BinaryOpNode{Op: op, Left: left, Right: right}
	}
	return left, nil
}

func (p *Parser) parseUnary() (Node, error) {
	if p.current().Type == TokenMinus || p.current().Type == TokenPlus {
		op := p.advance().Type
		operand, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return &UnaryNode{Op: op, Operand: operand}, nil
	}
	return p.parsePostfix()
}

func (p *Parser) parsePostfix() (Node, error) {
	primary, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	for p.current().Type == TokenPercent {
		p.advance()
		primary = &PercentNode{Operand: primary}
	}
	return primary, nil
}

func (p *Parser) parsePrimary() (Node, error) {
	tok := p.current()
	switch tok.Type {
	case TokenNumber:
		p.advance()
		value, err := strconv.ParseFloat(strings.ReplaceAll(tok.Value, "_", ""), 64)
		if err != nil {
			return nil, fmt.Errorf("angka tidak valid %q pada posisi %d", tok.Value, tok.Pos)
		}
		return &NumberNode{Value: value}, nil
	case TokenIdent:
		p.advance()
		return &VariableNode{Name: tok.Value}, nil
	case TokenLParen:
		p.advance()
		inner, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if _, err := p.expect(TokenRParen); err != nil {
			return nil, err
		}
		return inner, nil
	default:
		return nil, fmt.Errorf("ekspresi tidak valid: token tak terduga %q pada posisi %d", tok.Value, tok.Pos)
	}
}
