package lexer

import (
	"git.robaertschi.xyz/robaertschi/thorgot/token"
)

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
	lexer.line = 1

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

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z'
}

func isValidIdentChar(ch byte) bool {
	return ch == '_' || isLetter(ch)
}

func (l *Lexer) skipWhitespace() {
	for l.ch == '\r' || l.ch == '\b' || l.ch == '\t' || l.ch == ' ' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() token.Token {
	loc := token.Loc{Line: l.line, Col: l.col}
	pos := l.pos

	l.readChar()

	for isDigit(l.ch) || isLetter(l.ch) {
		l.readChar()
	}

	t := token.LookupKeyword(l.input[pos:l.pos])

	return token.Token{Token: t, Loc: loc, Literal: l.input[pos:l.pos]}
}

func (l *Lexer) readNumber() token.Token {
	pos := l.pos
	loc := token.Loc{Line: l.line, Col: l.col}

	for isDigit(l.ch) {
		l.readChar()
	}

	return token.Token{Token: token.Integer, Loc: loc, Literal: l.input[pos:l.pos]}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	var tok token.Token
	tok.Literal = string(l.ch)

	switch l.ch {
	case '\n':
		tok.Token = token.EndLine
	case ';':
		tok.Token = token.Semicolon
	case ':':
		tok.Token = token.Colon
	case '=':
		tok.Token = token.Equal
	case '{':
		tok.Token = token.LBrace
	case '}':
		tok.Token = token.RBrace
	case '(':
		tok.Token = token.LParen
	case ')':
		tok.Token = token.RParen

	case 0:
		return l.makeToken(token.Eof, "")

	default:
		if isValidIdentChar(l.ch) {
			return l.readIdentifier()
		} else if isDigit(l.ch) {
			return l.readNumber()
		}

		tok.Token = token.Illegal
	}

	l.readChar()
	return tok
}
