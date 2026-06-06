package lox

import "fmt"

type StmtVisitor[R any] interface {
	VisitBlockStmt(stmt BlockStmt) (R, error)
	VisitBreakStmt(stmt BreakStmt) (R, error)
	VisitExpressionStmt(stmt ExpressionStmt) (R, error)
	VisitFunctionStmt(stmt FunctionStmt) (R, error)
	VisitIfStmt(stmt IfStmt) (R, error)
	VisitPrintStmt(stmt PrintStmt) (R, error)
	VisitReturnStmt(stmt ReturnStmt) (R, error)
	VisitVarStmt(stmt VarStmt) (R, error)
	VisitWhileStmt(stmt WhileStmt) (R, error)
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

type BreakStmt struct {
	Keyword Token
}

func (BreakStmt) sKind() string {
	return "BreakStmt"
}

type BreakStmtAcceptor[R any] BreakStmt

func (b BreakStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitBreakStmt(BreakStmt(b))
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

type FunctionStmt struct {
	Name Token
	Params []Token
	Body []Stmt
}

func (FunctionStmt) sKind() string {
	return "FunctionStmt"
}

type FunctionStmtAcceptor[R any] FunctionStmt

func (f FunctionStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitFunctionStmt(FunctionStmt(f))
}

type IfStmt struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt
}

func (IfStmt) sKind() string {
	return "IfStmt"
}

type IfStmtAcceptor[R any] IfStmt

func (i IfStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitIfStmt(IfStmt(i))
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

type ReturnStmt struct {
	Keyword Token
	Value Expr
}

func (ReturnStmt) sKind() string {
	return "ReturnStmt"
}

type ReturnStmtAcceptor[R any] ReturnStmt

func (r ReturnStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitReturnStmt(ReturnStmt(r))
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

type WhileStmt struct {
	Condition Expr
	Body Stmt
}

func (WhileStmt) sKind() string {
	return "WhileStmt"
}

type WhileStmtAcceptor[R any] WhileStmt

func (w WhileStmtAcceptor[R]) accept(vis StmtVisitor[R]) (R, error) {
	return vis.VisitWhileStmt(WhileStmt(w))
}

func asStmtAcceptor[R any](stmt Stmt) StmtAcceptor[R] {
	switch e := stmt.(type) {
	case BlockStmt:
		return BlockStmtAcceptor[R](e)
	case BreakStmt:
		return BreakStmtAcceptor[R](e)
	case ExpressionStmt:
		return ExpressionStmtAcceptor[R](e)
	case FunctionStmt:
		return FunctionStmtAcceptor[R](e)
	case IfStmt:
		return IfStmtAcceptor[R](e)
	case PrintStmt:
		return PrintStmtAcceptor[R](e)
	case ReturnStmt:
		return ReturnStmtAcceptor[R](e)
	case VarStmt:
		return VarStmtAcceptor[R](e)
	case WhileStmt:
		return WhileStmtAcceptor[R](e)
	}
	panic(fmt.Errorf("no acceptor for stmt %s", stmt.sKind()))
}
