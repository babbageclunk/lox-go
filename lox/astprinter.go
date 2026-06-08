package lox

import (
	"fmt"
	"strings"
)

type AstPrinter struct{}

func (AstPrinter) Print(expr Expr) (string, error) {
	return AcceptExpr(expr, AstPrinter{})
}

func (p AstPrinter) VisitAssignExpr(expr AssignExpr) (string, error) {
	return p.parenthesize(fmt.Sprintf("assign %q", expr.Name.Lexeme), expr.Value)
}

func (p AstPrinter) VisitBinaryExpr(expr BinaryExpr) (string, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p AstPrinter) VisitCallExpr(expr CallExpr) (string, error) {
	items := []Expr{expr.Callee}
	items = append(items, expr.Arguments...)
	return p.parenthesize("call", items...)
}

func (p AstPrinter) VisitFunctionExpr(expr FunctionExpr) (string, error) {
	params := make([]string, len(expr.Params))
	for i, p := range expr.Params {
		params[i] = p.Lexeme
	}
	return p.parenthesize(fmt.Sprintf("fun(%s)", strings.Join(params, ", ")))
}

func (p AstPrinter) VisitGroupingExpr(expr GroupingExpr) (string, error) {
	return p.parenthesize("group", expr.Expression)
}

func (p AstPrinter) VisitLiteralExpr(expr LiteralExpr) (string, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprint(expr.Value), nil
}

func (p AstPrinter) VisitLogicalExpr(expr LogicalExpr) (string, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

func (p AstPrinter) VisitUnaryExpr(expr UnaryExpr) (string, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

func (p AstPrinter) VisitVariableExpr(expr VariableExpr) (string, error) {
	return p.parenthesize("var", expr)
}

func (p AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(name)
	for _, expr := range exprs {
		b.WriteString(" ")
		res, err := AcceptExpr(expr, p)
		if err != nil {
			return "", err
		}
		b.WriteString(res)
	}
	b.WriteString(")")
	return b.String(), nil
}

var _ ExprVisitor[string] = AstPrinter{}
