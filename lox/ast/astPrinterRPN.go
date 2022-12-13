package ast

import "strconv"

type AstPrinterRPN struct{}

func (a *AstPrinterRPN) Print(expr Expr) (interface{}, error) {
	return expr.Accept(a)
}

func (a *AstPrinterRPN) VisitBinaryExpr(expr *Binary) (interface{}, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right), nil
}

func (a *AstPrinterRPN) VisitGroupingExpr(expr *Grouping) (interface{}, error) {
	return a.parenthesize("group", expr.Expression), nil
}

func (a *AstPrinterRPN) VisitLiteralExpr(expr *Literal) (interface{}, error) {
	switch v := expr.Value.(type) {
	case string:
		if v == "" {
			return "nil", nil
		}
	}
	return expr.Value, nil
}

func (a *AstPrinterRPN) VisitUnaryExpr(expr *Unary) (interface{}, error) {
	return a.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (a *AstPrinterRPN) parenthesize(name string, exprs ...Expr) string {
	tmpStr := ""
	for _, expr := range exprs {
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
		if name != "group" {
			tmpStr += " " + name + " "
		}

	}
	return tmpStr
}
