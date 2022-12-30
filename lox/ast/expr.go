package ast

import . "golox/lox/lexer"

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type Visitor interface {
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitCallExpr(expr *Call) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitLogicalExpr(expr *Logical) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
	VisitVariableExpr(expr *Variable) (interface{}, error)
	VisitAssignExpr(expr *Assign) (interface{}, error)
	VisitTernaryExpr(expr *Ternary) (interface{}, error)
	VisitFunctionExpr(expr *FunctionExpr) (interface{}, error)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (t *Binary) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinaryExpr(t)
}

type Call struct {
	Callee    Expr
	Paren     Token
	Arguments []Expr
}

func (t *Call) Accept(v Visitor) (interface{}, error) {
	return v.VisitCallExpr(t)
}

type FunctionExpr struct {
	Params []Token
	Body   []Stmt
}

func (t *FunctionExpr) Accept(v Visitor) (interface{}, error) {
	return v.VisitFunctionExpr(t)
}

type Grouping struct {
	Expression Expr
}

func (t *Grouping) Accept(v Visitor) (interface{}, error) {
	return v.VisitGroupingExpr(t)
}

type Literal struct {
	Type  TokenType // golang is static typed language cache type to avoid unnecessary type switch cost
	Value interface{}
}

func (t *Literal) Accept(v Visitor) (interface{}, error) {
	return v.VisitLiteralExpr(t)
}

type Logical struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (t *Logical) Accept(v Visitor) (interface{}, error) {
	return v.VisitLogicalExpr(t)
}

type Unary struct {
	Operator Token
	Right    Expr
}

func (t *Unary) Accept(v Visitor) (interface{}, error) {
	return v.VisitUnaryExpr(t)
}

type Variable struct {
	Name Token
}

func (t *Variable) Accept(v Visitor) (interface{}, error) {
	return v.VisitVariableExpr(t)
}

type Assign struct {
	Name  Token
	Value Expr
}

func (t *Assign) Accept(v Visitor) (interface{}, error) {
	return v.VisitAssignExpr(t)
}

type Ternary struct {
	ConditionalExpr Expr
	ThenExpr        Expr
	ElseExpr        Expr
}

func (t *Ternary) Accept(v Visitor) (interface{}, error) {
	return v.VisitTernaryExpr(t)
}
