package runtime

import (
	"strings"
	"testing"

	"afl/src/parser"
)

func TestEvaluateNumberLiteral(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.NumberLiteralExpression{Value: "42"}
	result := Evaluate(expr, env)

	if result.Type != "number" {
		t.Errorf("Expected type 'number', got '%s'", result.Type)
	}

	if result.Value.(float64) != 42.0 {
		t.Errorf("Expected value 42, got %v", result.Value)
	}
}

func TestEvaluateStringLiteral(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.StringLiteralExpression{Value: "hello"}
	result := Evaluate(expr, env)

	if result.Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", result.Type)
	}

	if result.Value.(string) != "hello" {
		t.Errorf("Expected value 'hello', got '%s'", result.Value)
	}
}

func TestStringInterpolation(t *testing.T) {
	t.Run("simple variable interpolation", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("name", RuntimeValue{Type: "string", Value: "John"}, false)

		expr := parser.StringLiteralExpression{Value: "Hello {name}"}
		result := Evaluate(expr, env)

		if result.Type != "string" {
			t.Errorf("Expected type 'string', got '%s'", result.Type)
		}

		if result.Value.(string) != "Hello John" {
			t.Errorf("Expected 'Hello John', got '%s'", result.Value)
		}
	})

	t.Run("numeric variable interpolation", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("x", RuntimeValue{Type: "number", Value: 42.0}, false)

		expr := parser.StringLiteralExpression{Value: "The answer is {x}"}
		result := Evaluate(expr, env)

		if result.Value.(string) != "The answer is 42" {
			t.Errorf("Expected 'The answer is 42', got '%s'", result.Value)
		}
	})

	t.Run("multiple interpolations", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("first", RuntimeValue{Type: "string", Value: "John"}, false)
		env.declareVariable("last", RuntimeValue{Type: "string", Value: "Doe"}, false)

		expr := parser.StringLiteralExpression{Value: "{first} {last}"}
		result := Evaluate(expr, env)

		if result.Value.(string) != "John Doe" {
			t.Errorf("Expected 'John Doe', got '%s'", result.Value)
		}
	})

	t.Run("undefined variable keeps placeholder", func(t *testing.T) {
		env := NewEnvironment(nil)

		expr := parser.StringLiteralExpression{Value: "Hello {name}"}
		result := Evaluate(expr, env)

		if result.Value.(string) != "Hello {name}" {
			t.Errorf("Expected 'Hello {name}', got '%s'", result.Value)
		}
	})

	t.Run("string with no interpolation", func(t *testing.T) {
		env := NewEnvironment(nil)

		expr := parser.StringLiteralExpression{Value: "just a string"}
		result := Evaluate(expr, env)

		if result.Value.(string) != "just a string" {
			t.Errorf("Expected 'just a string', got '%s'", result.Value)
		}
	})

	t.Run("interpolation with float numbers", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("pi", RuntimeValue{Type: "number", Value: 3.14}, false)

		expr := parser.StringLiteralExpression{Value: "Pi is approximately {pi}"}
		result := Evaluate(expr, env)

		if result.Value.(string) != "Pi is approximately 3.14" {
			t.Errorf("Expected 'Pi is approximately 3.14', got '%s'", result.Value)
		}
	})

	t.Run("nested braces not supported (only first level)", func(t *testing.T) {
		env := NewEnvironment(nil)
		env.declareVariable("obj", RuntimeValue{Type: "object", Value: map[string]RuntimeValue{}}, false)

		expr := parser.StringLiteralExpression{Value: "Object: {obj}"}
		result := Evaluate(expr, env)

		// Should contain the string representation of the object
		if !strings.Contains(result.Value.(string), "map") {
			t.Logf("Object interpolation: %s", result.Value)
		}
	})
}

func TestEvaluateAddition(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.BinaryExpression{
		Left:     parser.NumberLiteralExpression{Value: "5"},
		Operator: parser.Token{Type: parser.Plus, Value: "+"},
		Right:    parser.NumberLiteralExpression{Value: "3"},
	}
	result := Evaluate(expr, env)

	if result.Type != "number" {
		t.Errorf("Expected type 'number', got '%s'", result.Type)
	}

	if result.Value.(float64) != 8.0 {
		t.Errorf("Expected value 8, got %v", result.Value)
	}
}

func TestEvaluateSubtraction(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.BinaryExpression{
		Left:     parser.NumberLiteralExpression{Value: "10"},
		Operator: parser.Token{Type: parser.Minus, Value: "-"},
		Right:    parser.NumberLiteralExpression{Value: "3"},
	}
	result := Evaluate(expr, env)

	if result.Value.(float64) != 7.0 {
		t.Errorf("Expected value 7, got %v", result.Value)
	}
}

func TestEvaluateMultiplication(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.BinaryExpression{
		Left:     parser.NumberLiteralExpression{Value: "4"},
		Operator: parser.Token{Type: parser.Asterisk, Value: "*"},
		Right:    parser.NumberLiteralExpression{Value: "5"},
	}
	result := Evaluate(expr, env)

	if result.Value.(float64) != 20.0 {
		t.Errorf("Expected value 20, got %v", result.Value)
	}
}

