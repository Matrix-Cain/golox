package lexer

import "strconv"

type Token struct {
	Type0   TokenType
	Lexeme  string
	Literal interface{}
	Line    int
}

func NewToken(type0 TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type0:   type0,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	if t.Type0 == NUMBER {
		return TokenTypeMapper[int(t.Type0)] + " " + t.Lexeme + " " + strconv.FormatFloat(t.Literal.(float64), 'f', -1, 64)
	} else {
		return TokenTypeMapper[int(t.Type0)] + " " + t.Lexeme + " " + t.Literal.(string)
	}

}
