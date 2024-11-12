package lexer

import "robaertschi.xyz/robaertschi/thorgot/token"

type Lexer struct {
	input   string
	ch      byte
	pos     int
	readPos int

	// Loc
	col  int
	line int
}

func New(input string) Lexer {
	lexer := Lexer{input: input}

	lexer.readChar()

	return lexer
}

func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPos]
	}

	if l.ch == '\n' {
		l.col = 0
		l.line += 1
	}

	l.pos = l.readPos
	l.readPos += 1
	l.col += 1
}

func (l *Lexer) makeToken(t token.TokenType, literal string) token.Token {
	return token.Token{Token: t, Literal: literal, Loc: token.Loc{Line: l.line, Col: l.col}}
}

func (l *Lexer) NextToken() token.Token {
	var token token.Token

	switch l.ch {
	case 0:

	}

	return token
}
