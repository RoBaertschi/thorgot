package ast

import (
	"strings"

	"git.robaertschi.xyz/robaertschi/thorgot/token"
)

// Statements should start with the specified Indentation, Expression should only do that on new lines
type Indentation int

func (i Indentation) indent() string {
	return strings.Repeat(" ", int(i*4))
}

type Node interface {
	TokenLiteral() string
	String(Indentation) string
}

type ExpressionNode interface {
	Node
	expressionNode()
}

type StatementNode interface {
	Node
	statementNode()
}

type Type string

func (t Type) String() string {
	return string(t)
}

type Block struct {
	Token      token.Token // the RBrace token
	Statements []StatementNode
}

func (b *Block) TokenLiteral() string { return b.Token.Literal }
func (b *Block) String(i Indentation) string {
	var out strings.Builder

	ind := i.indent()

	out.WriteString(ind + "{\n")
	for _, statement := range b.Statements {
		out.WriteString(statement.String(i + 1))
	}
	out.WriteString(ind + "}\n")

	return out.String()
}
func (b *Block) statementNode() {}

type FunctionArgument struct {
	Name string
	Type Type
}

type Function struct {
	Token         token.Token // the Fn token
	Name          string
	Arguments     []FunctionArgument
	ReturnType    Type
	HasReturnType bool
	Block         Block
}

func (f *Function) TokenLiteral() string { return f.Token.Literal }
func (f *Function) String(i Indentation) string {
	var out strings.Builder

	ind := i.indent()
	out.WriteString(ind + "fn " + f.Name + "(")
	for i, arg := range f.Arguments {
		out.WriteString(arg.Name + " " + arg.Type.String())
		if i != len(f.Arguments)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")

	if f.HasReturnType {
		out.WriteString(f.ReturnType.String() + " ")
	}

	out.WriteString(f.Block.String(i))

	return out.String()
}
func (f *Function) statementNode() {}

type ImplicitVariableDefiniton struct {
	Token token.Token // The Identifier token
	Name  string
	Value ExpressionNode
}

func (ivd *ImplicitVariableDefiniton) TokenLiteral() string {
	return ivd.Token.Literal
}
func (ivd *ImplicitVariableDefiniton) String(i Indentation) string {
	var out strings.Builder

	out.WriteString(i.indent() + ivd.Name + " := ")
	out.WriteString(ivd.Value.String(i))
	out.WriteString("\n")

	return out.String()
}

func (ivd *ImplicitVariableDefiniton) statementNode() {}

type IntegerLiteral struct {
	Token token.Token
	Value int64
}

func (il *IntegerLiteral) TokenLiteral() string        { return il.Token.Literal }
func (il *IntegerLiteral) String(i Indentation) string { return il.Token.Literal }
func (il *IntegerLiteral) expressionNode()             {}
