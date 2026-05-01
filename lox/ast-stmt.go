package lox

import "fmt"

type StmtVisitor[R any] interface {
	VisitExpressionStmt(stmt ExpressionStmt) (R, error)
	VisitPrintStmt(stmt PrintStmt) (R, error)
}

type ExpressionStmt struct {
	Expression Expr
}

func (ExpressionStmt) sKind() string {
	return "ExpressionStmt"
}

type ExpressionStmtAcceptor[R any] ExpressionStmt

func (e ExpressionStmtAcceptor[R]) accept(v StmtVisitor[R]) (R, error) {
	return v.VisitExpressionStmt(ExpressionStmt(e))
}

type PrintStmt struct {
	Expression Expr
}

func (PrintStmt) sKind() string {
	return "PrintStmt"
}

type PrintStmtAcceptor[R any] PrintStmt

func (p PrintStmtAcceptor[R]) accept(v StmtVisitor[R]) (R, error) {
	return v.VisitPrintStmt(PrintStmt(p))
}

func asStmtAcceptor[R any](stmt Stmt) StmtAcceptor[R] {
	switch e := stmt.(type) {
	case ExpressionStmt:
		return ExpressionStmtAcceptor[R](e)
	case PrintStmt:
		return PrintStmtAcceptor[R](e)
	}
	panic(fmt.Errorf("no acceptor for stmt %s", stmt.sKind()))
}
