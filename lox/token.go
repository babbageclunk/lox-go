package lox

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal any
	Line    int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		Type:    tokenType,
		Lexeme:  lexeme,
		Literal: literal,
		Line:    line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.Type, t.Lexeme, t.Literal)
}

type TokenType string

const (
	// Single-character tokens.
	TokenLeftParen  TokenType = "LEFT_PAREN"
	TokenRightParen TokenType = "RIGHT_PAREN"
	TokenLeftBrace  TokenType = "LEFT_BRACE"
	TokenRightBrace TokenType = "RIGHT_BRACE"
	TokenComma      TokenType = "COMMA"
	TokenDot        TokenType = "DOT"
	TokenMinus      TokenType = "MINUS"
	TokenPlus       TokenType = "PLUS"
	TokenSemicolon  TokenType = "SEMICOLON"
	TokenSlash      TokenType = "SLASH"
	TokenStar       TokenType = "STAR"

	// One or two character tokens.
	TokenBang         TokenType = "BANG"
	TokenBangEqual    TokenType = "BANG_EQUAL"
	TokenEqual        TokenType = "EQUAL"
	TokenEqualEqual   TokenType = "EQUAL_EQUAL"
	TokenGreater      TokenType = "GREATER"
	TokenGreaterEqual TokenType = "GREATER_EQUAL"
	TokenLess         TokenType = "LESS"
	TokenLessEqual    TokenType = "LESS_EQUAL"

	// Literals.
	TokenIdentifier TokenType = "IDENTIFIER"
	TokenString     TokenType = "STRING"
	TokenNumber     TokenType = "NUMBER"

	// Keywords.
	TokenAnd    TokenType = "AND"
	TokenBreak  TokenType = "BREAK"
	TokenClass  TokenType = "CLASS"
	TokenElse   TokenType = "ELSE"
	TokenFalse  TokenType = "FALSE"
	TokenFun    TokenType = "FUN"
	TokenFor    TokenType = "FOR"
	TokenIf     TokenType = "IF"
	TokenNil    TokenType = "NIL"
	TokenOr     TokenType = "OR"
	TokenPrint  TokenType = "PRINT"
	TokenReturn TokenType = "RETURN"
	TokenSuper  TokenType = "SUPER"
	TokenThis   TokenType = "THIS"
	TokenTrue   TokenType = "TRUE"
	TokenVar    TokenType = "VAR"
	TokenWhile  TokenType = "WHILE"

	TokenEof TokenType = "EOF"
)

var keywords = map[string]TokenType{
	"and":    TokenAnd,
	"break":  TokenBreak,
	"class":  TokenClass,
	"else":   TokenElse,
	"false":  TokenFalse,
	"for":    TokenFor,
	"fun":    TokenFun,
	"if":     TokenIf,
	"nil":    TokenNil,
	"or":     TokenOr,
	"print":  TokenPrint,
	"return": TokenReturn,
	"super":  TokenSuper,
	"this":   TokenThis,
	"true":   TokenTrue,
	"var":    TokenVar,
	"while":  TokenWhile,
}
