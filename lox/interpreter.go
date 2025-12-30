package lox

type Interpreter struct{}

func (i Interpreter) VisitLiteralExpr(expr Literal) any {
	return expr.Value
}

func (i Interpreter) VisitUnaryExpr(expr Unary) any {
	right := i.Evaluate(expr.Right)
	switch expr.Operator.Type {
	case TokenMinus:
		return -(right.(float64))
	case TokenBang:
		return !i.isTruthy(right)
	}
	// Unreachable.
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

func (i Interpreter) VisitGroupingExpr(expr Grouping) any {
	return i.Evaluate(expr.Expression)
}

func (i Interpreter) VisitBinaryExpr(expr Binary) any {
	left := i.Evaluate(expr.Left)
	right := i.Evaluate(expr.Right)
	switch expr.Operator.Type {
	case TokenGreater:
		return left.(float64) > right.(float64)
	case TokenGreaterEqual:
		return left.(float64) >= right.(float64)
	case TokenLess:
		return left.(float64) < right.(float64)
	case TokenLessEqual:
		return left.(float64) <= right.(float64)
	case TokenMinus:
		return left.(float64) - right.(float64)
	case TokenSlash:
		return left.(float64) / right.(float64)
	case TokenStar:
		return left.(float64) * right.(float64)
	case TokenPlus:
		lFloat, lIsFloat := left.(float64)
		rFloat, rIsFloat := right.(float64)
		if lIsFloat && rIsFloat {
			return lFloat + rFloat
		}
		lStr, lIsStr := left.(string)
		rStr, rIsStr := right.(string)
		if lIsStr && rIsStr {
			return lStr + rStr
		}
	case TokenBangEqual:
		return !i.isEqual(left, right)
	case TokenEqualEqual:
		return i.isEqual(left, right)
	}
	// Unreachable.
	return nil
}

func (i Interpreter) Evaluate(expr Expr) any {
	return Accept(expr, i)
}
