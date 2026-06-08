package lox

import (
	"errors"
	"fmt"
)

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
	name    string
	expr    FunctionExpr
	closure *Environment
}

func newFunction(name string, expr FunctionExpr, closure *Environment) function {
	return function{
		name:    name,
		expr:    expr,
		closure: closure,
	}
}

func (f function) arity() int {
	return len(f.expr.Params)
}

func (f function) call(interpreter *Interpreter, args []any) (any, error) {
	environment := NewNestedEnvironment(f.closure)
	for i, param := range f.expr.Params {
		environment.define(param.Lexeme, args[i])
	}
	err := interpreter.ExecuteBlock(f.expr.Body, environment)
	if ret, ok := errors.AsType[returnControlFlow](err); ok {
		return ret.value, nil
	} else if brk, ok := errors.AsType[breakControlFlow](err); ok {
		// We shouldn't let a break propagate across a function call boundary.
		// Convert it into a regular error rather than a control-flow one.
		return nil, newTokenError(brk.stmt.Keyword, "Break statement outside loop.")
	} else if err != nil {
		return nil, err
	}
	return nil, nil
}

func (f function) String() string {
	if f.name == "" {
		return "<anon fn>"
	}
	return fmt.Sprintf("<fn %s>", f.name)
}

var _ callable = function{}
