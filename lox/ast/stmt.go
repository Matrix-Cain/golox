package ast

import . "golox/lox/lexer"

type Stmt interface {
	Accept(v StmtVisitor) (interface{}, error)
}

type StmtVisitor interface {
	VisitBlockStmt(stmt *Block) (interface{}, error)
	VisitExpressionStmt(stmt *Expression) (interface{}, error)
	VisitFunctionStmt(stmt *Function) (interface{}, error)
	VisitIfStmt(stmt *If) (interface{}, error)
	VisitPrintStmt(stmt *Print) (interface{}, error)
	VisitReturnStmt(stmt *Return) (interface{}, error)
	VisitVarStmt(stmt *Var) (interface{}, error)
	VisitWhileStmt(stmt *While) (interface{}, error)
	VisitBreakStmt(stmt *Break) (interface{}, error)
	VisitContinueStmt(stmt *Continue) (interface{}, error)
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

type Function struct {
	Name   Token
	Params []Token
	Body   []Stmt
}

func (t *Function) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitFunctionStmt(t)
}

type If struct {
	Condition  Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (t *If) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitIfStmt(t)
}

type Print struct {
	Expression Expr
}

func (t *Print) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitPrintStmt(t)
}

type Return struct {
	KeyWord Token
	Value   Expr
}

func (t *Return) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitReturnStmt(t)
}

type Var struct {
	Name        Token
	Initializer Expr
}

func (t *Var) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitVarStmt(t)
}

type While struct {
	Condition      Expr
	Body           Stmt
	OptionalMutate Expr
}

func (t *While) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitWhileStmt(t)
}

type Break struct {
}

func (t *Break) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitBreakStmt(t)
}

type Continue struct {
}

func (t *Continue) Accept(v StmtVisitor) (interface{}, error) {
	return v.VisitContinueStmt(t)
}
