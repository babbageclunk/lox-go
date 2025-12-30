package lox

import "fmt"

type Visitor[R any] interface {
	visitBinaryExpr(expr Binary) R
	visitGroupingExpr(expr Grouping) R
	visitLiteralExpr(expr Literal) R
	visitUnaryExpr(expr Unary) R
}

type Binary struct {
	left Expr
	operator Token
	right Expr
}

func (Binary) kind() string {
	return "Binary"
}

type BinaryAcceptor[R any] Binary

func (b BinaryAcceptor[R]) accept(v Visitor[R]) R {
	return v.visitBinaryExpr(Binary(b))
}

type Grouping struct {
	expression Expr
}

func (Grouping) kind() string {
	return "Grouping"
}

type GroupingAcceptor[R any] Grouping

func (g GroupingAcceptor[R]) accept(v Visitor[R]) R {
	return v.visitGroupingExpr(Grouping(g))
}

type Literal struct {
	value any
}

func (Literal) kind() string {
	return "Literal"
}

type LiteralAcceptor[R any] Literal

func (l LiteralAcceptor[R]) accept(v Visitor[R]) R {
	return v.visitLiteralExpr(Literal(l))
}

type Unary struct {
	operator Token
	right Expr
}

func (Unary) kind() string {
	return "Unary"
}

type UnaryAcceptor[R any] Unary

func (u UnaryAcceptor[R]) accept(v Visitor[R]) R {
	return v.visitUnaryExpr(Unary(u))
}

func asAcceptor[R any](expr Expr) Acceptor[R] {
	switch e := expr.(type) {
	case Binary:
		return BinaryAcceptor[R](e)
	case Grouping:
		return GroupingAcceptor[R](e)
	case Literal:
		return LiteralAcceptor[R](e)
	case Unary:
		return UnaryAcceptor[R](e)
	}
	panic(fmt.Errorf("no acceptor for expr %s", expr.kind()))
}
