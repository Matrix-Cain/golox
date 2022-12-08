package main

type Token struct {
	type0   TokenType
	lexeme  string
	literal string
	line    int
}

func NewToken(type0 TokenType, lexeme string, literal string, line int) *Token {
	return &Token{
		type0:   type0,
		lexeme:  lexeme,
		literal: literal,
		line:    line,
	}
}

func (t Token) String() string {
	literalVal := t.literal
	if literalVal == "" {
		literalVal = "null"
	}
	return TokenTypeMapper[int(t.type0)] + " " + t.lexeme + " " + literalVal
}
