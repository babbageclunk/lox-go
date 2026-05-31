package lox

type callable interface {
	call(*Interpreter, []any) (any, error)
	arity() int
}
