package parser

import (
	"testing"

	"git.robaertschi.xyz/robaertschi/thorgot/ast"
	"git.robaertschi.xyz/robaertschi/thorgot/lexer"
)

func TestParseFunction(t *testing.T) {

}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors

	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %v", msg)
	}
	t.FailNow()
}

func testFunction(t *testing.T, s ast.StatementNode, name string, args []ast.FunctionArgument, returnType ast.Type, block ast.Block) {

}
