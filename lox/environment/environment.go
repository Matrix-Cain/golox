package environment

import (
	"golox/lox/common"
	"golox/lox/lexer"
)

type Environment struct {
	values    map[string]interface{}
	enclosing *Environment
}

func GetEnvironment() *Environment {
	return &Environment{values: make(map[string]interface{}, 0)}
}

func GetEnclosingEnvironment(nested *Environment) *Environment {
	return &Environment{values: make(map[string]interface{}, 0), enclosing: nested}
}

func (e *Environment) Get(name lexer.Token) (interface{}, error) {
	val, ok := e.values[name.Lexeme]
	if ok {
		return val, nil
	}
	if e.enclosing != nil {
		return e.enclosing.Get(name)
	}
	return nil, common.RuntimeError{HasError: true, Token: name, Reason: "Undefined variable '" + name.Lexeme + "'"}
}

func (e *Environment) GetAt(distance int, name string) interface{} {
	return e.ancestor(distance)[name]
}

func (e *Environment) Define(name string, value interface{}) {
	e.values[name] = value
}

func (e *Environment) Assign(name lexer.Token, value interface{}) error {
	_, ok := e.values[name.Lexeme]
	if ok {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.Assign(name, value)
	}
	return common.RuntimeError{HasError: true, Token: name, Reason: "Undefined variable '" + name.Lexeme + "'"}
}
