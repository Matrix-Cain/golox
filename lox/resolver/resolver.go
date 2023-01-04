package resolver

import (
	"errors"
	"golox/lox/ast"
	"golox/lox/interpreter"
	"golox/lox/lexer"
	"golox/utils"
)

type Resolver struct {
	interpreter interpreter.Interpreter
	scopes      []map[string]bool
}

func (i *Resolver) VisitExpressionStmt(stmt *ast.Expression) (interface{}, error) {
	_, err := i.resolve(stmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitIfStmt(stmt *ast.If) (interface{}, error) {
	var err error
	_, err = i.resolve(stmt.Condition)
	if err != nil {
		return nil, err
	}
	_, err = i.resolve(stmt.ThenBranch)
	if err != nil {
		return nil, err
	}
	if stmt.ElseBranch != nil {
		_, err = i.resolve(stmt.ThenBranch)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Resolver) VisitPrintStmt(stmt *ast.Print) (interface{}, error) {
	_, err := i.resolve(stmt.Expression)
	return nil, err
}

func (i *Resolver) VisitReturnStmt(stmt *ast.Return) (interface{}, error) {
	if stmt.Value != nil {
		_, err := i.resolve(stmt.Value)
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitWhileStmt(stmt *ast.While) (interface{}, error) {
	_, err := i.resolve(stmt.Condition)
	if err != nil {
		return nil, err
	}
	_, err = i.resolve(stmt.Body)
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
	_, err := i.resolve(expr.Left)
	if err != nil {
		return nil, err
	}
	_, err = i.resolve(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitCallExpr(expr *ast.Call) (interface{}, error) {
	_, err := i.resolve(expr.Callee)
	if err != nil {
		return nil, err
	}
	for _, argument := range expr.Arguments {
		_, err = i.resolve(argument)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Resolver) VisitGroupingExpr(expr *ast.Grouping) (interface{}, error) {
	_, err := i.resolve(expr.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitLiteralExpr(_ *ast.Literal) (interface{}, error) {
	return nil, nil
}

func (i *Resolver) VisitLogicalExpr(expr *ast.Logical) (interface{}, error) {
	_, err := i.resolve(expr.Left)
	if err != nil {
		return nil, err
	}
	_, err = i.resolve(expr.Right)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) VisitUnaryExpr(expr *ast.Unary) (interface{}, error) {
	_, err := i.resolve(expr.Right)
	return nil, err
}

func (i *Resolver) VisitAssignExpr(expr *ast.Assign) (interface{}, error) {
	_, err := i.resolve(expr.Value)
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
	_, err := i.resolve(expr.ConditionalExpr)
	if err != nil {
		return nil, err
	}
	_, err = i.resolve(expr.ElseExpr)
	if err != nil {
		return nil, err
	}
	if expr.ThenExpr != nil {
		_, err = i.resolve(expr.ThenExpr)
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
	_, err := i.resolve(expr.Body)
	if err != nil {
		return nil, err
	}
	i.endScope()
	return nil, nil
}

func NewResolver(interpreter interpreter.Interpreter) *Resolver {
	return &Resolver{interpreter: interpreter}
}

func (i *Resolver) VisitBlockStmt(stmt *ast.Block) (interface{}, error) {
	i.beginScope()
	_, err := i.resolve(stmt.Statements)
	if err != nil {
		return nil, err
	}
	i.endScope()
	return nil, nil
}

func (i *Resolver) VisitVarStmt(stmt *ast.Var) (interface{}, error) {
	i.declare(stmt.Name)
	if stmt.Initializer != nil {
		_, err := i.resolve(stmt.Initializer)
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
		if ok, v := scope[expr.Name.Lexeme]; ok {
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
	_, err := i.resolveFunction(stmt)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Resolver) resolve(obj interface{}) (interface{}, error) {
	switch v := obj.(type) {
	case []ast.Stmt:
		for _, statement := range v {
			_, err := i.resolve(statement)
			if err != nil {
				return nil, err
			}
		}
	case ast.Stmt:
		return v.Accept(i)
	case ast.Expr:
		return v.Accept(i)
	}
	utils.Report(-1, "Resolver#resolve", "Impossible code reached")
	return nil, errors.New("impossible code reached")
}

func (i *Resolver) resolveFunction(function *ast.Function) (interface{}, error) {
	i.beginScope()
	for _, param := range function.Params {
		i.declare(param)
		i.define(param)
	}
	_, err := i.resolve(function.Body)
	if err != nil {
		return nil, err
	}
	i.endScope()
	return nil, nil
}

func (i *Resolver) resolveLocal(expr ast.Expr, name lexer.Token) (interface{}, error) {
	for index := len(i.scopes) - 1; index >= 0; index-- {
		if ok, _ := i.scopes[index][name.Lexeme]; ok {
			i.interpreter.Resolve(expr, len(i.scopes)-1-index)
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
