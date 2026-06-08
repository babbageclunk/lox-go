package lox

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type Interpreter struct {
	globals, environment *Environment
}

func NewInterpreter() Interpreter {
	globals := NewEnvironment()
	globals.define("clock", newBuiltin(0, func(args []any) (any, error) {
		return time.Now().Unix(), nil
	}))
	return Interpreter{
		environment: globals,
		globals:     globals,
	}
}

func (i *Interpreter) Interpret(statements []Stmt) error {
	for _, statement := range statements {
		flow, err := i.Execute(statement)
		if err != nil {
			return err
		}
		// The parser prevents this situation, but keep the check as a
		// backstop for manually constructed ASTs.
		switch stmt := flow.stmt.(type) {
		case BreakStmt:
			return newTokenError(
				stmt.Keyword, "Break statement outside loop.")
		case ReturnStmt:
			return newTokenError(
				stmt.Keyword, "Return statement outside function or method.")
		}
	}
	return nil
}

func stringify(val any) string {
	if val == nil {
		return "nil"
	}
	if fval, ok := val.(float64); ok {
		str := strconv.FormatFloat(fval, 'f', -1, 64)
		if strings.HasSuffix(str, ".0") {
			str = str[:len(str)-2]
		}
		return str
	}

	return fmt.Sprintf("%v", val)
}

func (i *Interpreter) VisitLiteralExpr(expr LiteralExpr) (any, error) {
	return expr.Value, nil
}

func (i *Interpreter) VisitLogicalExpr(expr LogicalExpr) (any, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	if expr.Operator.Type == TokenOr {
		if i.isTruthy(left) {
			return left, nil
		}
	} else {
		if !i.isTruthy(left) {
			return left, nil
		}
	}
	return i.Evaluate(expr.Right)
}

func (i *Interpreter) VisitUnaryExpr(expr UnaryExpr) (any, error) {
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case TokenMinus:
		if err := i.checkNumberOperand(expr.Operator, right); err != nil {
			return nil, err
		}
		return -(right.(float64)), nil
	case TokenBang:
		return !i.isTruthy(right), nil
	}
	// Unreachable.
	return nil, nil
}

func (i *Interpreter) VisitVariableExpr(expr VariableExpr) (any, error) {
	return i.environment.get(expr.Name)
}

func (i *Interpreter) checkNumberOperand(operator Token, operand any) error {
	if _, ok := operand.(float64); !ok {
		return newTokenError(operator, "Operand must be a number.")
	}
	return nil
}

func (i *Interpreter) checkNumberOperands(operator Token, left, right any) error {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if !leftOk || !rightOk {
		return newTokenError(operator, "Operands must be numbers.")
	}
	return nil
}

func (i *Interpreter) isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return true
}

func (i *Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i *Interpreter) VisitGroupingExpr(expr GroupingExpr) (any, error) {
	return i.Evaluate(expr.Expression)
}

func (i *Interpreter) VisitBinaryExpr(expr BinaryExpr) (any, error) {
	left, err := i.Evaluate(expr.Left)
	if err != nil {
		return nil, err
	}
	right, err := i.Evaluate(expr.Right)
	if err != nil {
		return nil, err
	}
	switch expr.Operator.Type {
	case TokenGreater:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) > right.(float64), nil
	case TokenGreaterEqual:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) >= right.(float64), nil
	case TokenLess:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) < right.(float64), nil
	case TokenLessEqual:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) <= right.(float64), nil
	case TokenMinus:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) - right.(float64), nil
	case TokenSlash:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		denominator := right.(float64)
		if denominator == 0 {
			return nil, newTokenError(expr.Operator, "Division by zero.")
		}
		return left.(float64) / right.(float64), nil
	case TokenStar:
		if err := i.checkNumberOperands(expr.Operator, left, right); err != nil {
			return nil, err
		}
		return left.(float64) * right.(float64), nil
	case TokenPlus:
		lFloat, lIsFloat := left.(float64)
		rFloat, rIsFloat := right.(float64)
		if lIsFloat && rIsFloat {
			return lFloat + rFloat, nil
		}
		lStr, lIsStr := left.(string)
		rStr, rIsStr := right.(string)
		if lIsStr && rIsStr {
			return lStr + rStr, nil
		}
		return nil, newTokenError(expr.Operator, "Operands must be two numbers or two strings.")
	case TokenBangEqual:
		return !i.isEqual(left, right), nil
	case TokenEqualEqual:
		return i.isEqual(left, right), nil
	}
	// Unreachable.
	return nil, nil
}

