package lox

type Interpreter struct{}

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
		return -(right.(float64)), nil
	case TokenBang:
		return !i.isTruthy(right), nil
	}
	// Unreachable.
	return nil, nil
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
		return left.(float64) > right.(float64), nil
	case TokenGreaterEqual:
		return left.(float64) >= right.(float64), nil
	case TokenLess:
		return left.(float64) < right.(float64), nil
	case TokenLessEqual:
		return left.(float64) <= right.(float64), nil
	case TokenMinus:
		return left.(float64) - right.(float64), nil
	case TokenSlash:
		return left.(float64) / right.(float64), nil
	case TokenStar:
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
