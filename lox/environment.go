package lox

type Environment struct {
	enclosing *Environment
	values    map[string]any
}

func NewEnvironment() *Environment {
	return &Environment{
		enclosing: nil,
		values:    make(map[string]any),
	}
}

func NewNestedEnvironment(enclosing *Environment) *Environment {
	e := NewEnvironment()
	e.enclosing = enclosing
	return e
}

func (e *Environment) define(name string, value any) {
	e.values[name] = value
}

func (e *Environment) get(name Token) (any, error) {
	val, found := e.values[name.Lexeme]
	if found {
		return val, nil
	}

	if e.enclosing != nil {
		return e.enclosing.get(name)
	}

	return nil, newTokenError(name, "Undefined variable %q.", name.Lexeme)
}

func (e *Environment) assign(name Token, value any) error {
	if _, found := e.values[name.Lexeme]; found {
		e.values[name.Lexeme] = value
		return nil
	}
	if e.enclosing != nil {
		return e.enclosing.assign(name, value)
	}
	return newTokenError(name, "Undefined variable %q.", name.Lexeme)
}
