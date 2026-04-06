package parser

import (
	"errors"
	"testing"
)

func TestIsAlpha(t *testing.T) {
	if !isAlpha('a') || !isAlpha('Z') || isAlpha('1') {
		t.Errorf("IsAlpha failed")
	}
}

func TestIsDigit(t *testing.T) {
	if !isDigit('0') || !isDigit('9') || isDigit('a') {
		t.Errorf("IsDigit failed")
	}
}

func TestIsWhitespace(t *testing.T) {
	if !isWhitespace(' ') || !isWhitespace('\n') || isWhitespace('a') {
		t.Errorf("IsWhitespace failed")
	}
}

func TestIsKeyword(t *testing.T) {
	if !isKeyword("if") || !isKeyword("return") || !isKeyword("fn") || isKeyword("x") {
		t.Errorf("IsKeyword failed")
	}
}

func TestQuoteHandling(t *testing.T) {
	input := `print("Hello, World!")`
	expected := []Token{
		{Type: "Identifier", Value: "print"},
		{Type: "OpenParen", Value: "("},
		{Type: "StringLiteral", Value: "Hello, World!"},
		{Type: "CloseParen", Value: ")"},
		{Type: "EOF", Value: "EndOfFile"},
	}

	tokens, err := TokenizeInput(input)

	if err != nil {
		t.Errorf("TokenizeInput failed with error: %s", err)
	}

	if len(tokens) != len(expected) {
		t.Errorf("TokenizeInput failed: expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("TokenizeInput failed: expected token '%s', got '%s'", expected[i], token)
		}
	}
}

func TestInvalidQuoteHandling(t *testing.T) {
	input := `print("Hello, World!)`
	tokens, err := TokenizeInput(input)

	wantErr := errors.New("Unterminated string literal")

	if err.Error() != wantErr.Error() {
		t.Errorf("TokenizeInput failed: expected error 'Unterminated string literal', got '%s'", err)
	}

	if tokens != nil {
		t.Errorf("TokenizeInput failed: expected nil tokens, got %v", tokens)
	}
}

func TestUnregonizedToken(t *testing.T) {
	input := `@`
	tokens, err := TokenizeInput(input)

	wantErr := errors.New("Unrecognized character: @")

	if err.Error() != wantErr.Error() {
		t.Errorf("TokenizeInput failed: expected error 'Unrecognized character: @', got '%s'", err)
	}

	if tokens != nil {
		t.Errorf("TokenizeInput failed: expected nil tokens, got %v", tokens)
	}
}

func TestComments(t *testing.T) {
	input := `// This is a line comment
/* This is a block comment */
x = 5 /* This is another block comment */
/* Unclosed block comment`
	expected := []Token{
		{Type: "LineComment", Value: "//"},
		{Type: "BlockCommentStart", Value: "/*"},
		{Type: "BlockCommentEnd", Value: "*/"},
		{Type: "Identifier", Value: "x"},
		{Type: "Equals", Value: "="},
		{Type: "NumberLiteral", Value: "5"},
		{Type: "BlockCommentStart", Value: "/*"},
		{Type: "BlockCommentEnd", Value: "*/"},
		{Type: "BlockCommentStart", Value: "/*"},
		{Type: "EOF", Value: "EndOfFile"},
	}

	tokens, err := TokenizeInput(input)

	if err != nil {
		t.Errorf("TokenizeInput failed with error: %s", err)
	}

	if len(tokens) != len(expected) {
		t.Errorf("TokenizeInput failed: expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("TokenizeInput failed: expected token '%s', got '%s'", expected[i], token)
		}
	}
}

