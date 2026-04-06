package parser

import (
	"testing"
)

func TestParseVariableDeclaration(t *testing.T) {
	tests := []struct {
		name            string
		input           string
		expectedName    string
		expectedMutable bool
		shouldError     bool
	}{
		{
			name:            "simple variable declaration",
			input:           "x = 5",
			expectedName:    "x",
			expectedMutable: false,
			shouldError:     false,
		},
		{
			name:            "variable with type annotation",
			input:           "y: int = 10",
			expectedName:    "y",
			expectedMutable: false,
			shouldError:     false,
		},
		{
			name:            "variable with float type",
			input:           "z: float = 3.14",
			expectedName:    "z",
			expectedMutable: false,
			shouldError:     false,
		},
		{
			name:            "mutable variable",
			input:           "a? = 20",
			expectedName:    "a",
			expectedMutable: true,
			shouldError:     false,
		},
		{
			name:            "mutable with type annotation",
			input:           "b?: int = 15",
			expectedName:    "b",
			expectedMutable: true,
			shouldError:     false,
		},
		{
			name:            "variable with string value",
			input:           "name = \"Alice\"",
			expectedName:    "name",
			expectedMutable: false,
			shouldError:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}

				varDecl, ok := ast[0].(VariableDeclaration)
				if !ok {
					t.Errorf("Expected VariableDeclaration, got %T", ast[0])
					return
				}

				if varDecl.Name != tt.expectedName {
					t.Errorf("Expected name %s, got %s", tt.expectedName, varDecl.Name)
				}

				if varDecl.IsMutable != tt.expectedMutable {
					t.Errorf("Expected mutable %v, got %v", tt.expectedMutable, varDecl.IsMutable)
				}
			}
		})
	}
}

func TestParseFunctionDeclaration(t *testing.T) {
	tests := []struct {
		name               string
		input              string
		expectedFuncName   string
		expectedParamCount int
		shouldError        bool
	}{
		{
			name:               "simple function",
			input:              "fn add() { return 5 }",
			expectedFuncName:   "add",
			expectedParamCount: 0,
			shouldError:        false,
		},
		{
			name:               "function with return",
			input:              "fn getValue() { return 42 }",
			expectedFuncName:   "getValue",
			expectedParamCount: 0,
			shouldError:        false,
		},
		{
			name:               "function with print",
			input:              `fn greet() { print("Hello") }`,
			expectedFuncName:   "greet",
			expectedParamCount: 0,
			shouldError:        false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}

				funcDecl, ok := ast[0].(FunctionDeclaration)
				if !ok {
					t.Errorf("Expected FunctionDeclaration, got %T", ast[0])
					return
				}

				if funcDecl.Name != tt.expectedFuncName {
					t.Errorf("Expected function name %s, got %s", tt.expectedFuncName, funcDecl.Name)
				}

				if len(funcDecl.Parameters) != tt.expectedParamCount {
					t.Errorf("Expected %d parameters, got %d", tt.expectedParamCount, len(funcDecl.Parameters))
				}

				if len(funcDecl.Body) == 0 {
					t.Errorf("Expected function body, got empty")
				}
			}
		})
	}
}

func TestParseIfStatement(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "simple if statement",
			input:       "if (x > 5) { print(\"x is greater\") }",
			shouldError: false,
		},
		{
			name:        "if with equality",
			input:       "if (x == 10) { print(\"equal\") }",
			shouldError: false,
		},
		{
			name:        "if else statement",
			input:       "if (x > 5) { print(\"greater\") } else { print(\"not greater\") }",
			shouldError: false,
		},
		{
			name:        "if elseif else statement",
			input:       "if (x > 5) { print(\"greater\") } elseif (x == 5) { print(\"equal\") } else { print(\"less\") }",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}

				ifStmt, ok := ast[0].(IfStatement)
				if !ok {
					t.Errorf("Expected IfStatement, got %T", ast[0])
					return
				}

				if ifStmt.Condition == nil {
					t.Errorf("Expected condition in if statement")
				}

				if len(ifStmt.IfBranch.Body) == 0 {
					t.Errorf("Expected if branch body")
				}
			}
		})
	}
}

func TestParseForLoop(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedIterator string
		shouldError      bool
	}{
		{
			name:             "simple for loop",
			input:            "for (i in 0->5) { print(i) }",
			expectedIterator: "i",
			shouldError:      false,
		},
		{
			name:             "for loop with range",
			input:            "for (x in 1->10) { print(x) }",
			expectedIterator: "x",
			shouldError:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}

				forLoop, ok := ast[0].(ForLoop)
				if !ok {
					t.Errorf("Expected ForLoop, got %T", ast[0])
					return
				}

				if forLoop.Iterator != tt.expectedIterator {
					t.Errorf("Expected iterator %s, got %s", tt.expectedIterator, forLoop.Iterator)
				}

				if forLoop.Start == nil {
					t.Errorf("Expected start expression")
				}

				if forLoop.End == nil {
					t.Errorf("Expected end expression")
				}

				if len(forLoop.Body) == 0 {
					t.Errorf("Expected loop body")
				}
			}
		})
	}
}

