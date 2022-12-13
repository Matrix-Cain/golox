package lexer

type LexerError struct {
	HasError bool
	Line     int
	Reason   string
}
