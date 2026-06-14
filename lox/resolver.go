package lox

type Resolver struct {
	interpreter *Interpreter
	scopes      Stack[map[string]bool]
}

var (
	_ StmtVisitor[Void] = &Resolver{}
	_ ExprVisitor[Void] = &Resolver{}
)

func (r *Resolver) resolveStatements(statements []Stmt) error {
	for _, statement := range statements {
		if err := r.resolveStmt(statement); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) VisitBlockStmt(stmt BlockStmt) (Void, error) {
	r.beginScope()
	if err := r.resolveStatements(stmt.Statements); err != nil {
		return void, err
	}
	r.endScope()
	return void, nil
}

func (r *Resolver) VisitBreakStmt(stmt BreakStmt) (Void, error) {
	return void, nil
}

func (r *Resolver) VisitExpressionStmt(stmt ExpressionStmt) (Void, error) {
	return void, r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitFunctionStmt(stmt FunctionStmt) (Void, error) {
	r.declare(stmt.Name)
	r.define(stmt.Name)
	return void, r.resolveExpr(stmt.Function)
}

func (r *Resolver) VisitIfStmt(stmt IfStmt) (Void, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return void, err
	}
	if err := r.resolveStmt(stmt.ThenBranch); err != nil {
		return void, err
	}
	if stmt.ElseBranch != nil {
		if err := r.resolveStmt(stmt.ElseBranch); err != nil {
			return void, err
		}
	}
	return void, nil
}

func (r *Resolver) VisitPrintStmt(stmt PrintStmt) (Void, error) {
	return void, r.resolveExpr(stmt.Expression)
}

func (r *Resolver) VisitReturnStmt(stmt ReturnStmt) (Void, error) {
	if stmt.Value != nil {
		return void, r.resolveExpr(stmt.Value)
	}
	return void, nil
}

func (r *Resolver) VisitVarStmt(stmt VarStmt) (Void, error) {
	r.declare(stmt.Name)
	if stmt.Initializer != nil {
		if err := r.resolveExpr(stmt.Initializer); err != nil {
			return void, err
		}
	}
	r.define(stmt.Name)
	return void, nil
}

func (r *Resolver) VisitWhileStmt(stmt WhileStmt) (Void, error) {
	if err := r.resolveExpr(stmt.Condition); err != nil {
		return void, err
	}
	return void, r.resolveStmt(stmt.Body)
}

func (r *Resolver) VisitAssignExpr(expr AssignExpr) (Void, error) {
	if err := r.resolveExpr(expr.Value); err != nil {
		return void, err
	}
	r.resolveLocal(expr, expr.Name)
	return void, nil
}

func (r *Resolver) VisitBinaryExpr(expr BinaryExpr) (Void, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return void, err
	}
	return void, r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitCallExpr(expr CallExpr) (Void, error) {
	if err := r.resolveExpr(expr.Callee); err != nil {
		return void, err
	}
	for _, argument := range expr.Arguments {
		if err := r.resolveExpr(argument); err != nil {
			return void, err
		}
	}
	return void, nil
}

func (r *Resolver) VisitFunctionExpr(expr FunctionExpr) (Void, error) {
	return void, r.resolveFunction(expr)
}

func (r *Resolver) VisitGroupingExpr(expr GroupingExpr) (Void, error) {
	return void, r.resolveExpr(expr.Expression)
}

func (r *Resolver) VisitLiteralExpr(expr LiteralExpr) (Void, error) {
	return void, nil
}

func (r *Resolver) VisitLogicalExpr(expr LogicalExpr) (Void, error) {
	if err := r.resolveExpr(expr.Left); err != nil {
		return void, err
	}
	return void, r.resolveExpr(expr.Right)
}

func (r *Resolver) VisitVariableExpr(expr VariableExpr) (Void, error) {
	if !r.scopes.isEmpty() && !r.scopes.peek()[expr.Name.Lexeme] {
		return void, newTokenError(
			expr.Name, "Can't read local variable in its own initializer.")
	}
	r.resolveLocal(expr, expr.Name)
	return void, nil
}

func (r *Resolver) VisitUnaryExpr(expr UnaryExpr) (Void, error) {
	return void, r.resolveExpr(expr.Right)
}

func (r *Resolver) resolveStmt(stmt Stmt) error {
	_, err := AcceptStmt(stmt, r)
	return err
}

func (r *Resolver) resolveExpr(expr Expr) error {
	_, err := AcceptExpr(expr, r)
	return err
}

func (r *Resolver) resolveFunction(expr FunctionExpr) error {
	r.beginScope()
	for _, param := range expr.Params {
		r.declare(param)
		r.define(param)
	}
	if err := r.resolveStatements(expr.Body); err != nil {
		return err
	}
	r.endScope()
	return nil
}

func (r *Resolver) beginScope() {
	r.scopes.push(make(map[string]bool))
}

func (r *Resolver) endScope() {
	r.scopes.pop()
}

func (r *Resolver) declare(name Token) {
	if r.scopes.isEmpty() {
		return
	}
	r.scopes.peek()[name.Lexeme] = false
}

func (r *Resolver) define(name Token) {
	if r.scopes.isEmpty() {
		return
	}
	r.scopes.peek()[name.Lexeme] = true
}

func (r *Resolver) resolveLocal(expr Expr, name Token) {
	for i, scope := range r.scopes.walk() {
		if _, found := scope[name.Lexeme]; found {
			r.interpreter.resolve(expr, r.scopes.len()-1-i)
			return
		}
	}
}
