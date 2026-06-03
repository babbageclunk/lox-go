package lox

import (
	"errors"
	"slices"
)

var ParseError = errors.New("parse error")

type Parser struct {
	tokens  []Token
	current int
	err     error
}

func NewParser(tokens []Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}

func (p *Parser) parse() (statements []Stmt, err error) {
	p.err = nil
	for !p.isAtEnd() && p.err == nil {
		statements = append(statements, p.declaration())
	}
	return statements, p.err
}

func (p *Parser) expression() Expr {
	return p.assignment()
}

func (p *Parser) safeExpression() (expr Expr, ok bool) {
	defer func() {
		val := recover()
		if val == nil {
			return
		}
		if err, _ := val.(error); errors.Is(err, ParseError) {
			expr = nil
			ok = false
			return
		}
		panic(val)
	}()
	return p.expression(), true
}

func (p *Parser) declaration() (result Stmt) {
	defer func() {
		val := recover()
		if val == nil {
			return
		}
		if err, _ := val.(error); errors.Is(err, ParseError) {
			p.synchronize()
			result = nil
			return
		}
		panic(val)
	}()
	if p.match(TokenVar) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) statement() Stmt {
	switch {
	case p.match(TokenBreak):
		p.consume(TokenSemicolon, "Expect ';' after 'break'.")
		return BreakStmt{}

	case p.match(TokenFor):
		return p.forStatement()

	case p.match(TokenIf):
		return p.ifStatement()

	case p.match(TokenPrint):
		return p.printStatement()

	case p.match(TokenWhile):
		return p.whileStatement()

	case p.match(TokenLeftBrace):
		return BlockStmt{Statements: p.block()}
	default:
		return p.expressionStatement()
	}
}

func (p *Parser) forStatement() Stmt {
	p.consume(TokenLeftParen, "Expect '(' after 'for'.")

	var initializer Stmt
	switch {
	case p.match(TokenSemicolon):
		initializer = nil
	case p.match(TokenVar):
		initializer = p.varDeclaration()
	default:
		initializer = p.expressionStatement()
	}

	var condition Expr
	if !p.check(TokenSemicolon) {
		condition = p.expression()
	}
	p.consume(TokenSemicolon, "Expect ';' after loop condition.")

	var increment Expr
	if !p.check(TokenRightParen) {
		increment = p.expression()
	}
	p.consume(TokenRightParen, "Expect ')' after for clauses.")
	body := p.statement()

	if increment != nil {
		body = BlockStmt{[]Stmt{body, ExpressionStmt{increment}}}
	}
	if condition == nil {
		condition = LiteralExpr{true}
	}
	body = WhileStmt{
		Condition: condition,
		Body:      body,
	}
	if initializer != nil {
		body = BlockStmt{[]Stmt{initializer, body}}
	}

	return body
}

func (p *Parser) ifStatement() Stmt {
	p.consume(TokenLeftParen, "Expect '(' after 'if'.")
	condition := p.expression()
	p.consume(TokenRightParen, "Exoect ')' after if condition.")

	thenBranch := p.statement()
	var elseBranch Stmt
	if p.match(TokenElse) {
		elseBranch = p.statement()
	}
	return IfStmt{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
	}
}

func (p *Parser) printStatement() Stmt {
	value := p.expression()
	p.consume(TokenSemicolon, "Expect ';' after value.")
	return PrintStmt{Expression: value}
}

func (p *Parser) varDeclaration() Stmt {
	name := p.consume(TokenIdentifier, "Expect variable name.")

	var initializer Expr
	if p.match(TokenEqual) {
		initializer = p.expression()
	}

	p.consume(TokenSemicolon, "Expect ';' after variable declaration.")
	return VarStmt{Name: name, Initializer: initializer}
}

