package lox

type callable interface {
	call(*Interpreter, []any) (any, error)
	arity() int
}

func newBuiltin(arity int, function func(args []any) (any, error)) builtin {
	return builtin{
		argCount: arity,
		function: function,
	}
}

type builtin struct {
	argCount int
	function func(args []any) (any, error)
}

func (b builtin) arity() int {
	return b.argCount
}

func (b builtin) call(_ *Interpreter, args []any) (any, error) {
	return b.function(args)
}

func (b builtin) String() string {
	return "<native fn>"
}

var _ callable = builtin{}
