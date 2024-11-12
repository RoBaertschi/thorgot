package token

type TokenType string

type Loc struct {
	Line int
	Col  int
}

type Token struct {
	Token   TokenType
	Literal string
	Loc     Loc
}

const (
	Illegal TokenType = "Illegal"
	Eof               = "Eof"

	EndLine = "EndLine"

	Semicolon = "Semicolon" // ;
	Colon     = "Colon"     // :
	Equal     = "Equal"     // =
	LBrace    = "LBrace"    // {
	RBrace    = "RBrace"    // }
	LParen    = "LParen"    // (
	RParen    = "RParen"    // )

	Identifier = "Identifier"

	// Keywords
	Fn = "Fn" // fn
)

var stringToToken = map[string]TokenType{
	"fn": Fn,
}

func LookupKeyword(literal string) TokenType {
	if token, ok := stringToToken[literal]; ok {
		return token
	}

	return Identifier
}