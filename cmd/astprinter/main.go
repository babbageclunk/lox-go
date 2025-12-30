package main

import (
	"fmt"
	"strings"

	"github.com/babbageclunk/lox-go/lox"
)

// public static void main(String[] args) {
//   Expr expression = new Expr.Binary(
//       new Expr.Unary(
//           new Token(TokenType.MINUS, "-", null, 1),
//           new Expr.Literal(123)),
//       new Token(TokenType.STAR, "*", null, 1),
//       new Expr.Grouping(
//           new Expr.Literal(45.67)));

//   System.out.println(new AstPrinter().print(expression));
// }

func main() {
	expr := lox.Binary{
		Left: lox.Unary{
			Operator: lox.Token{
				TokenType: lox.TokenMinus,
				Lexeme:    "-",
			},
			Right: lox.Literal{Value: 123},
		},
		Operator: lox.Token{
			TokenType: lox.TokenStar,
			Lexeme:    "*",
		},
		Right: lox.Grouping{Expression: lox.Literal{Value: 45.67}},
	}
	fmt.Println(print(expr))
}

// class AstPrinter implements Expr.Visitor<String> {
//   String print(Expr expr) {
//     return expr.accept(this);
//   }
// }

type AstPrinter struct{}

func print(expr lox.Expr) string {
	return lox.Accept(expr, AstPrinter{})
}

// @Override
// public String visitBinaryExpr(Expr.Binary expr) {
//   return parenthesize(expr.operator.lexeme,
//                       expr.left, expr.right);
// }

func (p AstPrinter) VisitBinaryExpr(expr lox.Binary) string {
	return p.parenthesize(expr.Operator.Lexeme, expr.Left, expr.Right)
}

// @Override
// public String visitGroupingExpr(Expr.Grouping expr) {
//   return parenthesize("group", expr.expression);
// }

func (p AstPrinter) VisitGroupingExpr(expr lox.Grouping) string {
	return p.parenthesize("group", expr.Expression)
}

// @Override
// public String visitLiteralExpr(Expr.Literal expr) {
//   if (expr.value == null) return "nil";
//   return expr.value.toString();
// }

func (p AstPrinter) VisitLiteralExpr(expr lox.Literal) string {
	if expr.Value == nil {
		return "nil"
	}
	return fmt.Sprint(expr.Value)
}

// @Override
// public String visitUnaryExpr(Expr.Unary expr) {
//   return parenthesize(expr.operator.lexeme, expr.right);
// }

func (p AstPrinter) VisitUnaryExpr(expr lox.Unary) string {
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

func (p AstPrinter) parenthesize(name string, exprs ...lox.Expr) string {
	var b strings.Builder
	b.WriteString("(")
	b.WriteString(name)
	for _, expr := range exprs {
		b.WriteString(" ")
		b.WriteString(lox.Accept(expr, p))
	}
	b.WriteString(")")
	return b.String()
}
