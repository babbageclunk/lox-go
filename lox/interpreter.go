package lox

import (
	"fmt"
	"strconv"
	"strings"
)

type Interpreter struct {
	environment *Environment
	breaking    bool
}

func NewInterpreter() Interpreter {
	return Interpreter{
		environment: NewEnvironment(),
	}
}

func (i *Interpreter) Interpret(statements []Stmt) error {
	for _, statement := range statements {
		if err := i.Execute(statement); err != nil {
			return err
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

func (i *Interpreter) Execute(stmt Stmt) error {
	_, err := AcceptStmt(stmt, i)
	return err
}

func (i *Interpreter) ExecuteBlock(statements []Stmt, env *Environment) error {
	prev := i.environment
	defer func() {
		i.environment = prev
	}()
	i.environment = env
	for _, statement := range statements {
		if err := i.Execute(statement); err != nil {
			return err
		}
		if i.breaking {
			return nil
		}
	}
	return nil
}

type Void struct{}

var void Void

func (i *Interpreter) VisitBreakStmt(stmt BreakStmt) (Void, error) {
	i.breaking = true
	return void, nil
}

func (i *Interpreter) VisitBlockStmt(stmt BlockStmt) (Void, error) {
	return void, i.ExecuteBlock(stmt.Statements, NewNestedEnvironment(i.environment))
}

func (i *Interpreter) VisitExpressionStmt(stmt ExpressionStmt) (Void, error) {
	if _, err := i.Evaluate(stmt.Expression); err != nil {
		return void, err
	}
	return void, nil
}

func (i *Interpreter) VisitIfStmt(stmt IfStmt) (Void, error) {
	value, err := i.Evaluate(stmt.Condition)
	if err != nil {
		return void, err
	}
	if i.isTruthy(value) {
		return void, i.Execute(stmt.ThenBranch)
	} else if stmt.ElseBranch != nil {
		return void, i.Execute(stmt.ElseBranch)
	}
	return void, nil
}

func (i *Interpreter) VisitPrintStmt(stmt PrintStmt) (Void, error) {
	value, err := i.Evaluate(stmt.Expression)
	if err != nil {
		return void, err
	}
	fmt.Println(stringify(value))
	return void, nil
}

func (i *Interpreter) VisitVarStmt(stmt VarStmt) (Void, error) {
	var value any
	if stmt.Initializer != nil {
		var err error
		value, err = i.Evaluate(stmt.Initializer)
		if err != nil {
			return void, err
		}
	} else {
		value = uninitialised
	}

	i.environment.define(stmt.Name.Lexeme, value)
	return void, nil
}

func (i *Interpreter) VisitWhileStmt(stmt WhileStmt) (Void, error) {
	for {
		val, err := i.Evaluate(stmt.Condition)
		if err != nil {
			return void, err
		}
		if !i.isTruthy(val) {
			break
		}
		if err := i.Execute(stmt.Body); err != nil {
			return void, err
		}
		if i.breaking {
			i.breaking = false
			break
		}
	}
	return void, nil
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
