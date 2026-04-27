package lox

import (
	"fmt"
	"strconv"
	"strings"
)

type Interpreter struct{}

func (i Interpreter) Interpret(expr Expr) {
	result, err := i.Evaluate(expr)
	if err != nil {
		runtimeError(err)
		return
	}
	fmt.Println(stringify(result))
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

	return fmt.Sprintf("%s", val)
}

func (i Interpreter) VisitLiteralExpr(expr Literal) (any, error) {
	return expr.Value, nil
}

func (i Interpreter) VisitUnaryExpr(expr Unary) (any, error) {
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

func (i Interpreter) checkNumberOperand(operator Token, operand any) error {
	if _, ok := operand.(float64); !ok {
		return newTokenError(operator, "Operand must be a number.")
	}
	return nil
}

func (i Interpreter) checkNumberOperands(operator Token, left, right any) error {
	_, leftOk := left.(float64)
	_, rightOk := right.(float64)
	if !leftOk || !rightOk {
		return newTokenError(operator, "Operands must be numbers.")
	}
	return nil
}

func (i Interpreter) isTruthy(val any) bool {
	if val == nil {
		return false
	}
	if b, ok := val.(bool); ok {
		return b
	}
	return true
}

func (i Interpreter) isEqual(a, b any) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}
	return a == b
}

func (i Interpreter) VisitGroupingExpr(expr Grouping) (any, error) {
	return i.Evaluate(expr.Expression)
}

func (i Interpreter) VisitBinaryExpr(expr Binary) (any, error) {
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

func (i Interpreter) Evaluate(expr Expr) (any, error) {
	return Accept(expr, i)
}
