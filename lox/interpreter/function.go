package interpreter

import (
	"golox/lox/ast"
	"golox/lox/environment"
)

type LoxFunction struct {
	Declaration *ast.Function
}

func NewLoxFunction(declaration *ast.Function) *LoxFunction {
	return &LoxFunction{Declaration: declaration}
}

func (t *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	localEnvironment := environment.GetEnclosingEnvironment(interpreter.global)
	for index, argument := range t.Declaration.Params {
		localEnvironment.Define(argument.Lexeme, arguments[index])
	}
	val, err := interpreter.executeBlock(t.Declaration.Body, localEnvironment)
	return val, err
}

func (t *LoxFunction) Arity() int {
	return len(t.Declaration.Params)
}

func (t *LoxFunction) String() string {
	return "<fn " + t.Declaration.Name.Lexeme + ">"
}