func TestEvaluateDivision(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.BinaryExpression{
		Left:     parser.NumberLiteralExpression{Value: "20"},
		Operator: parser.Token{Type: parser.Slash, Value: "/"},
		Right:    parser.NumberLiteralExpression{Value: "4"},
	}
	result := Evaluate(expr, env)

	if result.Value.(float64) != 5.0 {
		t.Errorf("Expected value 5, got %v", result.Value)
	}
}

func TestEvaluateComparison(t *testing.T) {
	env := NewEnvironment(nil)

	t.Run("less than", func(t *testing.T) {
		expr := parser.BinaryExpression{
			Left:     parser.NumberLiteralExpression{Value: "3"},
			Operator: parser.Token{Type: parser.LessThan, Value: "<"},
			Right:    parser.NumberLiteralExpression{Value: "5"},
		}
		result := Evaluate(expr, env)
		if result.Value.(float64) != 1.0 {
			t.Errorf("Expected true (1), got %v", result.Value)
		}
	})

	t.Run("greater than", func(t *testing.T) {
		expr := parser.BinaryExpression{
			Left:     parser.NumberLiteralExpression{Value: "5"},
			Operator: parser.Token{Type: parser.GreaterThan, Value: ">"},
			Right:    parser.NumberLiteralExpression{Value: "3"},
		}
		result := Evaluate(expr, env)
		if result.Value.(float64) != 1.0 {
			t.Errorf("Expected true (1), got %v", result.Value)
		}
	})

	t.Run("equals", func(t *testing.T) {
		expr := parser.BinaryExpression{
			Left:     parser.NumberLiteralExpression{Value: "5"},
			Operator: parser.Token{Type: parser.EqualsEquals, Value: "=="},
			Right:    parser.NumberLiteralExpression{Value: "5"},
		}
		result := Evaluate(expr, env)
		if result.Value.(float64) != 1.0 {
			t.Errorf("Expected true (1), got %v", result.Value)
		}
	})
}

func TestEvaluateStringConcatenation(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.BinaryExpression{
		Left:     parser.StringLiteralExpression{Value: "hello"},
		Operator: parser.Token{Type: parser.Plus, Value: "+"},
		Right:    parser.StringLiteralExpression{Value: "world"},
	}
	result := Evaluate(expr, env)

	if result.Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", result.Type)
	}

	if result.Value.(string) != "helloworld" {
		t.Errorf("Expected 'helloworld', got '%s'", result.Value)
	}
}

func TestEvaluateNegation(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.BinaryExpression{
		Left:     nil,
		Operator: parser.Token{Type: parser.Minus, Value: "-"},
		Right:    parser.NumberLiteralExpression{Value: "5"},
	}
	result := Evaluate(expr, env)

	if result.Value.(float64) != -5.0 {
		t.Errorf("Expected value -5, got %v", result.Value)
	}
}

func TestEvaluateVariableDeclaration(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.VariableDeclaration{
		Name:      "x",
		IsMutable: false,
		Value:     parser.NumberLiteralExpression{Value: "42"},
	}
	result := Evaluate(expr, env)

	if result.Type != "number" || result.Value.(float64) != 42.0 {
		t.Errorf("Expected number 42, got %v", result.Value)
	}

	// Verify variable was stored
	val, err := env.getVariable("x")
	if err != nil {
		t.Errorf("Variable 'x' not found: %v", err)
	}

	runtimeVal := val.(RuntimeValue)
	if runtimeVal.Type != "number" || runtimeVal.Value.(float64) != 42.0 {
		t.Errorf("Expected variable x to be 42, got %v", runtimeVal.Value)
	}
}

func TestEvaluateIdentifier(t *testing.T) {
	env := NewEnvironment(nil)

	// Declare a variable
	env.declareVariable("x", RuntimeValue{Type: "number", Value: 42.0}, false)

	// Evaluate the identifier
	expr := parser.IdentifierExpression{Symbol: "x"}
	result := Evaluate(expr, env)

	if result.Type != "number" || result.Value.(float64) != 42.0 {
		t.Errorf("Expected 42, got %v", result.Value)
	}
}

func TestEvaluateFunctionDeclaration(t *testing.T) {
	env := NewEnvironment(nil)
	expr := parser.FunctionDeclaration{
		Name:       "add",
		Parameters: []string{"a", "b"},
		Body: []parser.Expression{
			parser.BinaryExpression{
				Left:     parser.IdentifierExpression{Symbol: "a"},
				Operator: parser.Token{Type: parser.Plus, Value: "+"},
				Right:    parser.IdentifierExpression{Symbol: "b"},
			},
		},
	}

	Evaluate(expr, env)

	// Verify function was stored
	val, err := env.getVariable("add")
	if err != nil {
		t.Errorf("Function 'add' not found: %v", err)
	}

	runtimeVal := val.(RuntimeValue)
	if runtimeVal.Type != "function" {
		t.Errorf("Expected type 'function', got '%s'", runtimeVal.Type)
	}
}

