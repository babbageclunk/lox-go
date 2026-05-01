package lox

import "fmt"

type ExprVisitor[R any] interface {
	VisitBinaryExpr(expr BinaryExpr) (R, error)
	VisitGroupingExpr(expr GroupingExpr) (R, error)
	VisitLiteralExpr(expr LiteralExpr) (R, error)
	VisitUnaryExpr(expr UnaryExpr) (R, error)
}

type BinaryExpr struct {
	Left Expr
	Operator Token
	Right Expr
}

func (BinaryExpr) eKind() string {
	return "BinaryExpr"
}

type BinaryExprAcceptor[R any] BinaryExpr

func (b BinaryExprAcceptor[R]) accept(v ExprVisitor[R]) (R, error) {
	return v.VisitBinaryExpr(BinaryExpr(b))
}

type GroupingExpr struct {
	Expression Expr
}

func (GroupingExpr) eKind() string {
	return "GroupingExpr"
}

type GroupingExprAcceptor[R any] GroupingExpr

func (g GroupingExprAcceptor[R]) accept(v ExprVisitor[R]) (R, error) {
	return v.VisitGroupingExpr(GroupingExpr(g))
}

type LiteralExpr struct {
	Value any
}

func (LiteralExpr) eKind() string {
	return "LiteralExpr"
}

type LiteralExprAcceptor[R any] LiteralExpr

func (l LiteralExprAcceptor[R]) accept(v ExprVisitor[R]) (R, error) {
	return v.VisitLiteralExpr(LiteralExpr(l))
}

type UnaryExpr struct {
	Operator Token
	Right Expr
}

func (UnaryExpr) eKind() string {
	return "UnaryExpr"
}

type UnaryExprAcceptor[R any] UnaryExpr

func (u UnaryExprAcceptor[R]) accept(v ExprVisitor[R]) (R, error) {
	return v.VisitUnaryExpr(UnaryExpr(u))
}

func asExprAcceptor[R any](expr Expr) ExprAcceptor[R] {
	switch e := expr.(type) {
	case BinaryExpr:
		return BinaryExprAcceptor[R](e)
	case GroupingExpr:
		return GroupingExprAcceptor[R](e)
	case LiteralExpr:
		return LiteralExprAcceptor[R](e)
	case UnaryExpr:
		return UnaryExprAcceptor[R](e)
	}
	panic(fmt.Errorf("no acceptor for expr %s", expr.eKind()))
}
