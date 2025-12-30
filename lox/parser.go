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

func (p *Parser) parse() (e Expr) {
	defer func() {
		val := recover()
		if val == nil {
			return
		}
		err, ok := val.(error)
		if !ok {
			panic(val)
		}
		if !errors.Is(err, ParseError) {
			panic(err)
		}
		e = nil
	}()
	return p.expression()
}

// private Expr expression() {
//   return equality();
// }

func (p *Parser) expression() Expr {
	return p.equality()
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
		expr = Binary{
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
		expr = Binary{
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
		expr = Binary{
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
		expr = Binary{
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
		return Unary{
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
		return Literal{Value: false}
	case p.match(TokenTrue):
		return Literal{Value: true}
	case p.match(TokenNil):
		return Literal{Value: nil}
	case p.match(TokenNumber, TokenString):
		return Literal{Value: p.previous().Literal}
	case p.match(TokenLeftParen):
		expr := p.expression()
		p.consume(TokenRightParen, "Expect ')' after expression.")
		return Grouping{Expression: expr}
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
	loxError(token, message)
	return ParseError
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
