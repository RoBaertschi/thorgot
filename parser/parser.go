package parser

import (
	"fmt"
	"strconv"

	"git.robaertschi.xyz/robaertschi/thorgot/ast"
	"git.robaertschi.xyz/robaertschi/thorgot/lexer"
	"git.robaertschi.xyz/robaertschi/thorgot/token"
)

const (
	_ int = iota
	Lowest
	// Equals
	// LessGreater
	// Sum
	// Product
	// Prefix
	// Call
	// Index
)

var precedences = map[token.TokenType]int{}

type prefixFunction func() ast.ExpressionNode
type infixFunction func(ast.ExpressionNode) ast.ExpressionNode

type Parser struct {
	lexer     lexer.Lexer
	curToken  token.Token
	peekToken token.Token

	prefixFunctions map[token.TokenType]prefixFunction
	infixFunctions  map[token.TokenType]infixFunction

	Errors []error
}

func New(lexer lexer.Lexer) Parser {
	p := Parser{}

	p.prefixFunctions = make(map[token.TokenType]prefixFunction)
	p.registerPrefix(token.Integer, p.parseIntegerLiteral)

	p.infixFunctions = make(map[token.TokenType]infixFunction)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixFunction) {
	p.prefixFunctions[tokenType] = fn
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
	err := fmt.Errorf("Expected token %q to be %q", t, p.peekToken.Type)
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

func (p *Parser) parseImplicitVariableDefiniton() *ast.ImplicitVariableDefiniton {
	iv := &ast.ImplicitVariableDefiniton{Token: p.curToken}
	iv.Name = p.curToken.Literal

	if !p.expectPeek(token.Colon) {
		return nil
	}

	if !p.expectPeek(token.Equal) {
		return nil
	}

	// move onto the expression
	p.nextToken()

	iv.Value = p.parseExpression(Lowest)

	if iv.Value == nil {
		return nil
	}

	if !p.peekTokenIs(token.NewLine) && !p.peekTokenIs(token.Semicolon) {
		p.error(fmt.Errorf("variable definiton expected either an new line or an semicolon to end the definiton"))
		return nil
	}

	p.nextToken()
	p.nextToken()

	return iv
}

// Expressions

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return Lowest
}

func (p *Parser) noPrefixFunction(tokenType token.TokenType) {
	p.error(fmt.Errorf("could not find prefix expression function for token %q", tokenType))
}

func (p *Parser) parseExpression(precedence int) ast.ExpressionNode {
	prefix := p.prefixFunctions[p.curToken.Type]
	if prefix == nil {
		p.noPrefixFunction(p.curToken.Type)
		return nil
	}
	leftExpr := prefix()

	for !p.peekTokenIs(token.Semicolon) && precedence < p.peekPrecedence() {
		infix := p.infixFunctions[p.peekToken.Type]
		if infix == nil {
			return leftExpr
		}

		p.nextToken()

		leftExpr = infix(leftExpr)
	}

	return leftExpr
}

func (p *Parser) parseIntegerLiteral() ast.ExpressionNode {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		p.error(fmt.Errorf("could not parse %q as integer", p.curToken.Literal))
		return nil
	}

	lit.Value = value

	return lit
}
