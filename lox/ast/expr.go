package ast

import . "golox/lox/lexer"

type Expr interface {
	Accept(v Visitor) (interface{}, error)
}

type Visitor interface {
	VisitBinaryExpr(expr *Binary) (interface{}, error)
	VisitGroupingExpr(expr *Grouping) (interface{}, error)
	VisitLiteralExpr(expr *Literal) (interface{}, error)
	VisitUnaryExpr(expr *Unary) (interface{}, error)
}

type Binary struct {
	Left     Expr
	Operator Token
	Right    Expr
}

func (t *Binary) Accept(v Visitor) (interface{}, error) {
	return v.VisitBinaryExpr(t)
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

type Unary struct {
	Operator Token
	Right    Expr
}

func (t *Unary) Accept(v Visitor) (interface{}, error) {
	return v.VisitUnaryExpr(t)
}
