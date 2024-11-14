package parser

import (
	"fmt"

	"git.robaertschi.xyz/robaertschi/thorgot/ast"
	"git.robaertschi.xyz/robaertschi/thorgot/lexer"
	"git.robaertschi.xyz/robaertschi/thorgot/token"
)

type Parser struct {
	lexer     lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	Errors []error
}

func New(lexer lexer.Lexer) Parser {
	p := Parser{}

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) error(err error) ast.StatementNode {
	p.Errors = append(p.Errors, err)
	return nil
}

func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) peekError(t token.TokenType) {
	err := fmt.Errorf("Expected token %v to be %v", t, p.peekToken.Type)
	p.Errors = append(p.Errors, err)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.lexer.NextToken()
}

func (p *Parser) ParseProgram() ast.Program {
	program := ast.Program{}
	program.Statements = make([]ast.StatementNode, 0)

	for p.curToken.Type != token.Eof {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
	}

	return program
}

func (p *Parser) parseStatement() ast.StatementNode {
	switch p.curToken.Type {
	case token.Fn:
		return p.parseFunction()
	}

	return p.error(fmt.Errorf("Invalid token %v found with literal %v", p.curToken.Type, p.curToken.Literal))
}

func (p *Parser) parseFunctionArguments() []ast.FunctionArgument {
	args := make([]ast.FunctionArgument, 0)

	for p.peekTokenIs(token.Identifier) {
		p.nextToken()
		name := p.curToken.Literal

		if !p.expectPeek(token.Identifier) {
			return nil
		}

		args = append(args, ast.FunctionArgument{Name: name, Type: ast.Type(p.curToken.Literal)})

		if !p.peekTokenIs(token.Comma) {
			break
		}
		p.nextToken()
	}

	if !p.expectPeek(token.RParen) {
		return nil
	}

	return args
}

func (p *Parser) parseFunction() *ast.Function {
	f := &ast.Function{Token: p.curToken}

	if !p.expectPeek(token.Identifier) {
		return nil
	}

	f.Name = p.curToken.Literal

	if !p.expectPeek(token.LParen) {
		return nil
	}

	args := p.parseFunctionArguments()

	if args == nil {
		return nil
	}

	f.Arguments = args

	if p.peekTokenIs(token.Identifier) {
		p.nextToken()
		f.ReturnType = ast.Type(p.curToken.Literal)
		f.HasReturnType = true
	}

	if !p.expectPeek(token.LBrace) {
		return nil
	}

	// parse block

	f.Block = p.parseBlock()
	if f.Block == nil {
		return nil
	}

	return f
}

func (p *Parser) parseBlock() *ast.Block {
	b := &ast.Block{Token: p.curToken}
	// skip {
	p.nextToken()

	for p.curToken.Type != token.RBrace {
		stmt := p.parseStatement()
		if stmt == nil {
			return nil
		}
		b.Statements = append(b.Statements, stmt)
	}

	p.nextToken()

	return b
}
