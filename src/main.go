package main

import (
	"flag"
	"fmt"
	"os"

	"afl/src/parser"
	"afl/src/runtime"
)

func main() {
	astFlag := flag.Bool("ast", false, "print AST")
	tokenFlag := flag.Bool("tokens", false, "print Tokens")
	flag.Parse()

	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: afl [options] <source-file>")
		return
	}

	filename := args[0]

	if len(filename) < 4 || filename[len(filename)-4:] != ".afl" {
		fmt.Printf("Error: File must have .afl extension\n")
		return
	}

	source, err := os.ReadFile(filename)

	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	tokens, err := parser.TokenizeInput(string(source))

	if *tokenFlag {
		fmt.Printf("Tokens: %v\n", tokens)
	}

	parser := parser.NewParser(tokens)

	if err != nil {
		fmt.Printf("Error initializing parser: %s\n", err)
		return
	}

	ast, err := parser.GetAst()

	if err != nil {
		fmt.Printf("Error parsing input: %s\n", err)
		return
	}

	if *astFlag {
		fmt.Printf("AST: %v\n", ast)
	}

	env := runtime.NewEnvironment(nil)
	runtime.EvaluateProgram(ast, env)
}
