package lox

//go:generate go run ../tool/generate-ast.go ast

type Expr interface {
	eKind() string
}

type ExprAcceptor[R any] interface {
	accept(ExprVisitor[R]) (R, error)
}

func AcceptExpr[R any](expr Expr, visitor ExprVisitor[R]) (R, error) {
	return asExprAcceptor[R](expr).accept(visitor)
}

type Stmt interface {
	sKind() string
}

type StmtAcceptor[R any] interface {
	accept(StmtVisitor[R]) (R, error)
}

func AcceptStmt[R any](stmt Stmt, visitor StmtVisitor[R]) (R, error) {
	return asStmtAcceptor[R](stmt).accept(visitor)
}
