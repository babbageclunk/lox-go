package lox

import (
	"fmt"
	"strings"
)

// class AstPrinter implements Expr.Visitor<String> {
//   String print(Expr expr) {
//     return expr.accept(this);
//   }
// }

type AstPrinter struct{}

func Print(expr Expr) (string, error) {
	return Accept(expr, AstPrinter{})
}

// @Override
// public String visitBinaryExpr(Expr.Binary expr) {
//   return parenthesize(expr.operator.lexeme,
//                       expr.left, expr.right);
// }

func (p AstPrinter) VisitBinaryExpr(expr Binary) (string, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

// @Override
// public String visitGroupingExpr(Expr.Grouping expr) {
//   return parenthesize("group", expr.expression);
// }

func (p AstPrinter) VisitGroupingExpr(expr Grouping) (string, error) {
	return p.parenthesize("group", expr.Expression)
}

// @Override
// public String visitLiteralExpr(Expr.Literal expr) {
//   if (expr.value == null) return "nil";
//   return expr.value.toString();
// }

func (p AstPrinter) VisitLiteralExpr(expr Literal) (string, error) {
	if expr.Value == nil {
		return "nil", nil
	}
	return fmt.Sprint(expr.Value), nil
}

// @Override
// public String visitUnaryExpr(Expr.Unary expr) {
//   return parenthesize(expr.operator.lexeme, expr.right);
// }

func (p AstPrinter) VisitUnaryExpr(expr Unary) (string, error) {
	return p.parenthesize(expr.Operator.Lexeme, expr.Right)
}

// private String parenthesize(String name, Expr... exprs) {
//   StringBuilder builder = new StringBuilder();

//   builder.append("(").append(name);
//   for (Expr expr : exprs) {
//     builder.append(" ");
//     builder.append(expr.accept(this));
//   }
//   builder.append(")");

//   return builder.toString();
// }

func (p AstPrinter) parenthesize(name string, exprs ...Expr) (string, error) {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(name)
	for _, expr := range exprs {
		b.WriteString(" ")
		res, err := Accept(expr, p)
		if err != nil {
			return "", err
		}
		b.WriteString(res)
	}
	b.WriteString(")")
	return b.String(), nil
}
