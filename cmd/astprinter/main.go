package main

import (
	"fmt"

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
	expr := lox.BinaryExpr{
		Left: lox.UnaryExpr{
			Operator: lox.Token{
				Type:   lox.TokenMinus,
				Lexeme: "-",
			},
			Right: lox.LiteralExpr{Value: 123},
		},
		Operator: lox.Token{
			Type:   lox.TokenStar,
			Lexeme: "*",
		},
		Right: lox.GroupingExpr{Expression: lox.LiteralExpr{Value: 45.67}},
	}
	fmt.Println(lox.AstPrinter{}.Print(expr))
}
