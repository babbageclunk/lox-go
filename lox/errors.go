package lox

import "fmt"

type tokenError struct {
	token   Token
	message string
}

func (e tokenError) Error() string {
	return e.message
}

func newTokenError(token Token, message string, vals ...any) error {
	if len(vals) != 0 {
		message = fmt.Sprintf(message, vals...)
	}
	return tokenError{
		token:   token,
		message: message,
	}
}
