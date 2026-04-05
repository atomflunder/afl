package main

import (
	"fmt"
	"os"

	"AtomflundersProgrammingLanguage/src/parser"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run main.go <source-file>")
		return
	}

	args := os.Args[1]

	if args == "" {
		fmt.Println("Usage: go run main.go <source-file>")
		return
	}

	source, err := os.ReadFile(args)

	if len(args) < 4 || args[len(args)-4:] != ".afl" {
		fmt.Printf("Error: File must have .afl extension\n")
		return
	}

	if err != nil {
		fmt.Printf("Error reading file: %s\n", err)
		return
	}

	parser, err := parser.NewParser(string(source))

	if err != nil {
		fmt.Printf("Error initializing parser: %s\n", err)
		return
	}

	ast, err := parser.GetAst()

	if err != nil {
		fmt.Printf("Error tokenizing input: %s\n", err)
		return
	}

	fmt.Printf("Parsed AST: %+v\n", ast)
}