func TestTokenizeInput(t *testing.T) {
	input := "if x > 10 return x else return 0 \n \n 99 == 1 myVar = 5 << || [ ] + - * / %< ! , != <= >=  >> && ++ -- /* */ <-"
	expected := []Token{
		{Type: "If", Value: "if"},
		{Type: "Identifier", Value: "x"},
		{Type: "GreaterThan", Value: ">"},
		{Type: "NumberLiteral", Value: "10"},
		{Type: "Return", Value: "return"},
		{Type: "Identifier", Value: "x"},
		{Type: "Else", Value: "else"},
		{Type: "Return", Value: "return"},
		{Type: "NumberLiteral", Value: "0"},
		{Type: "NumberLiteral", Value: "99"},
		{Type: "EqualsEquals", Value: "=="},
		{Type: "NumberLiteral", Value: "1"},
		{Type: "Identifier", Value: "myVar"},
		{Type: "Equals", Value: "="},
		{Type: "NumberLiteral", Value: "5"},
		{Type: "LeftShift", Value: "<<"},
		{Type: "LogicalOr", Value: "||"},
		{Type: "OpenSquare", Value: "["},
		{Type: "CloseSquare", Value: "]"},
		{Type: "Plus", Value: "+"},
		{Type: "Minus", Value: "-"},
		{Type: "Asterisk", Value: "*"},
		{Type: "Slash", Value: "/"},
		{Type: "Percent", Value: "%"},
		{Type: "LessThan", Value: "<"},
		{Type: "ExclamationMark", Value: "!"},
		{Type: "Comma", Value: ","},
		{Type: "NotEquals", Value: "!="},
		{Type: "LessThanOrEqual", Value: "<="},
		{Type: "GreaterThanOrEqual", Value: ">="},
		{Type: "RightShift", Value: ">>"},
		{Type: "LogicalAnd", Value: "&&"},
		{Type: "Increment", Value: "++"},
		{Type: "Decrement", Value: "--"},
		{Type: "BlockCommentStart", Value: "/*"},
		{Type: "BlockCommentEnd", Value: "*/"},
		{Type: "ArrowLeft", Value: "<-"},
		{Type: "EOF", Value: "EndOfFile"},
	}

	tokens, err := TokenizeInput(input)

	if err != nil {
		t.Errorf("TokenizeInput failed with error: %s", err)
	}

	if len(tokens) != len(expected) {
		t.Errorf("TokenizeInput failed: expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("TokenizeInput failed: expected token '%s', got '%s'", expected[i], token)
		}
	}
}

