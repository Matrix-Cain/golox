package parser

import (
	"errors"
	"golox/lox/ast"
	"golox/lox/lexer"
	"golox/utils"
)

type Parser struct {
	tokens  []lexer.Token
	current int
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{tokens: tokens, current: 0}
}

func (p *Parser) Parse() ([]ast.Stmt, ParseError) {
	statements := make([]ast.Stmt, 0)
	for !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return statements, ParseError{HasError: true}
		}
		statements = append(statements, statement)
	}
	return statements, ParseError{HasError: false}
	//expr, err := p.expression()
	//if err != nil {
	//	return nil, ParseError{HasError: true}
	//} else {
	//	return expr, ParseError{HasError: false}
	//}
}

func (p *Parser) declaration() (ast.Stmt, error) {
	if p.match(lexer.VAR) {
		stmt, err := p.varDeclaration()
		if err != nil {
			p.synchronize()
			return nil, nil
		}
		return stmt, nil
	}
	stmt, err := p.statement()
	if err != nil {
		p.synchronize()
		return nil, nil
	}
	return stmt, nil
}

func (p *Parser) statement() (ast.Stmt, error) {
	if p.match(lexer.PRINT) {
		return p.printStatement()
	}
	if p.match(lexer.LEFT_BRACE) {
		statements, err := p.block()
		if err != nil {
			return nil, err
		}
		return &ast.Block{Statements: statements}, nil
	}
	return p.expressionStatement()
}

func (p *Parser) printStatement() (ast.Stmt, error) {
	value, err := p.expression()
	p.Consume(lexer.SEMICOLON, "Expect ';' after value.")
	return &ast.Print{Expression: value}, err
}

func (p *Parser) varDeclaration() (ast.Stmt, error) {
	name, err := p.Consume(lexer.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil, err
	}
	var initializer ast.Expr
	if p.match(lexer.EQUAL) {
		initializer, err = p.expression()
	}
	p.Consume(lexer.SEMICOLON, "Expect ';' after variable declaration.")
	return &ast.Var{Name: name, Initializer: initializer}, nil
}

func (p *Parser) expressionStatement() (ast.Stmt, error) {
	expr, err := p.expression()
	p.Consume(lexer.SEMICOLON, "Expect ';' after expression.")
	return &ast.Expression{Expression: expr}, err
}

func (p *Parser) block() ([]ast.Stmt, error) {
	statements := make([]ast.Stmt, 0)
	for !p.check(lexer.RIGHT_BRACE) && !p.isAtEnd() {
		statement, err := p.declaration()
		if err != nil {
			return statements, err
		}
		statements = append(statements, statement)
	}
	_, err := p.Consume(lexer.RIGHT_BRACE, "Expect '}' after block.")
	return statements, err
}

func (p *Parser) expression() (ast.Expr, error) {
	return p.assignment()
}

func (p *Parser) assignment() (ast.Expr, error) {
	expr, err := p.conditional()

	if p.match(lexer.EQUAL) {
		equals := p.previous()
		value, err := p.assignment()
		if v, ok := expr.(*ast.Variable); ok {
			name := v.Name
			return &ast.Assign{Name: name, Value: value}, err
		}
		utils.Report(equals.Line, " Parser ", "Invalid assignment target.")
	}
	return expr, err
}

func (p *Parser) conditional() (ast.Expr, error) {
	var expr, thenBranch, elseBranch ast.Expr
	var err error
	expr, err = p.equality()

	if p.match(lexer.QUESTION) {
		thenBranch, err = p.expression()
		p.Consume(lexer.COLON, "Expect ':' after then branch of conditional expression.")
		elseBranch, err = p.conditional()
		expr = &ast.Ternary{ConditionalExpr: expr, ThenExpr: thenBranch, ElseExpr: elseBranch}
	}
	return expr, err
}

func (p *Parser) equality() (ast.Expr, error) {
	var err error
	var right, expr ast.Expr

	expr, err = p.comparison()

	for p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL) {
		operator := p.previous()
		right, err = p.comparison()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) comparison() (ast.Expr, error) {
	var err error
	var right, expr ast.Expr

	expr, err = p.term()
	for p.match(lexer.GREATER, lexer.GREATER_EQUAL, lexer.LESS, lexer.LESS_EQUAL) {
		operator := p.previous()
		right, err = p.term()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) term() (ast.Expr, error) {
	var err error
	var right, expr ast.Expr

	expr, err = p.factor()

	for p.match(lexer.MINUS, lexer.PLUS) {
		operator := p.previous()
		right, err = p.factor()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) factor() (ast.Expr, error) {
	var err error
	var right, expr ast.Expr

	expr, err = p.unary()

	for p.match(lexer.SLASH, lexer.STAR) {
		operator := p.previous()
		right, err = p.unary()
		expr = &ast.Binary{Left: expr, Operator: operator, Right: right}
	}
	return expr, err
}

func (p *Parser) unary() (ast.Expr, error) {
	if p.match(lexer.BANG, lexer.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		return &ast.Unary{Operator: operator, Right: right}, err
	}
	return p.primary()
}

func (p *Parser) primary() (ast.Expr, error) {
	if p.match(lexer.FALSE) {
		return &ast.Literal{Type: lexer.FALSE, Value: false}, nil
	}
	if p.match(lexer.TRUE) {
		return &ast.Literal{Type: lexer.TRUE, Value: true}, nil
	}
	if p.match(lexer.NIL) {
		return &ast.Literal{Type: lexer.NIL, Value: nil}, nil
	}
	if p.match(lexer.NUMBER, lexer.STRING) {
		return &ast.Literal{Type: p.previous().Type0, Value: p.previous().Literal}, nil
	}
	if p.match(lexer.IDENTIFIER) {
		return &ast.Variable{Name: p.previous()}, nil
	}
	if p.match(lexer.LEFT_PAREN) {
		expr, err := p.expression()
		_, err = p.Consume(lexer.RIGHT_PAREN, "Expect ')' after expression.")
		return &ast.Grouping{Expression: expr}, err
	}
	if p.match(lexer.BANG_EQUAL, lexer.EQUAL_EQUAL, lexer.GREATER_EQUAL, lexer.GREATER, lexer.LESS, lexer.LESS_EQUAL, lexer.PLUS, lexer.SLASH, lexer.STAR) {
		return nil, p.raiseError(p.previous(), "Missing Left Hand Operand")
	}
	return nil, p.raiseError(p.peek(), "Expect expression.")
}

func (p *Parser) Consume(type0 lexer.TokenType, message string) (lexer.Token, error) {
	if p.check(type0) {
		return p.advance(), nil
	}
	return lexer.Token{}, p.raiseError(p.peek(), message)
}

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type0 == lexer.SEMICOLON {
			return
		}
		switch p.peek().Type0 {
		case lexer.CLASS:
		case lexer.FUN:
		case lexer.VAR:
		case lexer.FOR:
		case lexer.IF:
		case lexer.WHILE:
		case lexer.PRINT:
		case lexer.RETURN:
			return

		}
		p.advance()
	}
}

func (p *Parser) raiseError(token lexer.Token, message string) error {
	if token.Type0 == lexer.EOF {
		utils.Report(token.Line, " at end", message)
	} else {
		utils.Report(token.Line, " at '"+token.Lexeme+"'", message)
	}
	return errors.New("Parse Error")
}

func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, type0 := range types {
		if p.check(type0) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(type0 lexer.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type0 == type0
}

func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current += 1
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type0 == lexer.EOF
}

func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}