func (i *Interpreter) VisitCallExpr(expr CallExpr) (any, error) {
	callee, err := i.Evaluate(expr.Callee)
	if err != nil {
		return nil, err
	}
	arguments := make([]any, len(expr.Arguments))
	for n, argument := range expr.Arguments {
		value, err := i.Evaluate(argument)
		if err != nil {
			return nil, err
		}
		arguments[n] = value
	}
	function, ok := callee.(callable)
	if !ok {
		return nil, newTokenError(expr.Paren, "Can only call functions and classes.")
	}
	if len(arguments) != function.arity() {
		return nil, newTokenError(
			expr.Paren,
			"Expected %d arguments but got %d.",
			function.arity(),
			len(arguments),
		)
	}
	return function.call(i, arguments)
}

func (i *Interpreter) Evaluate(expr Expr) (any, error) {
	return AcceptExpr(expr, i)
}

func (i *Interpreter) Execute(stmt Stmt) (controlFlow, error) {
	return AcceptStmt(stmt, i)
}

func (i *Interpreter) ExecuteBlock(statements []Stmt, env *Environment) (controlFlow, error) {
	prev := i.environment
	defer func() {
		i.environment = prev
	}()
	i.environment = env
	for _, statement := range statements {
		flow, err := i.Execute(statement)
		if err != nil {
			return flowNone, err
		}
		if flow != flowNone {
			return flow, nil
		}
	}
	return flowNone, nil
}

type controlFlow struct {
	stmt  Stmt
	value any
}

var flowNone = controlFlow{}

func (i *Interpreter) VisitBreakStmt(stmt BreakStmt) (controlFlow, error) {
	return controlFlow{
		stmt: stmt,
	}, nil
}

func (i *Interpreter) VisitBlockStmt(stmt BlockStmt) (controlFlow, error) {
	return i.ExecuteBlock(stmt.Statements, NewNestedEnvironment(i.environment))
}

func (i *Interpreter) VisitExpressionStmt(stmt ExpressionStmt) (controlFlow, error) {
	if _, err := i.Evaluate(stmt.Expression); err != nil {
		return flowNone, err
	}
	return flowNone, nil
}

func (i *Interpreter) VisitFunctionStmt(stmt FunctionStmt) (controlFlow, error) {
	i.environment.define(stmt.Name.Lexeme, newFunction(stmt, i.environment))
	return flowNone, nil
}

func (i *Interpreter) VisitIfStmt(stmt IfStmt) (controlFlow, error) {
	value, err := i.Evaluate(stmt.Condition)
	if err != nil {
		return flowNone, err
	}
	if i.isTruthy(value) {
		return i.Execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return i.Execute(stmt.ElseBranch)
	}
	return flowNone, nil
}

func (i *Interpreter) VisitPrintStmt(stmt PrintStmt) (controlFlow, error) {
	value, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return flowNone, err
	}
	fmt.Println(stringify(value))
	return flowNone, nil
}

func (i *Interpreter) VisitReturnStmt(stmt ReturnStmt) (controlFlow, error) {
	var value any
	if stmt.Value != nil {
		var err error
		value, err = i.Evaluate(stmt.Value)
		if err != nil {
			return flowNone, err
		}
	}
	return controlFlow{
		stmt:  stmt,
		value: value,
	}, nil

}

func (i *Interpreter) VisitVarStmt(stmt VarStmt) (controlFlow, error) {
	var value any
	if stmt.Initializer != nil {
		var err error
		value, err = i.Evaluate(stmt.Initializer)
		if err != nil {
			return flowNone, err
		}
	} else {
		value = uninitialised
	}

	i.environment.define(stmt.Name.Lexeme, value)
	return flowNone, nil
}

func (i *Interpreter) VisitWhileStmt(stmt WhileStmt) (controlFlow, error) {
	for {
		val, err := i.Evaluate(stmt.Condition)
		if err != nil {
			return flowNone, err
		}
		if !i.isTruthy(val) {
			break
		}
		flow, err := i.Execute(stmt.Body)
		if err != nil {
			return flowNone, err
		}
		switch flow.stmt.(type) {
		case BreakStmt:
			return flowNone, nil // Consume the break here.
		case ReturnStmt:
			return flow, nil // Allow return value to propagate up.
		}
	}
	return flowNone, nil
}

func (i *Interpreter) VisitAssignExpr(expr AssignExpr) (any, error) {
	value, err := i.Evaluate(expr.Value)
	if err != nil {
		return nil, err
	}
	err = i.environment.assign(expr.Name, value)
	if err != nil {
		return nil, err
	}
	return value, nil
}

var _ ExprVisitor[any] = &Interpreter{}
var _ StmtVisitor[controlFlow] = &Interpreter{}