func (p *Parser) whileStatement() Stmt {
	p.consume(TokenLeftParen, "Expect '(' after 'while'.")
	condition := p.expression()
	p.consume(TokenRightParen, "Expect ')' after condition.")
	body := p.statement()
	return WhileStmt{
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) block() []Stmt {
	var statements []Stmt
	for !p.check(TokenRightBrace) && !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	p.consume(TokenRightBrace, "Expect '}' after block.")
	return statements
}

func (p *Parser) expressionStatement() Stmt {
	expr := p.expression()
	p.consume(TokenSemicolon, "Expect ';' after expression.")
	return ExpressionStmt{Expression: expr}
}

func (p *Parser) assignment() Expr {
	expr := p.or()

	if p.match(TokenEqual) {
		equals := p.previous()
		value := p.assignment()

		if varExpr, ok := expr.(VariableExpr); ok {
			return AssignExpr{Name: varExpr.Name, Value: value}
		}

		p.error(equals, "Invalid assignment target.")
	}
	return expr
}

func (p *Parser) or() Expr {
	expr := p.and()

	for p.match(TokenOr) {
		operator := p.previous()
		right := p.and()
		expr = LogicalExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) and() Expr {
	expr := p.equality()

	for p.match(TokenAnd) {
		operator := p.previous()
		right := p.equality()
		expr = LogicalExpr{Left: expr, Operator: operator, Right: right}
	}

	return expr
}

func (p *Parser) equality() Expr {
	expr := p.comparison()

	for p.match(TokenBangEqual, TokenEqualEqual) {
		operator := p.previous()
		right := p.comparison()
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) comparison() Expr {
	expr := p.term()
	for p.match(TokenGreater, TokenGreaterEqual, TokenLess, TokenLessEqual) {
		operator := p.previous()
		right := p.term()
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) term() Expr {
	expr := p.factor()
	for p.match(TokenMinus, TokenPlus) {
		operator := p.previous()
		right := p.factor()
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) factor() Expr {
	expr := p.unary()
	for p.match(TokenSlash, TokenStar) {
		operator := p.previous()
		right := p.unary()
		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr
}

func (p *Parser) unary() Expr {
	if p.match(TokenBang, TokenMinus) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{
			Operator: operator,
			Right:    right,
		}
	}
	return p.call()
}

func (p *Parser) call() Expr {
	expr := p.primary()
	for {
		if p.match(TokenLeftParen) {
			expr = p.finishCall(expr)
		} else {
			break
		}
	}

	return expr
}

func (p *Parser) finishCall(callee Expr) Expr {
	var arguments []Expr
	if !p.check(TokenRightParen) {
		for {
			if len(arguments) >= 255 {
				p.error(p.peek(), "Can't have more than 255 arguments.")
			}
			arguments = append(arguments, p.expression())
			if !p.match(TokenComma) {
				break
			}
		}
	}
	paren := p.consume(TokenRightParen, "Expect ')' after arguments.")
	return CallExpr{
		Callee:    callee,
		Paren:     paren,
		Arguments: arguments,
	}
}

func (p *Parser) primary() Expr {
	switch {
	case p.match(TokenFalse):
		return LiteralExpr{Value: false}
	case p.match(TokenTrue):
		return LiteralExpr{Value: true}
	case p.match(TokenNil):
		return LiteralExpr{Value: nil}
	case p.match(TokenIdentifier):
		return VariableExpr{Name: p.previous()}
	case p.match(TokenNumber, TokenString):
		return LiteralExpr{Value: p.previous().Literal}
	case p.match(TokenLeftParen):
		expr := p.expression()
		p.consume(TokenRightParen, "Expect ')' after expression.")
		return GroupingExpr{Expression: expr}
	}
	panic(p.error(p.peek(), "Expect expression."))
}

func (p *Parser) match(types ...TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) consume(tt TokenType, message string) Token {
	if p.check(tt) {
		return p.advance()
	}
	panic(p.error(p.peek(), message))
}

func (p *Parser) check(tt TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tt
}

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TokenEof
}

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

func (p *Parser) error(token Token, message string) error {
	p.err = errors.Join(ParseError, loxError(token, message))
	return p.err
}

func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		if p.previous().Type == TokenSemicolon {
			return
		}

		switch p.peek().Type {
		case TokenClass, TokenFun, TokenVar, TokenFor,
			TokenIf, TokenWhile, TokenPrint, TokenReturn:
			return
		}

		p.advance()
	}
}
