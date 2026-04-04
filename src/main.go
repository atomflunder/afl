package main

import (
	"fmt"
	"os"

	"AtomflundersProgrammingLanguage/src/parser"
)

func main() {

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

	tokens, err := parser.TokenizeInput(string(source))

	if err != nil {
		fmt.Printf("Error tokenizing input: %s\n", err)
		return
	}

	for _, token := range tokens {
		fmt.Printf("%s: '%s'\n", token.Type, token.Value)
	}

	// TODO: Parsing
}