func TestTokenizeProgram(t *testing.T) {
	input := `
	// Comment

/* 
    Multi-line comment
*/

x = 5           // Type inference
y: int = 6      // Explicit type annotation
z: float = 7.0  // Explicit type annotation

z = 8.5 // Forbidden

a? = 10 // Mutable
a = 15  // Allowed

fn a_function() {
    print("Hello, World!")

    return 42
}

a_function(); // Semicolon is optional

if (x > 3) {
    print("x is greater than 3 but actually {x}")
} elseif (x == 3) {
    print("x is equal to 3")
} else {
    print("x is less than 3")
}

for (i in 0->5) {
    print(i)
}
`

	expected := []Token{
		{Type: "LineComment", Value: "//"},
		{Type: "BlockCommentStart", Value: "/*"},
		{Type: "BlockCommentEnd", Value: "*/"},
		{Type: "Identifier", Value: "x"},
		{Type: "Equals", Value: "="},
		{Type: "NumberLiteral", Value: "5"},
		{Type: "LineComment", Value: "//"},
		{Type: "Identifier", Value: "y"},
		{Type: "Colon", Value: ":"},
		{Type: "IntType", Value: "int"},
		{Type: "Equals", Value: "="},
		{Type: "NumberLiteral", Value: "6"},
		{Type: "LineComment", Value: "//"},
		{Type: "Identifier", Value: "z"},
		{Type: "Colon", Value: ":"},
		{Type: "FloatType", Value: "float"},
		{Type: "Equals", Value: "="},
		{Type: "NumberLiteral", Value: "7"},
		{Type: "Period", Value: "."},
		{Type: "NumberLiteral", Value: "0"},
		{Type: "LineComment", Value: "//"},
		{Type: "Identifier", Value: "z"},
		{Type: "Equals", Value: "="},
		{Type: "NumberLiteral", Value: "8"},
		{Type: "Period", Value: "."},
		{Type: "NumberLiteral", Value: "5"},
		{Type: "LineComment", Value: "//"},
		{Type: "Identifier", Value: "a"},
		{Type: "QuestionMark", Value: "?"},
		{Type: "Equals", Value: "="},
		{Type: "NumberLiteral", Value: "10"},
		{Type: "LineComment", Value: "//"},
		{Type: "Identifier", Value: "a"},
		{Type: "Equals", Value: "="},
		{Type: "NumberLiteral", Value: "15"},
		{Type: "LineComment", Value: "//"},
		{Type: "Function", Value: "fn"},
		{Type: "Identifier", Value: "a_function"},
		{Type: "OpenParen", Value: "("},
		{Type: "CloseParen", Value: ")"},
		{Type: "OpenCurly", Value: "{"},
		{Type: "Identifier", Value: "print"},
		{Type: "OpenParen", Value: "("},
		{Type: "StringLiteral", Value: "Hello, World!"},
		{Type: "CloseParen", Value: ")"},
		{Type: "Return", Value: "return"},
		{Type: "NumberLiteral", Value: "42"},
		{Type: "CloseCurly", Value: "}"},
		{Type: "Identifier", Value: "a_function"},
		{Type: "OpenParen", Value: "("},
		{Type: "CloseParen", Value: ")"},
		{Type: "Semicolon", Value: ";"},
		{Type: "LineComment", Value: "//"},
		{Type: "If", Value: "if"},
		{Type: "OpenParen", Value: "("},
		{Type: "Identifier", Value: "x"},
		{Type: "GreaterThan", Value: ">"},
		{Type: "NumberLiteral", Value: "3"},
		{Type: "CloseParen", Value: ")"},
		{Type: "OpenCurly", Value: "{"},
		{Type: "Identifier", Value: "print"},
		{Type: "OpenParen", Value: "("},
		{Type: "StringLiteral", Value: "x is greater than 3 but actually {x}"},
		{Type: "CloseParen", Value: ")"},
		{Type: "CloseCurly", Value: "}"},
		{Type: "ElseIf", Value: "elseif"},
		{Type: "OpenParen", Value: "("},
		{Type: "Identifier", Value: "x"},
		{Type: "EqualsEquals", Value: "=="},
		{Type: "NumberLiteral", Value: "3"},
		{Type: "CloseParen", Value: ")"},
		{Type: "OpenCurly", Value: "{"},
		{Type: "Identifier", Value: "print"},
		{Type: "OpenParen", Value: "("},
		{Type: "StringLiteral", Value: "x is equal to 3"},
		{Type: "CloseParen", Value: ")"},
		{Type: "CloseCurly", Value: "}"},
		{Type: "Else", Value: "else"},
		{Type: "OpenCurly", Value: "{"},
		{Type: "Identifier", Value: "print"},
		{Type: "OpenParen", Value: "("},
		{Type: "StringLiteral", Value: "x is less than 3"},
		{Type: "CloseParen", Value: ")"},
		{Type: "CloseCurly", Value: "}"},
		{Type: "For", Value: "for"},
		{Type: "OpenParen", Value: "("},
		{Type: "Identifier", Value: "i"},
		{Type: "In", Value: "in"},
		{Type: "NumberLiteral", Value: "0"},
		{Type: "ArrowRight", Value: "->"},
		{Type: "NumberLiteral", Value: "5"},
		{Type: "CloseParen", Value: ")"},
		{Type: "OpenCurly", Value: "{"},
		{Type: "Identifier", Value: "print"},
		{Type: "OpenParen", Value: "("},
		{Type: "Identifier", Value: "i"},
		{Type: "CloseParen", Value: ")"},
		{Type: "CloseCurly", Value: "}"},
		{Type: "EOF", Value: "EndOfFile"},
	}

	tokens, err := TokenizeInput(input)

	if err != nil {
		t.Errorf("TokenizeInput failed with error: %s", err)
	}

	if len(tokens) != len(expected) {
		t.Errorf("TokenizeInput failed: expected %d tokens, got %d", len(expected), len(tokens))
	}

	for i, token := range tokens {
		if token != expected[i] {
			t.Errorf("TokenizeInput failed: expected token '%s', got '%s'", expected[i], token)
		}
	}
}

func TestEmpty(t *testing.T) {
	input := ""

	_, err := TokenizeInput(input)

	if err == nil {
		t.Errorf("TokenizeInput should have failed with an error for empty input")
	}

}
