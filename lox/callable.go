package lox

import "fmt"

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

type function struct {
	declaration FunctionStmt
}

func newFunction(declaration FunctionStmt) function {
	return function{declaration: declaration}
}

func (f function) arity() int {
	return len(f.declaration.Params)
}

func (f function) call(interpreter *Interpreter, args []any) (any, error) {
	environment := NewNestedEnvironment(interpreter.globals)
	for i, param := range f.declaration.Params {
		environment.define(param.Lexeme, args[i])
	}
	return nil, interpreter.ExecuteBlock(f.declaration.Body, environment)
}

func (f function) String() string {
	return fmt.Sprintf("<fn %s>", f.declaration.Name.Lexeme)
}

var _ callable = function{}
