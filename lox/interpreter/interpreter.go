package interpreter

import (
	"fmt"
	"golox/lox/ast"
	"golox/lox/common"
	"golox/lox/environment"
	"golox/lox/lexer"
	"golox/utils"
	"strconv"
	"strings"
)

type Interpreter struct {
	environment *environment.Environment
}

func NewInterpreter() *Interpreter {
	return &Interpreter{environment: environment.GetEnvironment()}
}

func (i *Interpreter) Interpret(statements []ast.Stmt) common.RuntimeError {
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			if runtimeError, ok := err.(common.RuntimeError); ok {
				utils.RaiseError(runtimeError.Token.Line, runtimeError.Reason)
				return runtimeError
			}

		}
	}
	//value, err := i.evaluate(expression)
	//if err != nil {
	//	common.runtimeError := err.(common.RuntimeError)
	//	utils.RaiseError(common.runtimeError.Token.Line, common.runtimeError.Reason)
	//	return common.runtimeError
	//}
	//fmt.Println(i.stringify(value))
	return common.RuntimeError{HasError: false}
}

func (i *Interpreter) VisitLiteralExpr(expr *ast.Literal) (interface{}, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitGroupingExpr(expr *ast.Grouping) (interface{}, error) {
	return i.evaluate(expr.Expression)
}

func (i *Interpreter) VisitUnaryExpr(expr *ast.Unary) (interface{}, error) {
	right, err := i.evaluate(expr.Right)

	if err != nil {
		return nil, err
	}

	rightVal := ast.Literal{Value: right}

	switch right.(type) {
	case float64:
		rightVal.Type = lexer.NUMBER
	case string:
		rightVal.Type = lexer.STRING
	case bool:
		if right.(bool) == true {
			rightVal.Type = lexer.TRUE
		} else {
			rightVal.Type = lexer.FALSE
		}
	}

	switch expr.Operator.Type0 {
	case lexer.BANG:
		return !i.isTruthy(right), nil

	case lexer.MINUS:
		err := i.checkNumberOperand(expr.Operator, rightVal)
		if err != nil {
			return nil, err
		}
		return -rightVal.Value.(float64), nil
	}
	return nil, common.RuntimeError{HasError: true, Token: expr.Operator, Reason: "Unexpected error: VisitUnaryExpr unreachable"}
}

func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) (interface{}, error) {
	return i.environment.Get(expr.Name)
}

func (i *Interpreter) VisitBinaryExpr(expr *ast.Binary) (interface{}, error) {
	left, err := i.evaluate(expr.Left)
	right, err := i.evaluate(expr.Right)

	if err != nil {
		return nil, err
	}
	leftVal := ast.Literal{Value: left}
	rightVal := ast.Literal{Value: right}

	switch left.(type) {
	case float64:
		leftVal.Type = lexer.NUMBER
	case string:
		leftVal.Type = lexer.STRING
	case bool:
		if left.(bool) == true {
			leftVal.Type = lexer.TRUE
		} else {
			leftVal.Type = lexer.FALSE
		}
	}

	switch right.(type) {
	case float64:
		rightVal.Type = lexer.NUMBER
	case string:
		rightVal.Type = lexer.STRING
	case bool:
		if right.(bool) == true {
			rightVal.Type = lexer.TRUE
		} else {
			rightVal.Type = lexer.FALSE
		}
	}

	switch expr.Operator.Type0 {
	case lexer.GREATER:
		err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.Value.(float64) > rightVal.Value.(float64), nil
	case lexer.GREATER_EQUAL:
		err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.Value.(float64) >= rightVal.Value.(float64), nil
	case lexer.LESS:
		err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.Value.(float64) < rightVal.Value.(float64), nil
	case lexer.LESS_EQUAL:
		err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.Value.(float64) <= rightVal.Value.(float64), nil
	case lexer.BANG_EQUAL:
		return !i.isEqual(leftVal, rightVal), nil
	case lexer.EQUAL_EQUAL:
		return i.isEqual(leftVal, rightVal), nil
	case lexer.MINUS:
		err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.Value.(float64) - rightVal.Value.(float64), nil
	case lexer.PLUS:
		if leftVal.Type == lexer.STRING && rightVal.Type == lexer.NUMBER {
			return leftVal.Value.(string) + strconv.FormatFloat(rightVal.Value.(float64), 'f', -1, 64), nil
		}

		if leftVal.Type == lexer.NUMBER && rightVal.Type == lexer.STRING {
			return strconv.FormatFloat(leftVal.Value.(float64), 'f', -1, 64) + rightVal.Value.(string), nil
		}

		if leftVal.Type == lexer.STRING && rightVal.Type == lexer.STRING {
			return leftVal.Value.(string) + rightVal.Value.(string), nil
		}

		if leftVal.Type == lexer.NUMBER && rightVal.Type == lexer.NUMBER {
			return leftVal.Value.(float64) + rightVal.Value.(float64), nil
		}

		return nil, common.RuntimeError{HasError: true, Token: expr.Operator, Reason: "operands must be numbers or strings"}
	case lexer.SLASH:
		err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.Value.(float64) / rightVal.Value.(float64), nil
	case lexer.STAR:
		err := i.checkNumberOperands(expr.Operator, leftVal, rightVal)
		if err != nil {
			return nil, err
		}
		return leftVal.Value.(float64) * rightVal.Value.(float64), nil
	}

	return nil, common.RuntimeError{HasError: true, Token: expr.Operator, Reason: "Unexpected error: VisitBinaryExpr unreachable"}
}

