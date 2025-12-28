package lox

//go:generate go run ../tool/generate-ast.go ast

type Expr interface {
	kind() string
}
