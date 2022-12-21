package ast

import (
	"errors"
	"strconv"
)

type AstPrinter struct{}

func (a *AstPrinter) Print(expr Expr) (interface{}, error) {
	if expr == nil {
		return nil, errors.New("empty expressions")
	}
	return expr.Accept(a)
}

func (a *AstPrinter) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (a *AstPrinter) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return a.parenthesize("group", expr.Expression), nil
}

func (a *AstPrinter) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	switch v := expr.Value.(type) {
	case string:
		if v == "" {
			return "nil", nil
		}
	}
	return expr.Value, nil
}

func (a *AstPrinter) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (a *AstPrinter) parenthesize(name string, exprs ...Expr) string {
	tmpStr := "(" + name
	for _, expr := range exprs {
		tmpStr += " "
		str, _ := expr.Accept(a)
		switch v := str.(type) {
		case string:
			tmpStr += v
		case int:
			tmpStr += strconv.Itoa(v)
		case float64:
			tmpStr += strconv.FormatFloat(v, 'f', -1, 64)
		case float32:
			tmpStr += strconv.FormatFloat(float64(v), 'f', -1, 32)
		}

	}
	tmpStr += ")"
	return tmpStr
}
