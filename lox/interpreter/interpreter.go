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
	global        *environment.Environment
	environment   *environment.Environment
	locals        map[ast.Expr]int
	loopCnt       int
	breakState    bool
	continueState bool
}

func NewInterpreter() *Interpreter {
	interpreter := &Interpreter{global: environment.GetEnvironment()}
	interpreter.environment = interpreter.global

	interpreter.global.Define("clock", &Clock{})
	return interpreter
}

func (i *Interpreter) Resolve(expr ast.Expr, depth int) {
	i.locals[expr] = depth
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

func (i *Interpreter) lookUpVariable(name lexer.Token, expr ast.Expr) (interface{}, error) {

	if distance, ok := i.locals[expr]; ok {
		return i.environment.GetAt(distance, name.Lexeme)
	} else {
		return i.global.Get(name)
	}
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

	case lexer.INCREMENT:
		err := i.checkNumberOperand(expr.Operator, rightVal)
		if err != nil {
			return nil, err
		}
		return rightVal.Value.(float64) + 1, nil
	case lexer.DECREMENT:
		err := i.checkNumberOperand(expr.Operator, rightVal)
		if err != nil {
			return nil, err
		}
		return rightVal.Value.(float64) - 1, nil
	}
	return nil, common.RuntimeError{HasError: true, Token: expr.Operator, Reason: "Unexpected error: VisitUnaryExpr unreachable"}
}