func TestParseFunctionCall(t *testing.T) {
	tests := []struct {
		name                  string
		input                 string
		expectedFunctionName  string
		expectedArgumentCount int
		shouldError           bool
	}{
		{
			name:                  "function call with no args",
			input:                 "print()",
			expectedFunctionName:  "print",
			expectedArgumentCount: 0,
			shouldError:           false,
		},
		{
			name:                  "function call with string arg",
			input:                 `print("Hello, World!")`,
			expectedFunctionName:  "print",
			expectedArgumentCount: 1,
			shouldError:           false,
		},
		{
			name:                  "function call with multiple args",
			input:                 `print("x is", x, "and y is", y)`,
			expectedFunctionName:  "print",
			expectedArgumentCount: 4,
			shouldError:           false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}

				callExpr, ok := ast[0].(CallExpression)
				if !ok {
					t.Errorf("Expected CallExpression, got %T", ast[0])
					return
				}

				idExpr, ok := callExpr.Caller.(IdentifierExpression)
				if !ok {
					t.Errorf("Expected IdentifierExpression as caller, got %T", callExpr.Caller)
					return
				}

				if idExpr.Symbol != tt.expectedFunctionName {
					t.Errorf("Expected function %s, got %s", tt.expectedFunctionName, idExpr.Symbol)
				}

				if len(callExpr.Args) != tt.expectedArgumentCount {
					t.Errorf("Expected %d arguments, got %d", tt.expectedArgumentCount, len(callExpr.Args))
				}
			}
		})
	}
}

func TestParseBinaryExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "addition expression",
			input:       "x = a + b",
			shouldError: false,
		},
		{
			name:        "subtraction expression",
			input:       "x = a - b",
			shouldError: false,
		},
		{
			name:        "multiplication expression",
			input:       "x = a * b",
			shouldError: false,
		},
		{
			name:        "division expression",
			input:       "x = a / b",
			shouldError: false,
		},
		{
			name:        "comparison expression",
			input:       "x = a > b",
			shouldError: false,
		},
		{
			name:        "equality expression",
			input:       "x = a == b",
			shouldError: false,
		},
		{
			name:        "inequality expression",
			input:       "x = a != b",
			shouldError: false,
		},
		{
			name:        "logical and expression",
			input:       "x = a && b",
			shouldError: false,
		},
		{
			name:        "logical or expression",
			input:       "x = a || b",
			shouldError: false,
		},
		{
			name:        "compound expression",
			input:       "x = a + b * c",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}

				_, ok := ast[0].(VariableDeclaration)
				if !ok {
					t.Errorf("Expected VariableDeclaration, got %T", ast[0])
				}
			}
		})
	}
}

func TestParseUnaryExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "logical not",
			input:       "x = !a",
			shouldError: false,
		},
		{
			name:        "negation",
			input:       "x = -5",
			shouldError: false,
		},
		{
			name:        "not not",
			input:       "x = !!true",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}
			}
		})
	}
}

func TestParseMemberExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "simple member access",
			input:       "x = obj.property",
			shouldError: false,
		},
		{
			name:        "chained member access",
			input:       "x = obj.field.subfield",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}
			}
		})
	}
}

func TestParseAssignmentExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "simple assignment",
			input:       "x = 5",
			shouldError: false,
		},
		{
			name:        "chained assignment",
			input:       "x = y = 10",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}
			}
		})
	}
}

func TestParseCompleteProgram(t *testing.T) {
	input := `
	x = 5
	y: int = 10
	z? = 15
	
	fn add(a, b) {
	    return a + b
	}
	
	fn greet() {
	    print("Hello, World!")
	}
	
	result = add(x, y)
	
	if (result > 20) {
	    print("result is large")
	} elseif (result == 15) {
	    print("result is exact")
	} else {
	    print("result is small")
	}
	
	for (i in 0->5) {
	    print(i)
	}
	`

	parser, err := NewParser(input)
	if err != nil {
		t.Fatalf("Failed to create parser: %s", err)
	}

	ast, err := parser.GetAst()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if len(ast) == 0 {
		t.Fatalf("Expected non-empty AST")
	}

	varCount := 0
	funcCount := 0
	ifCount := 0
	forCount := 0

	for _, stmt := range ast {
		switch stmt.(type) {
		case VariableDeclaration:
			varCount++
		case FunctionDeclaration:
			funcCount++
		case IfStatement:
			ifCount++
		case ForLoop:
			forCount++
		}
	}

	if varCount < 3 {
		t.Errorf("Expected at least 3 variable declarations, got %d", varCount)
	}

	if funcCount != 2 {
		t.Errorf("Expected 2 function declarations, got %d", funcCount)
	}

	if ifCount != 1 {
		t.Errorf("Expected 1 if statement, got %d", ifCount)
	}

	if forCount != 1 {
		t.Errorf("Expected 1 for loop, got %d", forCount)
	}
}

