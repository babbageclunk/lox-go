package lox

import "fmt"

type Visitor[R any] interface {
	VisitBinaryExpr(expr Binary) (R, error)
	VisitGroupingExpr(expr Grouping) (R, error)
	VisitLiteralExpr(expr Literal) (R, error)
	VisitUnaryExpr(expr Unary) (R, error)
}

type Binary struct {
	Left Expr
	Operator Token
	Right Expr
}

func (Binary) kind() string {
	return "Binary"
}

type BinaryAcceptor[R any] Binary

func (b BinaryAcceptor[R]) accept(v Visitor[R]) (R, error) {
	return v.VisitBinaryExpr(Binary(b))
}

type Grouping struct {
	Expression Expr
}

func (Grouping) kind() string {
	return "Grouping"
}

type GroupingAcceptor[R any] Grouping

func (g GroupingAcceptor[R]) accept(v Visitor[R]) (R, error) {
	return v.VisitGroupingExpr(Grouping(g))
}

type Literal struct {
	Value any
}

func (Literal) kind() string {
	return "Literal"
}

type LiteralAcceptor[R any] Literal

func (l LiteralAcceptor[R]) accept(v Visitor[R]) (R, error) {
	return v.VisitLiteralExpr(Literal(l))
}

type Unary struct {
	Operator Token
	Right Expr
}

func (Unary) kind() string {
	return "Unary"
}

type UnaryAcceptor[R any] Unary

func (u UnaryAcceptor[R]) accept(v Visitor[R]) (R, error) {
	return v.VisitUnaryExpr(Unary(u))
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