func (i *Interpreter) VisitVariableExpr(expr *ast.Variable) (interface{}, error) {
	return i.lookUpVariable(expr.Name, expr)
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
	case int:
		leftVal.Type = lexer.NUMBER
	case int64:
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
	case int:
		rightVal.Type = lexer.NUMBER
	case int64:
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
		var leftFloat, rightFloat float64
		var ok bool
		if leftFloat, ok = utils.InterfaceToFloat64(leftVal.Value); ok {
		} else {
			return nil, common.RuntimeError{HasError: true, Reason: "cannot convert to float"}
		}
		if rightFloat, ok = utils.InterfaceToFloat64(rightVal.Value); ok {
		} else {
			return nil, common.RuntimeError{HasError: true, Reason: "cannot convert to float"}
		}
		return leftFloat - rightFloat, nil
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
			var leftFloat, rightFloat float64
			var ok bool
			if leftFloat, ok = utils.InterfaceToFloat64(leftVal.Value); ok {
			} else {
				return nil, common.RuntimeError{HasError: true, Reason: "cannot convert to float"}
			}
			if rightFloat, ok = utils.InterfaceToFloat64(rightVal.Value); ok {
			} else {
				return nil, common.RuntimeError{HasError: true, Reason: "cannot convert to float"}
			}
			return leftFloat + rightFloat, nil
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

func (i *Interpreter) VisitCallExpr(expr *ast.Call) (interface{}, error) {
	callee, err := i.evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	arguments := make([]interface{}, 0)
	for _, argument := range expr.Arguments {
		v, err := i.evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments = append(arguments, v)
	}
	if _, ok := callee.(LoxCallable); !ok {
		return nil, common.RuntimeError{HasError: true, Token: expr.Paren, Reason: "Can only call functions and classes"}
	}
	function := callee.(LoxCallable)
	if function.Arity() != len(arguments) {
		return nil, common.RuntimeError{HasError: true, Token: expr.Paren, Reason: "Expected " + strconv.Itoa(function.Arity()) + " arguments but got " + strconv.Itoa(len(arguments))}
	}
	return function.Call(i, arguments)
}

func (i *Interpreter) VisitTernaryExpr(expr *ast.Ternary) (interface{}, error) {
	var res interface{}
	var err error
	res, err = i.evaluate(expr.ConditionalExpr)
	if err != nil {
		return nil, err
	}
	if i.isTruthy(res) {
		_, err = i.evaluate(expr.ThenExpr)
	} else {
		if !i.isTruthy(res) {
			_, err = i.evaluate(expr.ElseExpr)
		}
	}
	return nil, err
}

func (i *Interpreter) VisitLogicalExpr(expr *ast.Logical) (interface{}, error) {
	var left interface{}
	var err error
	left, err = i.evaluate(expr.Left)
	if expr.Operator.Type0 == lexer.OR {
		if i.isTruthy(left) {
			return left, err
		}
	} else {
		if !i.isTruthy(left) {
			return left, err
		}
	}
	return i.evaluate(expr.Right)
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
			if v, ok := err.(*FuncReturn); ok {
				return v.Value, nil
			}
			return nil, err
		}

		if i.breakState {
			return nil, nil
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

func (i *Interpreter) VisitBreakStmt(stmt *ast.Break) (interface{}, error) {
	if i.loopCnt > 0 {
		i.breakState = true
	} else {
		return nil, common.RuntimeError{HasError: true, Token: lexer.Token{Type0: lexer.BREAK}, Reason: "'break' outside loop"}
	}
	return nil, nil
}

func (i *Interpreter) VisitContinueStmt(stmt *ast.Continue) (interface{}, error) {
	if i.loopCnt > 0 {
		i.continueState = true
	} else {
		return nil, common.RuntimeError{HasError: true, Token: lexer.Token{Type0: lexer.CONTINUE}, Reason: "'continue' outside loop"}
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

func (i *Interpreter) VisitFunctionStmt(stmt *ast.Function) (interface{}, error) {
	function := NewLoxFunction(&ast.FunctionExpr{Params: stmt.Params, Body: stmt.Body}, i.environment, stmt.Name.Lexeme)
	i.environment.Define(stmt.Name.Lexeme, function)
	return nil, nil
}

func (i *Interpreter) VisitFunctionExpr(expr *ast.FunctionExpr) (interface{}, error) {
	return NewLoxFunction(expr, i.environment, ""), nil
}

func (i *Interpreter) VisitIfStmt(stmt *ast.If) (interface{}, error) {
	var err error
	var res interface{}
	res, err = i.evaluate(stmt.Condition)
	if i.isTruthy(res) {
		_, err = i.execute(stmt.ThenBranch)
		if err != nil {
			return nil, err
		}
	} else if stmt.ElseBranch != nil { // Fix: elseBranch is optional
		_, err = i.execute(stmt.ElseBranch)
		if err != nil {
			return nil, err
		}
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

func (i *Interpreter) VisitReturnStmt(stmt *ast.Return) (interface{}, error) {
	var value interface{}
	var err error
	if stmt.Value != nil {
		value, err = i.evaluate(stmt.Value)
		if err != nil {
			return nil, err
		}
	}
	return nil, &FuncReturn{Value: value}
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

func (i *Interpreter) VisitWhileStmt(stmt *ast.While) (interface{}, error) {
	var result interface{}
	var err error
	i.loopCnt++
	defer func() {
		i.loopCnt--
	}()
	result, err = i.evaluate(stmt.Condition)
	if err != nil {
		return nil, err
	}
	for i.isTruthy(result) {
		_, err = i.execute(stmt.Body)

		if err != nil {
			return nil, err
		}

		if i.continueState {
			i.continueState = false
		} else if i.breakState {
			i.breakState = false
			return nil, nil
		}
		if stmt.OptionalMutate != nil { // for loop
			_, err = i.evaluate(stmt.OptionalMutate)
			if err != nil {
				return nil, err
			}
		}

		result, err = i.evaluate(stmt.Condition)
		if err != nil {
			return nil, err
		}
	}
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

func (i *Interpreter) isTruthy(object interface{}) bool {
	/*
		isTruthy follow ruby rule of judging true and false

			false/nil -> false
			others -> true
	*/
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

func (i *Interpreter) stringify(object interface{}) interface{} {
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
	if v, ok := object.(int64); ok {
		text := strconv.FormatInt(v, 10)
		return text
	}
	if v, ok := object.(int); ok {
		text := strconv.Itoa(v)
		return text
	}

	return object
}
