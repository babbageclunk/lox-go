package lox

type Binary struct {
	left Expr
	operator Token
	right Expr
}

func (Binary) kind() string {
	return "Binary"
}

type Grouping struct {
	expression Expr
}

func (Grouping) kind() string {
	return "Grouping"
}

type Literal struct {
	value any
}

func (Literal) kind() string {
	return "Literal"
}

type Unary struct {
	operator Token
	right Expr
}

func (Unary) kind() string {
	return "Unary"
}

