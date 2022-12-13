package tests

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"golox/lox/ast"
	"golox/lox/lexer"
	"testing"
)

func TestAST_Printer(t *testing.T) {
	expression := &ast.Binary{Left: &ast.Unary{Operator: lexer.Token{Type0: lexer.MINUS, Lexeme: "-", Line: 1}, Right: &ast.Literal{Value: 123}},
		Operator: lexer.Token{Type0: lexer.STAR, Lexeme: "*", Line: 1},
		Right:    &ast.Grouping{Expression: &ast.Literal{Value: 45.67}}}

	printer := ast.AstPrinter{}
	str, _ := printer.Print(expression)
	fmt.Println()
	fmt.Println(str.(string))
	fmt.Println()
	assert.Equal(t, "(* (- 123) (group 45.67))", str.(string))
}
