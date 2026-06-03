package lox

import (
	"fmt"
	"strconv"
)

type Scanner struct {
	source  []rune
	tokens  []Token
	start   int
	current int
	line    int
	err     error
}

func NewScanner(source string) *Scanner {
	return &Scanner{source: []rune(source), line: 1}
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	for !s.isAtEnd() && s.err == nil {
		// We're at the beginning of the next lexeme.
		s.start = s.current
		s.scanToken()
	}
	if s.err != nil {
		return nil, s.err
	}
	s.tokens = append(s.tokens, NewToken(TokenEof, "", nil, s.line))
	return s.tokens, nil
}

func (s *Scanner) report(line int, message string) {
	s.err = fmt.Errorf("[line %d] Error: %s", line, message)
}

func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(TokenLeftParen)
	case ')':
		s.addToken(TokenRightParen)
	case '{':
		s.addToken(TokenLeftBrace)
	case '}':
		s.addToken(TokenRightBrace)
	case ',':
		s.addToken(TokenComma)
	case '.':
		s.addToken(TokenDot)
	case '-':
		s.addToken(TokenMinus)
	case '+':
		s.addToken(TokenPlus)
	case ';':
		s.addToken(TokenSemicolon)
	case '*':
		s.addToken(TokenStar)

	case '!':
		tt := TokenBang
		if s.match('=') {
			tt = TokenBangEqual
		}
		s.addToken(tt)
	case '=':
		tt := TokenEqual
		if s.match('=') {
			tt = TokenEqualEqual
		}
		s.addToken(tt)
	case '<':
		tt := TokenLess
		if s.match('=') {
			tt = TokenLessEqual
		}
		s.addToken(tt)
	case '>':
		tt := TokenGreater
		if s.match('=') {
			tt = TokenGreaterEqual
		}
		s.addToken(tt)

	case '/':
		if s.match('/') {
			// A comment goes until the end of the line.
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else {
			s.addToken(TokenSlash)
		}

	case ' ', '\r', '\t':
		// Ignore whitespace.
	case '\n':
		s.line++

	case '"':
		s.string()

	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			s.report(s.line, "Unexpected character.")
		}
	}
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}
	text := s.source[s.start:s.current]
	tt, found := keywords[string(text)]
	if !found {
		tt = TokenIdentifier
	}
	s.addToken(tt)
}

func (s *Scanner) number() {
	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		// Consume the "."
		s.advance()
		for isDigit(s.peek()) {
			s.advance()
		}
	}

	value, err := strconv.ParseFloat(string(s.source[s.start:s.current]), 64)
	if err != nil {
		s.report(s.line, "parsing float: "+err.Error())
		return
	}
	s.addTokenLit(TokenNumber, value)
}

func (s *Scanner) string() {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			s.line++
		}
		s.advance()
	}

	if s.isAtEnd() {
		s.report(s.line, "Unterminated string.")
		return
	}

	// The closing ".
	s.advance()

	// Trim the surrounding quotes/
	value := s.source[s.start+1 : s.current-1]
	s.addTokenLit(TokenString, string(value))
}

func (s *Scanner) match(expected rune) bool {
	if s.isAtEnd() {
		return false
	}
	if s.source[s.current] != expected {
		return false
	}
	s.current++
	return true
}

func (s *Scanner) peek() rune {
	if s.isAtEnd() {
		return 0
	}
	return s.source[s.current]
}

func (s *Scanner) peekNext() rune {
	if s.current+1 >= len(s.source) {
		return 0
	}
	return s.source[s.current+1]
}

func isAlpha(c rune) bool {
	return (c >= 'a' && c <= 'z') ||
		(c >= 'A' && c <= 'Z') ||
		c == '_'
}

func isAlphaNumeric(c rune) bool {
	return isAlpha(c) || isDigit(c)
}

func isDigit(c rune) bool {
	return c >= '0' && c <= '9'
}

func (s *Scanner) advance() rune {
	r := s.source[s.current]
	s.current++
	return r
}

func (s *Scanner) addToken(tt TokenType) {
	s.addTokenLit(tt, nil)
}

func (s *Scanner) addTokenLit(tt TokenType, literal any) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, NewToken(tt, string(text), literal, s.line))
}
