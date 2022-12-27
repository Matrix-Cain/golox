package ast

import "golox/lox/lexer"

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitExpressionStmt(stmt *Expression) (interface{}, error)
	VisitPrintStmt(stmt *Print) (interface{}, error)
	VisitVarStmt(stmt *Var) (interface{}, error)
	VisitBlockStmt(stmt *Block) (interface{}, error)
}

type Block struct {
	Statements []Stmt
}

func (t *Block) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBlockStmt(t)
}

type Expression struct {
	Expression Expr
}

func (t *Expression) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitExpressionStmt(t)
}

type Print struct {
	Expression Expr
}

func (t *Print) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitPrintStmt(t)
}

type Var struct {
	Name        lexer.Token
	Initializer Expr
}

func (t *Var) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitVarStmt(t)
}