func (i *Interpreter) VisitTernaryExpr(expr *ast.Ternary) (interface{}, error) {
	return nil, nil
}

func (i *Interpreter) evaluate(expr ast.Expr) (interface{}, error) {
	return expr.Accept(i)
}

func (i *Interpreter) execute(stmt ast.Stmt) (interface{}, error) {
	return stmt.Accept(i)
}

func (i *Interpreter) executeBlock(statements []ast.Stmt, environment *environment.Environment) (interface{}, error) {
	previousEnviron := i.environment
	defer func() {
		i.environment = previousEnviron
	}()
	i.environment = environment
	for _, statement := range statements {
		_, err := i.execute(statement)
		if err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func (i *Interpreter) VisitBlockStmt(stmt *ast.Block) (interface{}, error) {
	_, err := i.executeBlock(stmt.Statements, environment.GetEnclosingEnvironment(i.environment)) // introduce enclosing var scope
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitExpressionStmt(stmt *ast.Expression) (interface{}, error) {
	_, err := i.evaluate(stmt.Expression)
	if err != nil {
		return nil, err
	}
	return nil, nil
}

func (i *Interpreter) VisitPrintStmt(stmt *ast.Print) (interface{}, error) {
	value, err := i.evaluate(stmt.Expression)
	if err == nil {
		fmt.Println(i.stringify(value))
	}
	return nil, err
}

func (i *Interpreter) VisitVarStmt(stmt *ast.Var) (interface{}, error) {
	var value interface{}
	var err error
	if stmt.Initializer != nil {
		value, err = i.evaluate(stmt.Initializer)
		if err != nil {
			return nil, err
		}
	}
	i.environment.Define(stmt.Name.Lexeme, value)
	return nil, nil
}

func (i *Interpreter) VisitAssignExpr(expr *ast.Assign) (interface{}, error) {
	value, err := i.evaluate(expr.Value)
	if err == nil {
		err := i.environment.Assign(expr.Name, value)
		if err != nil {
			return nil, err
		}
	}
	return nil, err

}

/*
isTruthy follow ruby rule of judging true and false

	false/nil -> false
	others -> true
*/
func (i *Interpreter) isTruthy(object interface{}) bool {
	if object == nil {
		return false
	}
	if v, ok := object.(bool); ok {
		return v
	}
	return true
}

func (i *Interpreter) isEqual(objectA ast.Literal, objectB ast.Literal) bool {
	if objectA.Type == objectB.Type && objectA.Value == objectB.Value {
		return true
	} else {
		return false
	}
}

func (i *Interpreter) checkNumberOperand(operator lexer.Token, operand ast.Literal) error {
	if operand.Type == lexer.NUMBER {
		return nil
	}
	return common.RuntimeError{HasError: true, Token: operator, Reason: "operand must be a number"}
}

func (i *Interpreter) checkNumberOperands(operator lexer.Token, operandLeft ast.Literal, operandRight ast.Literal) error {
	if operandLeft.Type == lexer.NUMBER && operandRight.Type == lexer.NUMBER {
		return nil
	}
	return common.RuntimeError{HasError: true, Token: operator, Reason: "operand must be a number"}
}

func (i *Interpreter) stringify(object interface{}) string {
	if object == nil {
		return "nil"
	}

	if v, ok := object.(float64); ok {
		text := strconv.FormatFloat(v, 'f', -1, 64)
		if strings.HasSuffix(text, ".0") {
			text = text[0 : len(text)-2]
		}
		return text
	}
	if v, ok := object.(string); ok {
		return v
	}
	return ""
}
