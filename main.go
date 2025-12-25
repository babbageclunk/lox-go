package main

import (
	"fmt"
	"os"

	"github.com/babbageclunk/lox-go/lox"
)

func main() {
	switch {
	case len(os.Args) > 2:
		fmt.Fprintln(os.Stderr, "Usage: lox-go [script]")
		os.Exit(64)
	case len(os.Args) == 2:
		if err := lox.RunFile(os.Args[1]); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	default:
		if err := lox.RunPrompt(); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
