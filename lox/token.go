package lox

import "fmt"

// class Token {
//   final TokenType type;
//   final String lexeme;
//   final Object literal;
//   final int line;

//   Token(TokenType type, String lexeme, Object literal, int line) {
//     this.type = type;
//     this.lexeme = lexeme;
//     this.literal = literal;
//     this.line = line;
//   }

//   public String toString() {
//     return type + " " + lexeme + " " + literal;
//   }
// }

type Token struct {
	TokenType TokenType
	Lexeme    string
	Literal   any
	Line      int
}

func NewToken(tokenType TokenType, lexeme string, literal any, line int) Token {
	return Token{
		TokenType: tokenType,
		Lexeme:    lexeme,
		Literal:   literal,
		Line:      line,
	}
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %v", t.TokenType, t.Lexeme, t.Literal)
}

type TokenType string

// enum TokenType {
//   // Single-character tokens.
//   LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE,
//   COMMA, DOT, MINUS, PLUS, SEMICOLON, SLASH, STAR,

//   // One or two character tokens.
//   BANG, BANG_EQUAL,
//   EQUAL, EQUAL_EQUAL,
//   GREATER, GREATER_EQUAL,
//   LESS, LESS_EQUAL,

//   // Literals.
//   IDENTIFIER, STRING, NUMBER,

//   // Keywords.
//   AND, CLASS, ELSE, FALSE, FUN, FOR, IF, NIL, OR,
//   PRINT, RETURN, SUPER, THIS, TRUE, VAR, WHILE,

//   EOF
// }

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

// static {
//   keywords = new HashMap<>();
//   keywords.put("and",    AND);
//   keywords.put("class",  CLASS);
//   keywords.put("else",   ELSE);
//   keywords.put("false",  FALSE);
//   keywords.put("for",    FOR);
//   keywords.put("fun",    FUN);
//   keywords.put("if",     IF);
//   keywords.put("nil",    NIL);
//   keywords.put("or",     OR);
//   keywords.put("print",  PRINT);
//   keywords.put("return", RETURN);
//   keywords.put("super",  SUPER);
//   keywords.put("this",   THIS);
//   keywords.put("true",   TRUE);
//   keywords.put("var",    VAR);
//   keywords.put("while",  WHILE);
// }

var keywords = map[string]TokenType{
	"and":    TokenAnd,
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
