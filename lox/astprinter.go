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

func Print(expr Expr) string {
	return Accept(expr, AstPrinter{})
}

// @Override
// public String visitBinaryExpr(Expr.Binary expr) {
//   return parenthesize(expr.operator.lexeme,
//                       expr.left, expr.right);
// }

func (p AstPrinter) VisitBinaryExpr(expr Binary) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

// @Override
// public String visitGroupingExpr(Expr.Grouping expr) {
//   return parenthesize("group", expr.expression);
// }

func (p AstPrinter) VisitGroupingExpr(expr Grouping) string {
	return p.parenthesize("group", expr.Expression)
}

// @Override
// public String visitLiteralExpr(Expr.Literal expr) {
//   if (expr.value == null) return "nil";
//   return expr.value.toString();
// }

func (p AstPrinter) VisitLiteralExpr(expr Literal) string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Value)
}

// @Override
// public String visitUnaryExpr(Expr.Unary expr) {
//   return parenthesize(expr.operator.lexeme, expr.right);
// }

func (p AstPrinter) VisitUnaryExpr(expr Unary) string {
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

func (p AstPrinter) parenthesize(name string, exprs ...Expr) string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(name)
	for _, expr := range exprs {
		b.WriteString(" ")
		b.WriteString(Accept(expr, p))
	}
	b.WriteString(")")
	return b.String()
}