func TestEvaluateFunctionCall(t *testing.T) {
	env := NewEnvironment(nil)

	// Declare function
	env.declareVariable("add", RuntimeValue{
		Type: "function",
		Value: &FunctionValue{
			Parameters: []string{"a", "b"},
			Body: []parser.Expression{
				parser.BinaryExpression{
					Left:     parser.IdentifierExpression{Symbol: "a"},
					Operator: parser.Token{Type: parser.Plus, Value: "+"},
					Right:    parser.IdentifierExpression{Symbol: "b"},
				},
			},
			Env: env,
		},
	}, false)

	// Call function
	expr := parser.CallExpression{
		Caller: parser.IdentifierExpression{Symbol: "add"},
		Args: []parser.Expression{
			parser.NumberLiteralExpression{Value: "3"},
			parser.NumberLiteralExpression{Value: "5"},
		},
	}

	result := Evaluate(expr, env)

	if result.Type != "number" || result.Value.(float64) != 8.0 {
		t.Errorf("Expected 8, got %v", result.Value)
	}
}

func TestEvaluateIfStatement(t *testing.T) {
	env := NewEnvironment(nil)

	t.Run("condition true", func(t *testing.T) {
		expr := parser.IfStatement{
			Condition: parser.NumberLiteralExpression{Value: "1"},
			IfBranch: parser.Branch{
				Condition: parser.NumberLiteralExpression{Value: "1"},
				Body: []parser.Expression{
					parser.NumberLiteralExpression{Value: "42"},
				},
			},
		}
		result := Evaluate(expr, env)

		if result.Type != "number" || result.Value.(float64) != 42.0 {
			t.Errorf("Expected 42, got %v", result.Value)
		}
	})

	t.Run("condition false, has else", func(t *testing.T) {
		expr := parser.IfStatement{
			Condition: parser.NumberLiteralExpression{Value: "0"},
			IfBranch: parser.Branch{
				Condition: parser.NumberLiteralExpression{Value: "0"},
				Body: []parser.Expression{
					parser.NumberLiteralExpression{Value: "10"},
				},
			},
			ElseBranch: parser.Branch{
				Body: []parser.Expression{
					parser.NumberLiteralExpression{Value: "20"},
				},
			},
		}
		result := Evaluate(expr, env)

		if result.Type != "number" || result.Value.(float64) != 20.0 {
			t.Errorf("Expected 20, got %v", result.Value)
		}
	})
}

func TestEvaluateProgram(t *testing.T) {
	env := NewEnvironment(nil)

	ast := []parser.Expression{
		parser.VariableDeclaration{
			Name:      "x",
			IsMutable: false,
			Value:     parser.NumberLiteralExpression{Value: "10"},
		},
		parser.VariableDeclaration{
			Name:      "y",
			IsMutable: false,
			Value:     parser.NumberLiteralExpression{Value: "20"},
		},
		parser.BinaryExpression{
			Left:     parser.IdentifierExpression{Symbol: "x"},
			Operator: parser.Token{Type: parser.Plus, Value: "+"},
			Right:    parser.IdentifierExpression{Symbol: "y"},
		},
	}

	result := EvaluateProgram(ast, env)

	if result.Type != "number" || result.Value.(float64) != 30.0 {
		t.Errorf("Expected 30, got %v", result.Value)
	}
}

func TestEvaluateObjectLiteral(t *testing.T) {
	env := NewEnvironment(nil)

	properties := []parser.Property{
		{
			Key:   "name",
			Value: parser.StringLiteralExpression{Value: "John"},
		},
		{
			Key:   "age",
			Value: parser.NumberLiteralExpression{Value: "30"},
		},
	}

	result := Evaluate(properties, env)

	if result.Type != "object" {
		t.Errorf("Expected type 'object', got '%s'", result.Type)
	}

	objMap := result.Value.(map[string]RuntimeValue)
	if len(objMap) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(objMap))
	}

	if objMap["name"].Value.(string) != "John" {
		t.Errorf("Expected name='John', got '%s'", objMap["name"].Value)
	}
}

func TestEvaluateMemberExpression(t *testing.T) {
	env := NewEnvironment(nil)

	// Create an object
	objMap := map[string]RuntimeValue{
		"name": {Type: "string", Value: "John"},
		"age":  {Type: "number", Value: 30.0},
	}
	env.declareVariable("person", RuntimeValue{Type: "object", Value: objMap}, false)

	// Access member
	expr := parser.MemberExpression{
		Object:   parser.IdentifierExpression{Symbol: "person"},
		Property: "name",
	}

	result := Evaluate(expr, env)

	if result.Type != "string" || result.Value.(string) != "John" {
		t.Errorf("Expected 'John', got %v", result.Value)
	}
}
