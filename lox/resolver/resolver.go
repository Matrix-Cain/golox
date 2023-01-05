package resolver

import (
	"errors"
	"golox/lox/ast"
	"golox/lox/interpreter"
	"golox/lox/lexer"
	"golox/utils"
)

type functionType int

const (
	NONE = iota
	FUNCTION
)

type Resolver struct {
	interpreter     *interpreter.Interpreter
	scopes          []map[string]bool
	currentFunction functionType
}

func NewResolver(interpreter *interpreter.Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter, scopes: make([]map[string]bool, 0), currentFunction: NONE}
}

func (i *Resolver) Resolve(obj interface{}) (interface{}, error) {
	switch v := obj.(type) {
	case []ast.Stmt:
		for _, statement := range v {
			_, err := i.Resolve(statement)
			if err != nil {
				return nil, err
			}
		}
		return nil, nil
	case ast.Stmt:
		return v.Accept(i)
	case ast.Expr:
		return v.Accept(i)
	}
	utils.Report(-1, "Resolver#Resolve", "Impossible code reached")
	return nil, errors.New("impossible code reached")
}

func (i *Resolver) VisitExpressionStmt(stmt *ast.Expression) (interface{}, error) {
	_, err := i.Resolve(stmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitIfStmt(stmt *ast.If) (interface{}, error) {
	var err error
	_, err = i.Resolve(stmt.Condition)
	if err != nil {
		return nil, err
	}
	_, err = i.Resolve(stmt.ThenBranch)
	if err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		_, err = i.Resolve(stmt.ThenBranch)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Resolver) VisitPrintStmt(stmt *ast.Print) (interface{}, error) {
	_, err := i.Resolve(stmt.Expression)
	return nil, err
}

func (i *Resolver) VisitReturnStmt(stmt *ast.Return) (interface{}, error) {
	if i.currentFunction == NONE {
		utils.Report(stmt.KeyWord.Line, "Resolver#VisitReturnStmt ", "`return` in top code")
	}
	if stmt.Value != nil {
		_, err := i.Resolve(stmt.Value)
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitWhileStmt(stmt *ast.While) (interface{}, error) {
	_, err := i.Resolve(stmt.Condition)
	if err != nil {
		return nil, err
	}

	if stmt.OptionalMutate != nil { // optionalMutate is not guaranteed to be present
		_, err = i.Resolve(stmt.OptionalMutate)
		if err != nil {
			return nil, err
		}
	}
	_, err = i.Resolve(stmt.Body)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitBreakStmt(_ *ast.Break) (interface{}, error) {
	return nil, nil
}

func (i *Resolver) VisitContinueStmt(_ *ast.Continue) (interface{}, error) {
	return nil, nil
}

func (i *Resolver) VisitBinaryExpr(expr *ast.Binary) (interface{}, error) {
	_, err := i.Resolve(expr.Left)
	if err != nil {
		return nil, err
	}
	_, err = i.Resolve(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitCallExpr(expr *ast.Call) (interface{}, error) {
	_, err := i.Resolve(expr.Callee)
	if err != nil {
		return nil, err
	}
	for _, argument := range expr.Arguments {
		_, err = i.Resolve(argument)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Resolver) VisitGroupingExpr(expr *ast.Grouping) (interface{}, error) {
	_, err := i.Resolve(expr.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitLiteralExpr(_ *ast.Literal) (interface{}, error) {
	return nil, nil
}

func (i *Resolver) VisitLogicalExpr(expr *ast.Logical) (interface{}, error) {
	_, err := i.Resolve(expr.Left)
	if err != nil {
		return nil, err
	}
	_, err = i.Resolve(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitUnaryExpr(expr *ast.Unary) (interface{}, error) {
	_, err := i.Resolve(expr.Right)
	return nil, err
}

func (i *Resolver) VisitAssignExpr(expr *ast.Assign) (interface{}, error) {
	_, err := i.Resolve(expr.Value)
	if err != nil {
		return nil, err
	}
	_, err = i.resolveLocal(expr, expr.Name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitTernaryExpr(expr *ast.Ternary) (interface{}, error) {
	_, err := i.Resolve(expr.ConditionalExpr)
	if err != nil {
		return nil, err
	}
	_, err = i.Resolve(expr.ElseExpr)
	if err != nil {
		return nil, err
	}
	if expr.ThenExpr != nil {
		_, err = i.Resolve(expr.ThenExpr)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Resolver) VisitFunctionExpr(expr *ast.FunctionExpr) (interface{}, error) {
	i.beginScope()
	for _, param := range expr.Params {
		i.declare(param)
		i.define(param)
	}
	_, err := i.Resolve(expr.Body)
	if err != nil {
		return nil, err
	}
	i.endScope()
	return nil, nil
}

func (i *Resolver) VisitBlockStmt(stmt *ast.Block) (interface{}, error) {
	i.beginScope()
	_, err := i.Resolve(stmt.Statements)
	if err != nil {
		return nil, err
	}
	i.endScope()
	return nil, nil
}

func (i *Resolver) VisitVarStmt(stmt *ast.Var) (interface{}, error) {
	i.declare(stmt.Name)
	if stmt.Initializer != nil {
		_, err := i.Resolve(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.define(stmt.Name)
	return nil, nil
}

func (i *Resolver) VisitVariableExpr(expr *ast.Variable) (interface{}, error) {
	if len(i.scopes) > 0 {
		scope := i.scopes[len(i.scopes)-1]
		if v, ok := scope[expr.Name.Lexeme]; ok {
			if !v {
				utils.Report(expr.Name.Line, "Resolver", "Can't read local variable in its own initializer")
			}
		}
	}
	_, err := i.resolveLocal(expr, expr.Name)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitFunctionStmt(stmt *ast.Function) (interface{}, error) {
	i.declare(stmt.Name)
	i.define(stmt.Name)
	_, err := i.resolveFunction(stmt, FUNCTION)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) resolveFunction(function *ast.Function, functionType0 functionType) (interface{}, error) {
	enclosingFunction := i.currentFunction
	i.currentFunction = functionType0
	i.beginScope()
	for _, param := range function.Params {
		i.declare(param)
		i.define(param)
	}
	_, err := i.Resolve(function.Body)
	if err != nil {
		return nil, err
	}
	i.endScope()
	i.currentFunction = enclosingFunction
	return nil, nil
}

func (i *Resolver) resolveLocal(expr ast.Expr, name lexer.Token) (interface{}, error) {
	for index := len(i.scopes) - 1; index >= 0; index-- {
		if _, ok := i.scopes[index][name.Lexeme]; ok {
			i.interpreter.Resolve(expr, len(i.scopes)-1-index)
			return nil, nil
		}
	}
	return nil, nil
}

func (i *Resolver) beginScope() {
	i.scopes = append(i.scopes, make(map[string]bool, 0))
}

func (i *Resolver) endScope() {
	i.scopes = i.scopes[:len(i.scopes)-1]
}

func (i *Resolver) declare(name lexer.Token) {
	if len(i.scopes) == 0 {
		return
	}
	scope := i.scopes[len(i.scopes)-1]
	if _, ok := scope[name.Lexeme]; ok {
		utils.Report(name.Line, "Resolver#declare ", "Multiple definition")
	}
	scope[name.Lexeme] = false
	i.scopes[len(i.scopes)-1] = scope

}

func (i *Resolver) define(name lexer.Token) {
	if len(i.scopes) == 0 {
		return
	}
	scope := i.scopes[len(i.scopes)-1]
	scope[name.Lexeme] = true
	i.scopes[len(i.scopes)-1] = scope
}
