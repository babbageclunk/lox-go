package lox

type tokenError struct {
	token   Token
	message string
}

func (e tokenError) Error() string {
	return e.message
}

func newTokenError(token Token, message string) error {
	return tokenError{
		token:   token,
		message: message,
	}
}