func TestParseExpression(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "number literal",
			input:       "42",
			shouldError: false,
		},
		{
			name:        "string literal",
			input:       `"hello"`,
			shouldError: false,
		},
		{
			name:        "identifier",
			input:       "x",
			shouldError: false,
		},
		{
			name:        "addition expression",
			input:       "a + b",
			shouldError: false,
		},
		{
			name:        "multiplication and addition",
			input:       "a + b * c",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil && len(ast) == 0 {
				t.Errorf("Expected non-empty AST")
			}
		})
	}
}

func TestParseSemicolons(t *testing.T) {
	tests := []struct {
		name                   string
		input                  string
		expectedStatementCount int
		shouldError            bool
	}{
		{
			name:                   "statement with semicolon",
			input:                  "x = 5;",
			expectedStatementCount: 1,
			shouldError:            false,
		},
		{
			name:                   "multiple statements with semicolons",
			input:                  "x = 5; y = 10;",
			expectedStatementCount: 2,
			shouldError:            false,
		},
		{
			name:                   "statements without semicolons",
			input:                  "x = 5 y = 10",
			expectedStatementCount: 2,
			shouldError:            false,
		},
		{
			name:                   "mixed with and without semicolons",
			input:                  "x = 5; y = 10 z = 15;",
			expectedStatementCount: 3,
			shouldError:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != tt.expectedStatementCount {
					t.Errorf("Expected %d statements, got %d", tt.expectedStatementCount, len(ast))
				}
			}
		})
	}
}

func TestParseIdentifiersAndLiterals(t *testing.T) {
	tests := []struct {
		name         string
		input        string
		expectedType string
		shouldError  bool
	}{
		{
			name:         "number literal",
			input:        "42",
			expectedType: "NumberLiteralExpression",
			shouldError:  false,
		},
		{
			name:         "string literal",
			input:        `"test string"`,
			expectedType: "StringLiteralExpression",
			shouldError:  false,
		},
		{
			name:         "identifier",
			input:        "myVar",
			expectedType: "IdentifierExpression",
			shouldError:  false,
		},
		{
			name:         "underscore identifier",
			input:        "_private",
			expectedType: "IdentifierExpression",
			shouldError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}
			}
		})
	}
}

func TestParserOperatorPrecedence(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "multiplication before addition",
			input:       "x = 2 + 3 * 4",
			shouldError: false,
		},
		{
			name:        "comparison operators",
			input:       "x = a > b && c < d",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}
			}
		})
	}
}

func TestParseEmptyStructures(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "function with empty body",
			input:       "fn empty() { }",
			shouldError: false,
		},
		{
			name:        "if with empty body",
			input:       "if (x > 5) { }",
			shouldError: false,
		},
		{
			name:        "for with empty body",
			input:       "for (i in 0->5) { }",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}
			}
		})
	}
}

func TestParseTypeAnnotations(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		shouldError bool
	}{
		{
			name:        "int type",
			input:       "x: int = 5",
			shouldError: false,
		},
		{
			name:        "float type",
			input:       "y: float = 3.14",
			shouldError: false,
		},
		{
			name:        "string type",
			input:       `s: string = "hello"`,
			shouldError: false,
		},
		{
			name:        "bool type",
			input:       "b: bool = true",
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := NewParser(tt.input)
			if err != nil {
				t.Fatalf("Failed to create parser: %s", err)
			}

			ast, err := parser.GetAst()
			if err != nil && !tt.shouldError {
				t.Fatalf("Unexpected error: %s", err)
			}

			if err == nil && tt.shouldError {
				t.Fatalf("Expected error but got none")
			}

			if err == nil {
				if len(ast) != 1 {
					t.Errorf("Expected 1 statement, got %d", len(ast))
				}

				varDecl, ok := ast[0].(VariableDeclaration)
				if !ok {
					t.Errorf("Expected VariableDeclaration, got %T", ast[0])
					return
				}

				if varDecl.Value == nil {
					t.Errorf("Expected value in variable declaration")
				}
			}
		})
	}
}

func TestParseWithComments(t *testing.T) {
	input := `
	// This is a comment
	x = 5 // inline comment
	// Another comment
	y = 10
	/* Block comment */
	z = 15
	`

	parser, err := NewParser(input)
	if err != nil {
		t.Fatalf("Failed to create parser: %s", err)
	}

	ast, err := parser.GetAst()
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}

	if len(ast) != 3 {
		t.Errorf("Expected 3 statements, got %d", len(ast))
	}
}
