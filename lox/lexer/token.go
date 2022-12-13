package lexer

type Token struct {
	Type0   TokenType
	Lexeme  string
	Literal string
	Line    int
}

func NewToken(type0 TokenType, lexeme string, literal string, line int) *Token {
	return &Token{
		Type0:   type0,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	literalVal := t.Literal
	if literalVal == "" {
		literalVal = "null"
	}
	return TokenTypeMapper[int(t.Type0)] + " " + t.Lexeme + " " + literalVal
}
