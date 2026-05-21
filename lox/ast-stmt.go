package lox

import "fmt"

type StmtVisitor[R any] interface {
	VisitBlockStmt(stmt BlockStmt) (R, error)
	VisitExpressionStmt(stmt ExpressionStmt) (R, error)
	VisitPrintStmt(stmt PrintStmt) (R, error)
	VisitVarStmt(stmt VarStmt) (R, error)
}

type BlockStmt struct {
	Statements []Stmt
}

func (BlockStmt) sKind() string {
	return "BlockStmt"
}

type BlockStmtAcceptor[R any] BlockStmt

func (b BlockStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitBlockStmt(BlockStmt(b))
}

type ExpressionStmt struct {
	Expression Expr
}

func (ExpressionStmt) sKind() string {
	return "ExpressionStmt"
}

type ExpressionStmtAcceptor[R any] ExpressionStmt

func (e ExpressionStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitExpressionStmt(ExpressionStmt(e))
}

type PrintStmt struct {
	Expression Expr
}

func (PrintStmt) sKind() string {
	return "PrintStmt"
}

type PrintStmtAcceptor[R any] PrintStmt

func (p PrintStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitPrintStmt(PrintStmt(p))
}

type VarStmt struct {
	Name Token
	Initializer Expr
}

func (VarStmt) sKind() string {
	return "VarStmt"
}

type VarStmtAcceptor[R any] VarStmt

func (v VarStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitVarStmt(VarStmt(v))
}

func asStmtAcceptor[R any](stmt Stmt) StmtAcceptor[R] {
	switch e := stmt.(type) {
	case BlockStmt:
		return BlockStmtAcceptor[R](e)
	case ExpressionStmt:
		return ExpressionStmtAcceptor[R](e)
	case PrintStmt:
		return PrintStmtAcceptor[R](e)
	case VarStmt:
		return VarStmtAcceptor[R](e)
	}
	panic(fmt.Errorf("no acceptor for stmt %s", stmt.sKind()))
}
