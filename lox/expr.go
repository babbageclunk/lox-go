package lox

//go:generate go run ../tool/generate-ast.go ast

type Expr interface {
	kind() string
}

type Acceptor[R any] interface {
	accept(Visitor[R]) R
}

func Accept[R any](expr Expr, visitor Visitor[R]) R {
	return asAcceptor[R](expr).accept(visitor)
}
