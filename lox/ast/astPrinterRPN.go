package ast

import (
	"strconv"
	"strings"
)

//将一般表达式转换成逆波兰式
//一般表达式	                    逆波兰式
//  E        （若为常数或变量）	E
//（E）	                        E的逆波兰式
// E1 op E2 （二元运算）	        E1的逆波兰式 E2的逆波兰式 op
// op E     （一元运算）	        E的逆波兰式 op

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
	if expr.Operator.Lexeme == "-" {
		return a.parenthesize("~", expr.Right), nil // handle special case where Negate and Minus conflict or say ambiguous
	}
	return a.parenthesize(expr.Operator.Lexeme, expr.Right), nil
}

func (a *AstPrinterRPN) parenthesize(name string, exprs ...Expr) string {
	tmpStr := ""
	for _, expr := range exprs {
		str, _ := expr.Accept(a)
		switch v := str.(type) {
		case string:
			tmpStr += v + " "
		case int:
			tmpStr += strconv.Itoa(v) + " "
		case float64:
			tmpStr += strconv.FormatFloat(v, 'f', -1, 64) + " "
		case float32:
			tmpStr += strconv.FormatFloat(float64(v), 'f', -1, 32) + " "
		}
	}
	if name != "group" {
		tmpStr += name
	} else {
		tmpStr = strings.TrimRight(tmpStr, " ")
	}
	return tmpStr
}
