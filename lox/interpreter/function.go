package interpreter

import (
	"golox/lox/ast"
	"golox/lox/environment"
)

type LoxFunction struct {
	Declaration *ast.FunctionExpr
	Name        string
	Closure     *environment.Environment
}

func NewLoxFunction(declaration *ast.FunctionExpr, closure *environment.Environment, name string) *LoxFunction {
	return &LoxFunction{Declaration: declaration, Closure: closure, Name: name}
}

func (t *LoxFunction) Call(interpreter *Interpreter, arguments []interface{}) (interface{}, error) {
	localEnvironment := environment.GetEnclosingEnvironment(t.Closure)
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
	if t.Name == "" {
		return "<fn anonymous>"
	}
	return "<fn " + t.Name + ">"
}

type FuncReturn struct {
	Value interface{}
}

func (r FuncReturn) Error() string { return "" }
