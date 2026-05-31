package lox

import "fmt"

type ExprVisitor[R any] interface {
	VisitAssignExpr(expr AssignExpr) (R, error)
	VisitBinaryExpr(expr BinaryExpr) (R, error)
	VisitCallExpr(expr CallExpr) (R, error)
	VisitGroupingExpr(expr GroupingExpr) (R, error)
	VisitLiteralExpr(expr LiteralExpr) (R, error)
	VisitLogicalExpr(expr LogicalExpr) (R, error)
	VisitUnaryExpr(expr UnaryExpr) (R, error)
	VisitVariableExpr(expr VariableExpr) (R, error)
}

type AssignExpr struct {
	Name Token
	Value Expr
}

func (AssignExpr) eKind() string {
	return "AssignExpr"
}

type AssignExprAcceptor[R any] AssignExpr

func (a AssignExprAcceptor[R]) accept(vis ExprVisitor[R]) (R, error) {
	return vis.VisitAssignExpr(AssignExpr(a))
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

func (b BinaryExprAcceptor[R]) accept(vis ExprVisitor[R]) (R, error) {
	return vis.VisitBinaryExpr(BinaryExpr(b))
}

type CallExpr struct {
	Callee Expr
	Paren Token
	Arguments []Expr
}

func (CallExpr) eKind() string {
	return "CallExpr"
}

type CallExprAcceptor[R any] CallExpr

func (c CallExprAcceptor[R]) accept(vis ExprVisitor[R]) (R, error) {
	return vis.VisitCallExpr(CallExpr(c))
}

type GroupingExpr struct {
	Expression Expr
}

func (GroupingExpr) eKind() string {
	return "GroupingExpr"
}

type GroupingExprAcceptor[R any] GroupingExpr

func (g GroupingExprAcceptor[R]) accept(vis ExprVisitor[R]) (R, error) {
	return vis.VisitGroupingExpr(GroupingExpr(g))
}

type LiteralExpr struct {
	Value any
}

func (LiteralExpr) eKind() string {
	return "LiteralExpr"
}

type LiteralExprAcceptor[R any] LiteralExpr

func (l LiteralExprAcceptor[R]) accept(vis ExprVisitor[R]) (R, error) {
	return vis.VisitLiteralExpr(LiteralExpr(l))
}

type LogicalExpr struct {
	Left Expr
	Operator Token
	Right Expr
}

func (LogicalExpr) eKind() string {
	return "LogicalExpr"
}

type LogicalExprAcceptor[R any] LogicalExpr

func (l LogicalExprAcceptor[R]) accept(vis ExprVisitor[R]) (R, error) {
	return vis.VisitLogicalExpr(LogicalExpr(l))
}

type UnaryExpr struct {
	Operator Token
	Right Expr
}

func (UnaryExpr) eKind() string {
	return "UnaryExpr"
}

type UnaryExprAcceptor[R any] UnaryExpr

func (u UnaryExprAcceptor[R]) accept(vis ExprVisitor[R]) (R, error) {
	return vis.VisitUnaryExpr(UnaryExpr(u))
}

type VariableExpr struct {
	Name Token
}

func (VariableExpr) eKind() string {
	return "VariableExpr"
}

type VariableExprAcceptor[R any] VariableExpr

func (v VariableExprAcceptor[R]) accept(vis ExprVisitor[R]) (R, error) {
	return vis.VisitVariableExpr(VariableExpr(v))
}

func asExprAcceptor[R any](expr Expr) ExprAcceptor[R] {
	switch e := expr.(type) {
	case AssignExpr:
		return AssignExprAcceptor[R](e)
	case BinaryExpr:
		return BinaryExprAcceptor[R](e)
	case CallExpr:
		return CallExprAcceptor[R](e)
	case GroupingExpr:
		return GroupingExprAcceptor[R](e)
	case LiteralExpr:
		return LiteralExprAcceptor[R](e)
	case LogicalExpr:
		return LogicalExprAcceptor[R](e)
	case UnaryExpr:
		return UnaryExprAcceptor[R](e)
	case VariableExpr:
		return VariableExprAcceptor[R](e)
	}
	panic(fmt.Errorf("no acceptor for expr %s", expr.eKind()))
}
