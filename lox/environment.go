package lox

type Environment map[string]any

func NewEnvironment() Environment {
	return make(Environment)
}

func (e Environment) define(name string, value any) {
	e[name] = value
}

func (e Environment) get(name Token) (any, error) {
	val, found := e[name.Lexeme]
	if found {
		return val, nil
	}

	return nil, newTokenError(name, "Undefined variable %q.", name.Lexeme)
}

func (e Environment) assign(name Token, value any) error {
	if _, found := e[name.Lexeme]; found {
		e[name.Lexeme] = value
		return nil
	}
	return newTokenError(name, "Undefined variable %q.", name.Lexeme)
}
