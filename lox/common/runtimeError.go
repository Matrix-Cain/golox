package common

import "golox/lox/lexer"

type RuntimeError struct {
	HasError bool
	Token    lexer.Token
	Reason   string
}

func (r RuntimeError) Error() string { return "" }
