package lexer_test

import (
	"testing"

	"git.robaertschi.xyz/robaertschi/thorgot/lexer"
	"git.robaertschi.xyz/robaertschi/thorgot/token"
)

func TestCorrectTokens(t *testing.T) {
	tests := []struct {
		expectedTokens []token.Token
		input          string
	}{{
		expectedTokens: []token.Token{{Token: token.Eof, Literal: "", Loc: token.Loc{Line: 1, Col: 1}}},
		input:          "",
	}, {input: "hello 1234 ; () {}",
		expectedTokens: []token.Token{
			{Token: token.Identifier, Literal: "hello", Loc: token.Loc{Line: 1, Col: 1}},
			{Token: token.Integer, Literal: "1234", Loc: token.Loc{Line: 1, Col: 7}},
			{Token: token.Semicolon, Literal: ";", Loc: token.Loc{Line: 1, Col: 12}},
			{Token: token.LParen, Literal: "(", Loc: token.Loc{Line: 1, Col: 14}},
			{Token: token.RParen, Literal: ")", Loc: token.Loc{Line: 1, Col: 15}},
		}}}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		for _, expected := range test.expectedTokens {
			actual := lexer.NextToken()

			if expected.Literal != actual.Literal {
				t.Errorf("Literal is not equal: actual = (%v) is not expected = (%v)", actual.Literal, expected.Literal)
			}

			if expected.Token != actual.Token {
				t.Errorf("Token is not equal: actual = (%v) is not expected = (%v)", actual.Token, expected.Token)
			}

			if expected.Loc.Line != actual.Loc.Line {
				t.Errorf("Loc Line is not equal: actual = (%v) is not expected = (%v)", actual.Loc.Line, expected.Loc.Line)
			}

			if expected.Loc.Col != actual.Loc.Col {
				t.Errorf("Loc Col is not equal: actual = (%v) is not expected = (%v)", actual.Loc.Col, expected.Loc.Col)
			}

		}
	}
}
