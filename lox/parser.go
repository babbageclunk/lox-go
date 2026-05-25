package lox

import (
	"errors"
	"slices"
)

var ParseError = errors.New("parse error")

// class Parser {
//   private final List<Token> tokens;
//   private int current = 0;

//   Parser(List<Token> tokens) {
//     this.tokens = tokens;
//   }
// }

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

// Expr parse() {
//   try {
//     return expression();
//   } catch (ParseError error) {
//     return null;
//   }
// }

func (p *Parser) parse() (statements []Stmt, err error) {
	p.err = nil
	for !p.isAtEnd() && p.err == nil {
		statements = append(statements, p.declaration())
	}
	return statements, p.err
}

// private Expr expression() {
//   return equality();
// }

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
	case p.match(TokenIf):
		return p.ifStatement()

	case p.match(TokenPrint):
		return p.printStatement()

	case p.match(TokenLeftBrace):
		return BlockStmt{Statements: p.block()}
	default:
		return p.expressionStatement()
	}
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
	expr := p.equality()

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

// private Expr equality() {
//   Expr expr = comparison();

//   while (match(BANG_EQUAL, EQUAL_EQUAL)) {
//     Token operator = previous();
//     Expr right = comparison();
//     expr = new Expr.Binary(expr, operator, right);
//   }

//   return expr;
// }

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

// private Expr comparison() {
//   Expr expr = term();

//   while (match(GREATER, GREATER_EQUAL, LESS, LESS_EQUAL)) {
//     Token operator = previous();
//     Expr right = term();
//     expr = new Expr.Binary(expr, operator, right);
//   }

//   return expr;
// }

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

// private Expr term() {
//   Expr expr = factor();

//   while (match(MINUS, PLUS)) {
//     Token operator = previous();
//     Expr right = factor();
//     expr = new Expr.Binary(expr, operator, right);
//   }

//   return expr;
// }

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

// private Expr factor() {
//   Expr expr = unary();

//   while (match(SLASH, STAR)) {
//     Token operator = previous();
//     Expr right = unary();
//     expr = new Expr.Binary(expr, operator, right);
//   }

//   return expr;
// }

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

// private Expr unary() {
//   if (match(BANG, MINUS)) {
//     Token operator = previous();
//     Expr right = unary();
//     return new Expr.Unary(operator, right);
//   }

//   return primary();
// }

func (p *Parser) unary() Expr {
	if p.match(TokenBang, TokenMinus) {
		operator := p.previous()
		right := p.unary()
		return UnaryExpr{
			Operator: operator,
			Right:    right,
		}
	}
	return p.primary()
}

// private Expr primary() {
//   if (match(FALSE)) return new Expr.Literal(false);
//   if (match(TRUE)) return new Expr.Literal(true);
//   if (match(NIL)) return new Expr.Literal(null);

//   if (match(NUMBER, STRING)) {
//     return new Expr.Literal(previous().literal);
//   }

//   if (match(LEFT_PAREN)) {
//     Expr expr = expression();
//     consume(RIGHT_PAREN, "Expect ')' after expression.");
//     return new Expr.Grouping(expr);
//   }
//   throw error(peek(), "Expect expression.");
// }

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

// private boolean match(TokenType... types) {
//   for (TokenType type : types) {
//     if (check(type)) {
//       advance();
//       return true;
//     }
//   }

//   return false;
// }

func (p *Parser) match(types ...TokenType) bool {
	if slices.ContainsFunc(types, p.check) {
		p.advance()
		return true
	}
	return false
}

// private Token consume(TokenType type, String message) {
//   if (check(type)) return advance();

//   throw error(peek(), message);
// }

func (p *Parser) consume(tt TokenType, message string) Token {
	if p.check(tt) {
		return p.advance()
	}
	panic(p.error(p.peek(), message))
}

// private boolean check(TokenType type) {
//   if (isAtEnd()) return false;
//   return peek().type == type;
// }

func (p *Parser) check(tt TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == tt
}

// private Token advance() {
//   if (!isAtEnd()) current++;
//   return previous();
// }

func (p *Parser) advance() Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// private boolean isAtEnd() {
//   return peek().type == EOF;
// }

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == TokenEof
}

// private Token peek() {
//   return tokens.get(current);
// }

func (p *Parser) peek() Token {
	return p.tokens[p.current]
}

// private Token previous() {
//   return tokens.get(current - 1);
// }

func (p *Parser) previous() Token {
	return p.tokens[p.current-1]
}

// private ParseError error(Token token, String message) {
//   Lox.error(token, message);
//   return new ParseError();
// }

func (p *Parser) error(token Token, message string) error {
	p.err = errors.Join(ParseError, loxError(token, message))
	return p.err
}

// private void synchronize() {
//   advance();

//   while (!isAtEnd()) {
//     if (previous().type == SEMICOLON) return;

//     switch (peek().type) {
//       case CLASS:
//       case FUN:
//       case VAR:
//       case FOR:
//       case IF:
//       case WHILE:
//       case PRINT:
//       case RETURN:
//         return;
//     }

//     advance();
//   }
// }

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
